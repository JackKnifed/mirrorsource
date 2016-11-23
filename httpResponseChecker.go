package main

import (
	"log"
	"net/http"
	"path"
)

type httpResponseChecker struct {
	Version versionType `json:"version"`
	URLBase string      `json:"urlBase"`
	Follow  bool        `json:"follow"`
	Actions []upstreamTarget
}

func (action httpResponseChecker) RunVersion(s string) {
	client := &http.Client{}
	if action.Follow {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	resp, err := client.Head(path.Join(action.UrlBase, s))
	if err != nil {
		log.Println(err)
		return
	}
	if resp.StatusCode == http.StatusOK {
		for _, each := range action.Actions {
			each.RunVersion(s)
		}
	}
	return
}
