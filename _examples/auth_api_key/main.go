package main

import (
	"context"
	"fmt"
	gokhttp "github.com/BRUHItsABunny/gOkHttp"
	"github.com/BRUHItsABunny/go-premiumize/api"
	"github.com/BRUHItsABunny/go-premiumize/client"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	hClient := gokhttp.GetHTTPClient(gokhttp.DefaultGOKHTTPOptions)
	_ = hClient.SetProxy("http://127.0.0.1:8888")

	session := &api.PremiumizeSession{SessionType: "apikey", AuthToken: "APIKEYHERE"}
	pClient := client.NewPremiumizeClient(session, hClient.Client)
	ctx := context.Background()

	if !pClient.ShouldAuthenticate() {
		fResp, err := pClient.FoldersList(ctx, &api.FolderListRequest{})
		if err == nil {
			fmt.Println(spew.Sdump(fResp))
		} else {
			fmt.Println("Error listing directory: ", err)
		}
	}
}
