package api

import (
	"context"
	"github.com/BRUHItsABunny/go-premiumize/constants"
	"net/http"
)

func Token(ctx context.Context, r *TokenRequest) (*http.Request, error) {
	var req *http.Request

	body, err := r.ToHTTPBody()
	if err == nil {
		req, err = http.NewRequestWithContext(ctx, "POST", constants.TokenURL, body)
		if err == nil {
			req = readyHeader(req, nil)
		}
	}

	return req, err
}

func FolderList(ctx context.Context, s *PremiumizeSession, r *FolderListRequest) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", constants.EndpointFolderList, nil)
	if err == nil {
		req = readyHeader(req, s)
		req.URL.RawQuery = r.ToQuery(s)
	}

	return req, err
}

func FolderCreate(ctx context.Context, s *PremiumizeSession, r *FolderCreateRequest) (*http.Request, error) {
	var req *http.Request

	body, err := r.ToHTTPBody(s)
	if err == nil {
		req, err = http.NewRequestWithContext(ctx, "POST", constants.EndpointFolderCreate, body)
		if err == nil {
			req = readyHeader(req, s)
		}
	}

	return req, err
}

func FolderRename(ctx context.Context, s *PremiumizeSession, r *FolderRenameRequest) (*http.Request, error) {
	var req *http.Request

	body, err := r.ToHTTPBody(s)
	if err == nil {
		req, err = http.NewRequestWithContext(ctx, "POST", constants.EndpointFolderRename, body)
		if err == nil {
			req = readyHeader(req, s)
		}
	}

	return req, err
}

func FolderPaste(ctx context.Context, s *PremiumizeSession, r *FolderPasteRequest) (*http.Request, error) {
	var req *http.Request

	body, err := r.ToHTTPBody(s)
	if err == nil {
		req, err = http.NewRequestWithContext(ctx, "POST", constants.EndpointFolderPaste, body)
		if err == nil {
			req = readyHeader(req, s)
		}
	}

	return req, err
}

func FolderDelete(ctx context.Context, s *PremiumizeSession, r *FolderPasteRequest) (*http.Request, error) {
	var req *http.Request

	body, err := r.ToHTTPBody(s)
	if err == nil {
		req, err = http.NewRequestWithContext(ctx, "POST", constants.EndpointFolderDelete, body)
		if err == nil {
			req = readyHeader(req, s)
		}
	}

	return req, err
}

func FolderUploadInfo(ctx context.Context, s *PremiumizeSession, r *FolderUploadInfoRequest) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", constants.EndpointFolderUploadInfo, nil)
	if err == nil {
		req = readyHeader(req, s)
		req.URL.RawQuery = r.ToQuery(s)
	}

	return req, err
}

func FolderSearch(ctx context.Context, s *PremiumizeSession, r *FolderSearchRequest) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", constants.EndpointFolderSearch, nil)
	if err == nil {
		req = readyHeader(req, s)
		req.URL.RawQuery = r.ToQuery(s)
	}

	return req, err
}
