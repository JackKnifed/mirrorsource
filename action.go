package mirrorsource

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Action interface {
	Do(*Version) error
}

type Sha1Verify struct {
	FileLoc    string // inherited
	FileFmt    string
	HashURLFmt string
}

func (a *Sha1Verify) Do(v *Version) error {
	filename := filepath.Join(a.FileLoc, v.Format(a.FileFmt))
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open [%s] for checksumming - %v", filename, err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("problem sha1'ing file - %s - %v", filename, err)
	}
	actualSum := fmt.Sprintf("%x", h.Sum(nil))

	resp, err := http.Get(v.Format(a.HashURLFmt))
	if err != nil {
		return fmt.Errorf("problem retrieving remote hash - %v", err)
	}
	defer resp.Body.Close()

	correctBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("problem retrieving remote hash - %v", err)
	}
	correctSum := strings.Fields(string(correctBytes))[0]

	if correctSum != actualSum {
		return fmt.Errorf("%s hash does not match", v.String())
	}
	return nil
}
