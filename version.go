package mirrorsource

import (
	"errors"
	"fmt"
	"sync"
	"unicode/utf8"
)

type Version struct {
	lock sync.RWMutex
	fmt  string
	val  []interface{}
}

func DecodeVersion(format string, encoded string) (*Version, error) {
	retVal := &Version{fmt: format}
	_, err := fmt.Sscanf(retVal.fmt, encoded, retVal.val...)
	if err != nil {
		return nil, err
	}
	return retVal, nil
}

func (v *Version) String() string {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return fmt.Sprintf(v.fmt, v.val...)
}

func (v *Version) Format(f string) string {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return fmt.Sprintf(f, v.val...)
}

func incrementInterface(in interface{}) (interface{}, error) {
	switch val := in.(type) {
	case bool:
		// if it's a bool, there is no carry over
		if !val {
			return true, nil
		}
	case int:
		return val + 1, nil
	case uint:
		return val + 1, nil
	case string:
		r, _ := utf8.DecodeLastRuneInString(val)
		return fmt.Sprintf("%s%s", val[:len(val)-1], string(r+1)), nil
	}
	return nil, errors.New("was used on a value that cannot be incremented")
}

func resetInterface(in interface{}) (interface{}, error) {
	switch in.(type) {
	case bool:
		return false, nil
	case int:
		return 0, nil
	case uint:
		return 0, nil
	case string:
		return "a", nil
	}
	return nil, errors.New("was used on a value that cannot be incremented")
}

// not sure if this should be a method of versions or if it should be it's own function
func IncrementVersion(v *Version) ([]*Version, error) {
	checkVer := v.val[:]
	nextVers := []*Version{}
	var err error

	for i := len(checkVer) - 1; i >= 0; i-- {
		checkVer[i], err = incrementInterface(checkVer[i])
		if err != nil {
			return nil, err
		}
		// add that version to the return stuff
		nextVers = append(nextVers, &Version{fmt: v.fmt, val: checkVer})
		// reset it to it's deafut for the next run
		checkVer[i], err = resetInterface(checkVer[i])
		if err != nil {
			return nil, err
		}
	}

	return nextVers, nil
}
