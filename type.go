package main

func main() {
	return
}

type upstreamTarget interface {
	RunVersion(string)
}

// these are probably junk?
// type outputType struct {
// 	active      bool   `json:"active"`
// 	gitRepo     string `json:"gitRepo"`
// 	location    string `json:"location"`
// 	outputID    int    `json:"outputID"`
// 	outputPrune int    `json:"outputPrune"`
// }

// type sourceConfig struct {
// 	name           string         "json:name"
// 	upstreamDriver string         "json:type"
// 	upstream       upstreamTarget "json:upstream"
// 	output         []outputType   "json:specFile"
// }
