package main

import (
	"net/http"
	"path"
)

type httpResponseTracker struct {
	Version versionType `json:"version"`
	UrlBase string      `json:"urlBase"`
	Follow  bool        `json:"follow"`
}

func (upstream httpResponseTracker) check() error {
	client := &http.Client{}
	if upstream.Follow {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}()
	}
	versions, err := upstream.possibleUpgrades()
	for _, each := range urls {
		resp, err := client.Head(path.Join(upstream.UrlBase, each))
		if resp.Code == http.StatusOK {
			upstream.Update(each)

		}
	}
}
