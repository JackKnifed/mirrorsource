package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type httpDownload struct {
	DisplayName      string
	Version          versionType `json:"version"`
	Follow           bool        `json:"follow"`
	StoragePath      string      `json:"storagePath"`
	URLBase          string      `json:"urlBase"`
	FileLocation     string
	FileFormat       string
	AddressFormat    string
	DownstreamFormat string
	UpstreamFormat   string
	Actions          []upstreamTarget
	localVersions    []httpDownloadHelper
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
		fmt.Sprintf(action.FileFormat, vParts))

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
			go each.RunVersion(s)
		}
	}
	return
}

type httpDownloadHelper []struct {
	Upstream   string
	Address    string
	Local      string
	Downstream string
}

func (h httpDownloadHelper) Len() int      { return len(h) }
func (h httpDownloadHelper) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h httpDownloadHelper) Less(i, j int) bool {
	return h[i].Address < h[j].Address
}

func (action httpDownload) listFiles() ([]httpDownloadHelper, error) {
	fileInfos, err := ioutil.ReadDir(action.FileLocation)
	if err != nil {
		return []httpDownloadHelper{}, fmt.Errorf("could not list location - %v", err)
	}

	var arr []httpDownloadHelper
	for _, f := range fileInfos {
		upstream, err := translateFormat(f.Name(),
			action.FileFormat, action.UpstreamFormat)
		if err != nil {
			continue
		}
		downstream, err := translateFormat(f.Name(),
			action.FileFormat, action.DownstreamFormat)
		if err != nil {
			continue
		}
		address, err := translateFormat(f.Name(),
			action.FileFormat, action.AddressFormat)
		if err != nil {
			continue
		}
		arr = append(arr,
			httpDownloadHelper{{Upstream: upstream, Downstream: downstream,
				Local: f.Name(), Address: address}})
	}
	sort.Sort(arr)
	return arr, nil
}

func (action httpDownload) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" || r.URL.Path == "." || r.URL.Path == "/" {
		// generate a list of all files stored on the server
		files, err := action.listFiles()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		action.localVersions = files
		err = templates.Lookup("packageList").Execute(w, action)
		if err != nil {
			log.Printf("problem listing files - %s - %v", r.URL.Path, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	file, err := os.Open(filepath.Join(action.FileLocation, r.URL.Path))
	switch {
	case err == os.ErrNotExist:
		log.Printf("client requested page that does not exist - %s", r.URL.Path)
		http.NotFound(w, r)
		return
	case err == os.ErrPermission:
		log.Printf("client requested page not authorized - %s", r.URL.Path)
		http.Error(w, "request not authorized", http.StatusForbidden)
		return
	case err != nil:
		log.Printf("unknown error opening path - %s - %v", r.URL.Path, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	default:
		fileInfo, err := file.Stat()
		if err != nil {
			log.Printf("failed to stat file - %s", r.URL.Path)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if fileInfo.IsDir() {
			log.Printf("client requested page is folder - %s", r.URL.Path)
			http.Error(w, "request not authorized", http.StatusForbidden)
			return
		}

		name, err := translateFormat(
			fileInfo.Name(), action.FileFormat, action.AddressFormat)
		if err != nil {
			log.Printf("file could not be translated to download name - %s - %v",
				r.URL.Path, err)
			http.Error(w, "could not rename download", http.StatusInternalServerError)
			return
		}
		http.ServeContent(w, r, name, fileInfo.ModTime(), file)
	}
}
