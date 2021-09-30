package api

import (
	"github.com/BRUHItsABunny/go-premiumize/constants"
	"net/http"
	"net/url"
)

func readyHeader(req *http.Request, session *PremiumizeSession) *http.Request {
	if req.Method == "POST" {
		req.Header["content-type"] = []string{constants.HeaderContentTypeForm}
	}
	req.Header["user-agent"] = []string{constants.HeaderUserAgent}
	req.Header = authHeader(session, req.Header)
	return req
}

func authHeader(session *PremiumizeSession, headers http.Header) http.Header {
	if session != nil && len(session.AuthToken) > 0 {
		if session.SessionType == "device_code" {
			headers["authorization"] = []string{"Bearer " + session.AuthToken}
		}
	}
	return headers
}

func authParams(session *PremiumizeSession, params url.Values) url.Values {
	if session != nil && len(session.AuthToken) > 0 {
		if session.SessionType == "apikey" {
			params["apikey"] = []string{session.AuthToken}
		}
	}
	return params
}
