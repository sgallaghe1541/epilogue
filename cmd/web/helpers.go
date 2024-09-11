package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (app *app) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
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

func readUserEmail(r io.Reader) (string, error) {
	u := struct {
		Name              string `json:"name"`
		Email             string `json:"mail"`
		FirstName         string `json:"givenName"`
		LastName          string `json:"surname"`
		NickName          string `json:"mailNickname"`
		UserPrincipalName string `json:"userPrincipalName"`
		Location          string `json:"usageLocation"`
	}{}

	err := json.NewDecoder(r).Decode(&u)
	if err != nil {
		return "", err
	}

	return strings.ToLower(u.Email), nil
}

func (app *app) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
