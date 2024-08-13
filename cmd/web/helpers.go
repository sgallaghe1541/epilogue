package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

func serverError(logger *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func getAbsURL(u url.URL) string {
	b := u.JoinPath(u.Scheme, u.Host)
	return fmt.Sprint(b, "/")
}

func getRecentWEDate() string {
	now := time.Now()

	now = now.Add(time.Hour * -24 * 5)

	for now.Weekday() != 6 {
		now = now.Add(time.Hour * -24)
	}

	return now.Format("2006-01-02")
}
