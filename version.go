package main

import (
	"strconv"
	"strings"
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

func main() {
	return
}

type versionType struct {
	Fmt   string   `json:"fmt"`
	Cur   string   `json:"cur"`
	Local []string `json:"local"`
}

func addFmt(numbers []int, format []string) []string {
	switch {
	case len(numbers) == 0:
		return []string{format[0]}
	case len(format) < 2:
		return prefixArray(addNum(numbers, []string{}), format[0])
	default:
		return prefixArray(addNum(numbers, format[1:]), format[0])
	}
}

func addNum(numbers []int, format []string) (out []string) {
	// if there is no formatting below, simply return the current number
	var childParts []string
	switch {
	case len(format) == 0:
		return []string{
			"0",
			strconv.Itoa(numbers[0] + 1),
		}
	case len(numbers) < 2:
		childParts = addFmt([]int{}, format)
	default:
		childParts = addFmt(numbers[1:], format)
	}

	out = []string{"0" + childParts[0]}
	out = append(out, strconv.Itoa(numbers[0]+1)+childParts[0])
	out = append(out, prefixArray(childParts[1:], strconv.Itoa(numbers[0]))...)
	return out
}

func (v versionType) possibleUpgrades() ([]string, error) {
	fmtChunks := trimEmpty(strings.Split(v.Fmt, "%v"))

	vPartsString, err := stringSubtraction(v.Cur, fmtChunks)
	if err != nil {
		return []string{}, err
	}
	vParts, err := stringsToInts(vPartsString)
	if err != nil {
		return []string{}, err
	}

	var possible []string
	if strings.HasPrefix(v.Fmt, "%v") {
		possible = addNum(vParts, fmtChunks)
	} else {
		possible = addFmt(vParts, fmtChunks)
	}
	return possible[1:], nil
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
	driver  string      `json:"driver"`
	url     string      `json:"url"`
}

// func (upstream httpResponseTracker) check() error {

// 	req, _ := http.NewRequest("GET", upstream.Replace(upstream.url), nil)

// }
