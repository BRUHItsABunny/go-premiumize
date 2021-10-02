package main

import (
	"context"
	"flag"
	"fmt"
	gokhttp "github.com/BRUHItsABunny/gOkHttp"
	"github.com/BRUHItsABunny/go-premiumize/api"
	"github.com/BRUHItsABunny/go-premiumize/client"
	"github.com/sheerun/queue"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// Sample application that I will end up actually using
// go run main.go -apikey APIKEY -folder BUNBUNBUN -recursion -threads 6
type CloneParameters struct {
	APIKey           string
	DownloadThreads  int
	Folder           string
	Recursive        bool
	ProgressInterval int
	Proxy            string
}

type DownloadableItem struct {
	Item   *api.PremiumizeItem
	Crumbs []*api.BreadCrumb
}

func main() {
	// Parse parameters
	programParams := new(CloneParameters)
	flag.StringVar(&programParams.APIKey, "apikey", "", "This is our APIKey - not needed, if missing it will authenticate via device code")
	flag.IntVar(&programParams.DownloadThreads, "threads", 1, "This is how many files we download in parallel (min=1, max=6)")
	flag.StringVar(&programParams.Folder, "folder", "", "This is the folder we will start crawling in")
	flag.BoolVar(&programParams.Recursive, "recursion", false, "This controls if we want all files inside all folders of the folder you selected or just all files in the folder you selected")
	flag.IntVar(&programParams.ProgressInterval, "pinterval", 5, "This is how many seconds we wait in between each progress print")
	flag.StringVar(&programParams.Proxy, "proxy", "", "This argument is for proxying this program (format: proto://ip:port)")

	flag.Parse()

	// Prepare our client
	var session *api.PremiumizeSession
	if len(programParams.APIKey) > 0 {
		session = &api.PremiumizeSession{SessionType: "apikey", AuthToken: programParams.APIKey}
	}
	hClient := gokhttp.GetHTTPDownloadClient(gokhttp.DefaultGOKHTTPOptions) // A client with sufficient timeouts for downloading
	if len(programParams.Proxy) > 0 {
		err := hClient.SetProxy(programParams.Proxy)
		if err != nil {
			panic(err)
		}
	}
	pClient := client.NewPremiumizeClient(session, hClient.Client)
	ctx := context.Background()

	// Start worker goroutines for downloads
	downloadQueue := queue.New()
	if programParams.DownloadThreads < 1 {
		programParams.DownloadThreads = 1
	} else {
		if programParams.DownloadThreads > 6 {
			programParams.DownloadThreads = 6
		}
	}
	stopChan := make(chan bool, programParams.DownloadThreads)
	wg := &sync.WaitGroup{}
	wg.Add(programParams.DownloadThreads)
	downloaded := int32(0)
	for i := 0; i < programParams.DownloadThreads; i++ {
		go Download(i, wg, stopChan, downloadQueue, &downloaded, pClient)
	}

	// Find the folder, check against name and id?
	offset := 0
	if strings.HasPrefix(programParams.Folder, "My Files/") || strings.HasPrefix(programParams.Folder, "/") {
		offset = 1
	}
	crumbs := strings.Split(programParams.Folder, "/")[offset:]
	folderID := "" // Start in root
	lastCrumb := len(crumbs) - 1
	for i, crumb := range crumbs {
		listResp, err := pClient.FoldersList(ctx, &api.FolderListRequest{ID: folderID})
		if err == nil {
			for _, item := range listResp.Content {
				if item.Type == "folder" {
					if item.Name == crumb {
						folderID = item.ID
						if i == lastCrumb {
							break
						}
						continue
					}
				}
			}
		} else {
			panic(err)
		}
	}

	// Start crawling and feeding links into download queue
	downloading, err := Crawl(ctx, downloadQueue, pClient, programParams.Recursive, folderID)
	if err == nil {
		fmt.Println(fmt.Sprintf("Found %d files to download", downloading))
		// Wait until download queue is empty and goroutines are finished
		sysChan := make(chan os.Signal, 1)
		signal.Notify(sysChan, os.Interrupt, syscall.SIGTERM)
		shouldStop := int(downloaded) == downloading
		if programParams.ProgressInterval < 1 {
			programParams.ProgressInterval = 1
		}
		ticker := time.NewTicker(time.Duration(programParams.ProgressInterval) * time.Second)
		for {
			if shouldStop {
				fmt.Println("Stopping workers and exiting loop")
				// Trigger workers exit
				for i := 0; i < programParams.DownloadThreads; i++ {
					stopChan <- true
				}
				break
			}
			select {
			case <-ticker.C: // Catch empty queue (eventually...)
				fmt.Println(fmt.Sprintf("[%s] Ticker: %d downloaded out of %d", time.Now().Format(time.RFC3339), downloaded, downloading))
				shouldStop = int(downloaded) == downloading
				break
			case <-sysChan: // Catch ctrl+C
				fmt.Println("Program kill detected, killing us softly... (tries to wrap up download)")
				shouldStop = true
				break
			}
		}
	} else {
		fmt.Println("Stopping workers and exiting loop")
		// Trigger workers exit
		for i := 0; i < programParams.DownloadThreads; i++ {
			stopChan <- true
		}
	}

	wg.Wait() // Wait for workers
	os.Exit(0)
}

func Download(id int, wg *sync.WaitGroup, stopper chan bool, queue *queue.Queue, counter *int32, pClient *client.PremiumizeClient) {
	// Single threaded resumable downloader
	id++
	shouldStop := false
	dir := ""
	var item *DownloadableItem
	ticker := time.NewTicker(time.Duration(1) * time.Second)
	time.Sleep(5 * time.Second)
	for {
		if shouldStop {
			fmt.Println(fmt.Sprintf("[%d] Worker exiting", id))
			break
		}
		select {
		case <-stopper:
			shouldStop = true
			// fmt.Println(fmt.Sprintf("[%d] Worker received stop", id))
			break
		case <-ticker.C:
			if queue.Length() > 0 {
				item = queue.Pop().(*DownloadableItem)
				if item.Item != nil && item.Item.Type == "file" {
					// Figure out what directory we are working in
					dir = ""
					for _, crumb := range item.Crumbs {
						dir += crumb.Name + "/"
					}
					_ = os.MkdirAll(dir[:len(dir)-1], 0777)
					fmt.Println(fmt.Sprintf("[%d] Downloading %s", id, dir+item.Item.Name))

					// Check if file exists
					start, end := int64(0), int64(*item.Item.Size)
					f, err := os.OpenFile(dir+item.Item.Name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
					if err == nil {
						stat, err := os.Stat(dir + item.Item.Name)
						if err == nil && stat != nil {
							start = stat.Size()
						}
					}
					// If file exists, update our range header for the resume
					req, err := http.NewRequest("GET", *item.Item.Link, nil)
					if err == nil {
						req.Header["range"] = []string{fmt.Sprintf("bytes=%d-%d", start, end)}
						resp, err := pClient.Client.Do(req)
						if err == nil {
							// Write to file
							_, _ = io.Copy(f, resp.Body)
							_ = f.Close()
							_ = resp.Body.Close()
						} else {
							fmt.Println("req fire err: ", err)
						}
					} else {
						fmt.Println("req err: ", err)
					}
					fmt.Println(fmt.Sprintf("[%d] Done downloading %s", id, item.Item.Name))
					atomic.AddInt32(counter, 1) // Count but no freeze
				}
			}
			break
		}
	}
	wg.Done()
	fmt.Println(fmt.Sprintf("[%d] Worker exited", id))
}

func Crawl(ctx context.Context, downloadQueue *queue.Queue, pClient *client.PremiumizeClient, recursive bool, folderID string) (int, error) {
	done, doneTemp := 0, 0
	resp, err := pClient.FoldersList(ctx, &api.FolderListRequest{ID: folderID, BreadCrumbs: true})
	if err == nil {
		crumbs := make([]*api.BreadCrumb, 0)
		if len(resp.BreadCrumbs) > 1 {
			crumbs = resp.BreadCrumbs[1:]
		}
		for _, item := range resp.Content {
			if item.Type == "folder" && recursive {
				time.Sleep(time.Duration(1) * time.Second)
				doneTemp, err = Crawl(ctx, downloadQueue, pClient, recursive, item.ID)
				if err != nil {
					fmt.Println("Recursion err: ", err)
					break
				}
				done += doneTemp
			} else {
				if item.Type == "file" {
					// fmt.Println("Found file: ", item.Name)
					downloadQueue.Append(&DownloadableItem{
						Item:   item,
						Crumbs: crumbs,
					})
					done++
				}
			}
		}
	}

	return done, err
}
