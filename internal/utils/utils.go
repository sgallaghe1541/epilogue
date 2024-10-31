package utils

import (
	"fmt"
	"strings"
)

func Split(s string) (string, string, error) {
	split := strings.Split(s, " -- ")
	if len(split) != 2 {
		return "", "", fmt.Errorf("could not split %s", s)
	}
	return split[0], split[1], nil
}

func SortIntKeys(keys []int) []int {
	for i := 0; i < len(keys)-1; i++ {
		for j := 0; j < len(keys)-i-1; j++ {
			if keys[j] > keys[j+1] {
				keys[j], keys[j+1] = keys[j+1], keys[j]
			}
		}
	}
	return keys
}

func SortIntKeysS(m map[int]string) []int {
	keys := []int{}

	for k := range m {
		keys = append(keys, k)
	}

	for i := 0; i < len(keys)-1; i++ {
		for j := 0; j < len(keys)-i-1; j++ {
			if keys[j] > keys[j+1] {
				keys[j], keys[j+1] = keys[j+1], keys[j]
			}
		}
	}
	return keys
}
