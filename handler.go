package main

import (
	"net/http"
	textTempl "text/template"
	"strconv"
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

func stringSubtraction(s string, chunks []string) ([]string, error) {
	var result []string
	remaining := strings.TrimLeft(chunks[0], s)

	for i := 1; i < len(chunks) ; i++{
		if !strings.Contains(remaining, chunks[i]) {
			return []string{}, errors.New(`chunk [%q] not found in [%q]
			original string [%q]
			cut chunks %#v
			parts gotten %#v`,
			chunks[i], remaining, s, chunks, result)
		}
		parts := strings.SplitN(remaining, chunks[i], 2)
		result = append(result, parts[0])
		remaining = parts[1]
	}

	if remaining != "" {
		result = append(result, remaining)
	}

	return result, nil
}

func stringsToInts(s []string) ([]int, error) {
	var results []int
	for _, s := range strings {
		i, err := strconv.Atoi(s)
		if err != nil {
			return []int{}, errors.New("failed - string [%q] is not an int - %v", s, err)
		}
		results = append(results, i)
	}
	return results
}

func (v versionType) nextPossible() (possible []string ) {
	fmtChunks := strings.Split(v.fmt, "%v")
	unprocessedPrefix := v.cur
	var processedSuffix []string

	for len(fmtChunks) > 0 {
		suffix := fmtChunks[len(fmtChunks)-1]
		prefix := strings.TrimSuffix(v.cur, fmtChunks[len(fmtChunks)-1])
		parts := string.Split(prefix, fmtChunks[len(fmtChunks)-2])
		versionNumber := parts[len(parts)-1]

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