package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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

func trimEmpty(parts []string) []string {
	if len(parts) > 1 && parts[0] == "" {
		parts = parts[1:]
	}
	if len(parts) > 1 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

func translateFormat(s, oldFmt, newFmt string) (string, error) {
	fmtChunks := trimEmpty(strings.Split(oldFmt, "%v"))
	vPartsString, err := stringSubtraction(s, fmtChunks)
	if err != nil {
		return "", errors.New("string does not match format")
	}
	return fmt.Sprintf(newFmt, vPartsString), nil
}
