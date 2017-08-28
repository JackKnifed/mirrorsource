package mirrorsource

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type MD5Verify struct {
	errCh      chan<- error // inherited
	fileLoc    string       // inherited
	BaseFmt    string       // inherited
	HashURLFmt string
}

func (v *MD5Verify) Verify(basename string) bool {
	filename := filepath.Join(v.fileLoc, basename)
	f, err := os.Open(filename)
	if err != nil {
		v.errCh <- fmt.Errorf("failed to open file for checksumming - %v", err)
		return false
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		v.errCh <- fmt.Errorf("problem md5'ing file - %s - %v", filename, err)
		return false
	}
	actualSum := fmt.Sprintf("%x", h.Sum(nil))

	var version []interface{}
	_, err = fmt.Sscanf(basename, v.BaseFmt, version...)
	if err != nil {
		v.errCh <- fmt.Errorf("problem decoding version information - %v", err)
		return false
	}

	resp, err := http.Get(fmt.Sprintf(v.HashURLFmt, version...))
	if err != nil {
		v.errCh <- fmt.Errorf("problem retrieving remote hash - %v", err)
		return false
	}
	defer resp.Body.Close()

	correctBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		v.errCh <- fmt.Errorf("problem retrieving remote hash - %v", err)
		return false
	}
	correctSum := strings.Fields(string(correctBytes))[0]

	return correctSum == actualSum
}
