package main

type sourceConfig struct {
	name           string            "json:name"
	upstreamDriver string            "json:type"
	upstream       upstreamInterface "json:upstream"
	version        versionType       "json:version"
	output         []outputType      "json:specFile"
}

type upstreamTarget interface {
	func check()bool
	func download() error
	func build() error
}

type versionType struct {
	fmt  string   "json:fmt"
	cur  string   "json:cur"
	past []string "json:past"
}

type outputType struct {
	active      bool   "json:active"
	gitRepo     string "json:gitRepo"
	location    string "json:location"
	outputID    int    "json:outputID"
	outputPrune int    "json:outputPrune"
}

type httpResponseTracker struct {
	driver string "json:driver"
	url    string "json:url"
}

func ServeSource() {
	var configs []sourceConfig

}
