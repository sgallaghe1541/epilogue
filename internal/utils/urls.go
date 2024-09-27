package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func GetAbsURL(u url.URL) string {
	b := u.JoinPath(u.Scheme, u.Host)
	return fmt.Sprint(b, "/")
}

func GetSelectNumDescription(s string) (string, string) {
	vals := strings.Split(s, " -- ")
	if len(vals) != 2 {
		return "", ""
	}
	return vals[0], vals[1]
}
