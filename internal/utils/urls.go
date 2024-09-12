package utils

import (
	"fmt"
	"net/url"
)

func GetAbsURL(u url.URL) string {
	b := u.JoinPath(u.Scheme, u.Host)
	return fmt.Sprint(b, "/")
}
