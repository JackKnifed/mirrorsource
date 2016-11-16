package main

func main() {
	return
}

type sourceConfig struct {
	name           string         "json:name"
	upstreamDriver string         "json:type"
	upstream       upstreamTarget "json:upstream"
	output         []outputType   "json:specFile"
}

type upstreamTarget interface {
	check() error
	version() string
	localVersions() []string
	download() error
}

type outputType struct {
	active      bool   `json:"active"`
	gitRepo     string `json:"gitRepo"`
	location    string `json:"location"`
	outputID    int    `json:"outputID"`
	outputPrune int    `json:"outputPrune"`
}
