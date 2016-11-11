package main

import (
	"net/http"
	textTempl "text/template"
)

// type sourceConfig struct {
// 	name           string            "json:name"
// 	upstreamDriver string            "json:type"
// 	upstream       upstreamTarget "json:upstream"
// 	output         []outputType      "json:specFile"
// }

// type upstreamTarget interface {
// 	check() error
// 	version() string
// 	download() error
// }

type versionType struct {
	fmt  string   `json:"fmt"`
	cur  string   `json:"cur"`
	past []string `json:"past"`
}

func (v versionType) nextPossible() (possible []string ) {
	fmtChunks := strings.Split(v.fmt, "%v")
	remainingChunks := fmtChunks

	var possibleSubs []string
	for len(fmtChunks) > 0 {
		prefix := strings.TrimSuffix(v.cur, fmtChunks[len(fmtChunks)-1])

	}
}


type outputType struct {
	active      bool   `json:"active"`
	gitRepo     string `json:"gitRepo"`
	location    string `json:"location"`
	outputID    int    `json:"outputID"`
	outputPrune int    `json:"outputPrune"`
}

type httpResponseTracker struct {
	version versionType `json:"version"`
	driver string `json:"driver"`
	url    string `json:"url"`
}

// func (upstream httpResponseTracker) check() error {


// 	req, _ := http.NewRequest("GET", upstream.Replace(upstream.url), nil)
	



// }