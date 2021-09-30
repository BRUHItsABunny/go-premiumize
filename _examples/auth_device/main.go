package main

import (
	"context"
	"fmt"
	gokhttp "github.com/BRUHItsABunny/gOkHttp"
	"github.com/BRUHItsABunny/go-premiumize/api"
	"github.com/BRUHItsABunny/go-premiumize/client"
	"github.com/davecgh/go-spew/spew"
	"time"
)

func main() {
	hClient := gokhttp.GetHTTPClient(gokhttp.DefaultGOKHTTPOptions)
	_ = hClient.SetProxy("http://127.0.0.1:8888")

	pClient := client.NewPremiumizeClient(nil, hClient.Client)
	ctx := context.Background()

	req := api.NewTokenRequest()
	resp, err := pClient.Token(ctx, req)
	if err == nil {
		if resp.Error == nil {
			if resp.UserCode != nil {
				fmt.Println("Enter this code: " + *resp.UserCode)
				fmt.Println("On the webpage: https://premiumize.me/device")
				req.SetCodeAndExpiration(*resp.DeviceCode, *resp.ExpiresIn)

				for {
					time.Sleep(time.Second * time.Duration(5))
					resp, err = pClient.Token(ctx, req)
					if err == nil {
						if resp.Error == nil {
							if resp.AccessToken != nil {
								fmt.Println("Logged in! token: " + *resp.AccessToken)
								pClient.Session.AuthToken = *resp.AccessToken
								break
							}
						} else {
							fmt.Println(*resp.ErrorDescription)
						}
					} else {
						fmt.Println("Error occurred: ", err)
					}
				}

				if !pClient.ShouldAuthenticate() {
					fResp, err := pClient.FoldersList(ctx, &api.FolderListRequest{})
					if err == nil {
						fmt.Println(spew.Sdump(fResp))
					} else {
						fmt.Println("Error listing directory: ", err)
					}
				}
			}
		} else {
			fmt.Println("Error occurred: " + *resp.Error)
		}
	} else {
		fmt.Println("Error occurred: ", err)
	}
}
