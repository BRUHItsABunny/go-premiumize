package api

import (
	"errors"
	"fmt"
	"github.com/BRUHItsABunny/go-premiumize/constants"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	ErrExpired = errors.New("expired")
	ErrEmpty   = errors.New("one or more values cannot be empty")
)

type PremiumizeSession struct {
	AuthToken   string `json:"auth_token"`
	SessionType string `json:"session_type"`
	// Expires time.Time
}

type PremiumizeAPIResponse struct {
	Status  string  `json:"status"`
	Message *string `json:"message,omitempty"`
}

type PremiumizeItem struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Type            string  `json:"type"` // file / folder
	Size            *int    `json:"size,omitempty"`
	CreatedAt       *int    `json:"created_at,omitempty"`
	MIMEType        *string `json:"mime_type,omitempty"`
	TranscodeStatus *string `json:"transcode_status,omitempty"`
	Link            *string `json:"link,omitempty"`
	StreamLink      *string `json:"stream_link,omitempty"`
	VirusScan       *string `json:"virus_scan,omitempty"`
}

type PremiumizeTransfer struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Status   string   `json:"status"`
	Message  *string  `json:"message,omitempty"`
	Progress *float64 `json:"progress,omitempty"`
	Source   *string  `json:"src,omitempty"`
	FolderID *string  `json:"folder_id,omitempty"`
	FileID   *string  `json:"file_id,omitempty"`
}

func NewTokenRequest() *TokenRequest {
	return &TokenRequest{
		ClientID:  constants.ClientID,
		GrantType: constants.TokenResponseType,
		Code:      "",
		Expires:   time.Time{},
	}
}

type TokenRequest struct {
	ClientID  string
	GrantType string
	Code      string
	Expires   time.Time
}

func (r *TokenRequest) SetCodeAndExpiration(code string, expires int) {
	r.Code = code
	r.Expires = time.Now().Add(time.Second * time.Duration(expires)) // Probably 600
}

func (r *TokenRequest) ToHTTPBody() (io.Reader, error) {
	var (
		err    error
		result io.Reader
	)

	if len(r.Code) > 0 {
		if !r.Expires.IsZero() && r.Expires.After(time.Now()) {
			// Generate request to check auth
			result = strings.NewReader(url.Values{
				"grant_type": {r.GrantType},
				"client_id":  {r.ClientID},
				"code":       {r.Code},
			}.Encode())
		} else {
			err = ErrExpired
		}
	} else {
		// Generate request to initialize auth
		result = strings.NewReader(url.Values{
			"response_type": {r.GrantType},
			"client_id":     {r.ClientID},
		}.Encode())
	}

	return result, err
}

type TokenResponse struct {
	UserCode         *string `json:"user_code,omitempty"`
	DeviceCode       *string `json:"device_code,omitempty"`
	ExpiresIn        *int    `json:"expires_in,omitempty"`
	Scope            *string `json:"scope,omitempty"`
	AccessToken      *string `json:"access_token,omitempty"`
	Error            *string `json:"error,omitempty"`
	ErrorDescription *string `json:"error_description,omitempty"`
}

type FolderListRequest struct {
	ID          string
	BreadCrumbs bool
}

func (r *FolderListRequest) ToQuery(session *PremiumizeSession) string {
	var (
		result = url.Values{}
	)

	if len(r.ID) > 0 {
		result["id"] = []string{r.ID}
	}
	if r.BreadCrumbs {
		result["includeBreadCrumbs"] = []string{strconv.FormatBool(r.BreadCrumbs)}
	}
	result = authParams(session, result)
	return result.Encode()
}

type FolderListResponse struct {
	PremiumizeAPIResponse
	Content  []*PremiumizeItem `json:"content"`
	Name     string            `json:"name"`
	ParentID string            `json:"parent_id"`
	FolderID string            `json:"folder_id"`
}

type FolderCreateRequest struct {
	Name   string
	Parent string
}

func (r *FolderCreateRequest) ToHTTPBody(session *PremiumizeSession) (io.Reader, error) {
	var (
		err    error
		result io.Reader
	)

	if len(r.Name) > 0 {
		params := url.Values{
			"name": {r.Name},
		}
		if len(r.Parent) > 0 {
			params["parent_id"] = []string{r.Parent}
		}
		params = authParams(session, params)
		result = strings.NewReader(params.Encode())
	} else {
		err = ErrEmpty
	}

	return result, err
}

type FolderRenameRequest struct {
	Name string
	ID   string
}

func (r *FolderRenameRequest) ToHTTPBody(session *PremiumizeSession) (io.Reader, error) {
	var (
		err    error
		result io.Reader
	)

	if len(r.Name) > 0 && len(r.ID) > 0 {
		params := url.Values{
			"name": {r.Name},
			"id":   {r.ID},
		}
		params = authParams(session, params)
		result = strings.NewReader(params.Encode())
	} else {
		err = ErrEmpty
	}

	return result, err
}

type FolderPasteRequest struct {
	Items               []*PremiumizeItem
	DestinationFolderID string
}

func (r *FolderPasteRequest) ToHTTPBody(session *PremiumizeSession) (io.Reader, error) {
	var (
		err    error
		result io.Reader
	)

	if len(r.Items) > 0 && len(r.DestinationFolderID) > 0 {
		items := url.Values{"id": {r.DestinationFolderID}}
		var baseKey string
		for i, item := range r.Items {
			baseKey = fmt.Sprintf("items[%d]", i)
			items[baseKey+"[id]"] = []string{item.ID}
			items[baseKey+"[type]"] = []string{item.Type}
		}
		items = authParams(session, items)
		result = strings.NewReader(items.Encode())
	} else {
		err = ErrEmpty
	}

	return result, err
}

type FolderDeleteRequest struct {
	FolderID string
}

func (r *FolderDeleteRequest) ToHTTPBody(session *PremiumizeSession) (io.Reader, error) {
	var (
		err    error
		result io.Reader
	)

	if len(r.FolderID) > 0 {
		params := url.Values{"id": {r.FolderID}}
		params = authParams(session, params)
		result = strings.NewReader(params.Encode())
	} else {
		err = ErrEmpty
	}

	return result, err
}

type FolderUploadInfoRequest struct {
	ID string
}

func (r *FolderUploadInfoRequest) ToQuery(session *PremiumizeSession) string {
	var (
		result = url.Values{}
	)

	if len(r.ID) > 0 {
		result["id"] = []string{r.ID}
	}

	result = authParams(session, result)
	return result.Encode()
}

type FolderSearchRequest struct {
	Query string
}

func (r *FolderSearchRequest) ToQuery(session *PremiumizeSession) string {
	var (
		result = url.Values{}
	)

	if len(r.Query) > 0 {
		result["q"] = []string{r.Query}
	}

	result = authParams(session, result)
	return result.Encode()
}
