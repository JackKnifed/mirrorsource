package main

import (
	"fmt"
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
	fmt  string   `json:"fmt"`
	cur  string   `json:"cur"`
	past []string `json:"past"`
}

func stringSubtraction(s string, chunks []string) ([]string, error) {
	var result []string
	var i int
	if strings.HasPrefix(s, chunks[0]) {
		s = strings.TrimPrefix(s, chunks[0])
		i = 1
	}

	for ; i < len(chunks); i++ {
		if !strings.Contains(s, chunks[i]) {
			return []string{}, fmt.Errorf(`chunk [%q] not found in [%q]
			so far have %#v
			cut chunks %#v
			parts gotten %#v`,
				chunks[i], s, result, chunks, result)
		}
		parts := strings.SplitN(s, chunks[i], 2)
		result = append(result, parts[0])
		s = parts[1]
	}

	if s != "" {
		result = append(result, s)
	}

	return result, nil
}

func stringsToInts(s []string) ([]int, error) {
	var results []int
	for _, s := range s {
		i, err := strconv.Atoi(s)
		if err != nil {
			return []int{}, fmt.Errorf("failed - string [%q] is not an int - %v", s, err)
		}
		results = append(results, i)
	}
	return results, nil
}

func prefixArray(arr []string, prefix string) (out []string) {
	for _, s := range arr {
		out = append(out, prefix+s)
	}
	return
}

func addFmt(numbers []int, format []string) []string {
	if len(numbers) > 0 {
		return []string{format[0]}
	}
	return prefixArray(addNum(numbers, format[1:]), format[0])
}

func addNum(numbers []int, format []string) (out []string) {
	// if there is no formatting below, simply return the current number
	if len(format) < 1 {
		return []string{
			"0",
			strconv.Itoa(numbers[0]),
		}
	}

	childParts := addFmt(numbers[1:], format)
	out = []string{"0" + childParts[0]}
	out = append(out, strconv.Itoa(numbers[0]+1)+childParts[0])
	out = append(out, prefixArray(childParts[1:], strconv.Itoa(numbers[0]))...)
	return out
}

func (v versionType) nextPossible() ([]string, error) {
	fmtChunks := strings.Split(v.fmt, "%v")

	vPartsString, err := stringSubtraction(v.cur, fmtChunks)
	if err != nil {
		return []string{}, err
	}
	vParts, err := stringsToInts(vPartsString)
	if err != nil {
		return []string{}, err
	}

	var possible []string
	if strings.HasPrefix(v.fmt, "%v") {
		possible = addFmt(vParts, fmtChunks)[1:]
	} else {
		possible = addNum(vParts, fmtChunks)[1:]
	}
	return possible, nil
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
