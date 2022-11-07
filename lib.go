package go_premiumize

import (
	"github.com/BRUHItsABunny/go-premiumize/api"
	"github.com/BRUHItsABunny/go-premiumize/client"
	"net/http"
)

func GetPremiumizeClient(session *api.PremiumizeSession, hClient *http.Client) *client.PremiumizeClient {
	return client.NewPremiumizeClient(session, hClient)
}

func GetPremiumizeAPISession(apiKey string) *api.PremiumizeSession {
	return &api.PremiumizeSession{SessionType: "apikey", AuthToken: apiKey}
}
