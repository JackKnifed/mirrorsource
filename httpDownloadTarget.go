package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type httpDownload struct {
	Version      versionType `json:"version"`
	URLBase      string      `json:"urlBase"`
	Follow       bool        `json:"follow"`
	StoragePath  string      `json:"storagePath"`
	OutputFormat string      `json:"outputFormat"`
	Actions      []upstreamTarget
}

func (action httpDownload) RunVersion(s string) {
	fmtChunks := trimEmpty(strings.Split(action.Version.Fmt, "%v"))
	vPartsString, err := stringSubtraction(action.Version.Cur, fmtChunks)
	if err != nil {
		log.Printf("input passed [%q] does not match input format [%q]",
			s, action.Version.Fmt)
		return
	}
	vParts, err := stringsToInts(vPartsString)
	if err != nil {
		log.Println(err)
		return
	}

	outputPath := filepath.Join(action.StoragePath,
		fmt.Sprintf(action.OutputFormat, vParts))

	outfile, err := os.Create(outputPath)
	if err != nil {
		log.Printf("count not open dest file [%q] - %v", outputPath, err)
		return
	}
	defer outfile.Close()

	client := &http.Client{}
	if action.Follow {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	resp, err := client.Get(path.Join(action.URLBase, s))
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		if _, err := io.Copy(outfile, resp.Body); err != nil {
			log.Println(err)
		}
		for _, each := range action.Actions {
			each.RunVersion(s)
		}
	}
	return
}
