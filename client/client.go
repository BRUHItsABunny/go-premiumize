package client

import (
	"context"
	"encoding/json"
	"github.com/BRUHItsABunny/go-premiumize/api"
	"github.com/BRUHItsABunny/go-premiumize/constants"
	"io"
	"net/http"
)

func NewPremiumizeClient(session *api.PremiumizeSession, client *http.Client) *PremiumizeClient {

	if session == nil {
		session = &api.PremiumizeSession{SessionType: constants.TokenResponseType}
	}
	if client == nil {
		client = http.DefaultClient
	}

	return &PremiumizeClient{Session: session, Client: client}
}

type PremiumizeClient struct {
	Session *api.PremiumizeSession
	Client  *http.Client
}

func (c *PremiumizeClient) ShouldAuthenticate() bool { // Does not do a remote check
	if c.Session == nil {
		c.Session = &api.PremiumizeSession{SessionType: constants.TokenResponseType}
	}
	return !(len(c.Session.AuthToken) > 0)
}

func (c *PremiumizeClient) Token(ctx context.Context, r *api.TokenRequest) (*api.TokenResponse, error) {
	var (
		req       *http.Request
		resp      *http.Response
		bodyBytes []byte
		result    = new(api.TokenResponse)
		err       error
	)

	req, err = api.Token(ctx, r)
	if err == nil {
		resp, err = c.Client.Do(req)
		if err == nil {
			bodyBytes, err = io.ReadAll(resp.Body)
			if err == nil {
				err = json.Unmarshal(bodyBytes, result)
			}
		}
	}

	return result, err
}

func (c *PremiumizeClient) FoldersList(ctx context.Context, r *api.FolderListRequest) (*api.FolderListResponse, error) {
	var (
		req       *http.Request
		resp      *http.Response
		bodyBytes []byte
		result    = new(api.FolderListResponse)
		err       error
	)

	req, err = api.FolderList(ctx, c.Session, r)
	if err == nil {
		resp, err = c.Client.Do(req)
		if err == nil {
			bodyBytes, err = io.ReadAll(resp.Body)
			if err == nil {
				err = json.Unmarshal(bodyBytes, result)
			}
		}
	}

	return result, err
}
