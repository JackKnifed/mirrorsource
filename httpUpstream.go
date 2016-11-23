package main

import (
	"net/http"
	"path"
)

type httpResponseTracker struct {
	Version versionType `json:"version"`
	UrlBase string      `json:"urlBase"`
	Follow  bool        `json:"follow"`
	Actions []upstreamTarget{}
}

func (upstream httpResponseTracker) RunVersion(versionString string) error {
	client := &http.Client{}
	if upstream.Follow {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	resp, err := client.Head(path.Join(upstream.UrlBase, versionString))
	if resp.Code == http.StatusOK {
		upstream.Update(each)
	}
}
