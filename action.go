package mirrorsource

import (
	"crypto/md5"
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
	Do(Version) error
}

type Sha1Verify struct {
	FileLoc    string // inherited
	FileFmt    string
	HashURLFmt string
}

func (a *Sha1Verify) Do(v Version) error {
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

type Md5Verify struct {
	FileLoc    string // inherited
	FileFmt    string
	HashURLFmt string
}

func (a *Md5Verify) Do(v Version) error {
	filename := filepath.Join(a.FileLoc, v.Format(a.FileFmt))
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open [%s] for checksumming - %v",
			filename, err)
	}
	defer f.Close()

	h := md5.New()
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

type CheckURL struct {
	URLFmt string
}

func (a *CheckURL) Do(v Version) error {
	resp, err := http.DefaultClient.Head(v.Format(a.URLFmt))
	if err != nil {
		return fmt.Errorf("problem checking %s - %v", v.Format(a.URLFmt), err)
	}
	// TODO: need some better checks to see what to do about sites that redirect newer stuff to the latest
	if resp.StatusCode != http.StatusFound {
		return fmt.Errorf("%s not found", v.String())
	}
	return nil
}

type GetURL struct {
	URLFmt string
	Output io.WriteCloser
}

func (a *GetURL) Do(v Version) error {
	resp, err := http.DefaultClient.Get(v.Format(a.URLFmt))
	if err != nil {
		return fmt.Errorf("problem retrieving %s - %v",
			v.Format(a.URLFmt), err)
	}
	if resp.StatusCode != http.StatusFound {
		resp.Body.Close()
		return fmt.Errorf("problem retrieving %s - %d status",
			v.String(), resp.StatusCode)
	}

	go func() {
		defer a.Output.Close()
		defer resp.Body.Close()
		io.Copy(a.Output, resp.Body)
	}()

	return nil
}

type SaveFile struct {
	FileLoc string
	FileFmt string
	Perm    os.FileMode
	In      io.ReadCloser
}

func (a *SaveFile) Do(v Version) error {
	defer a.In.Close()
	f, err := os.OpenFile(filepath.Join(a.FileLoc, v.Format(a.FileFmt)),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, a.Perm)
	if err != nil {
		return fmt.Errorf("problem opening file %s - %v",
			filepath.Join(a.FileLoc, v.Format(a.FileFmt)), err)
	}
	defer f.Close()
	v.AddRevertAction(&RemoveFile{FileLoc: a.FileLoc, FileFmt: a.FileFmt})
	_, err = io.Copy(f, a.In)
	return err
}

type RemoveFile struct {
	FileLoc string
	FileFmt string
	In      io.ReadCloser
}

func (a *RemoveFile) Do(v Version) error {
	err := os.Remove(filepath.Join(a.FileLoc, v.Format(a.FileFmt)))
	return err
}

type NextVersions struct{}

func (a *NextVersions) Do(v Version) error {
	err := v.NextVersions()
	if err != nil {
		return err
	}
	return nil
}
