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
