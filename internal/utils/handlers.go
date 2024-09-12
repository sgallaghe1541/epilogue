package utils

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

func ServerError(w http.ResponseWriter, r *http.Request, logger *slog.Logger, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func ReadUserEmail(r io.Reader) (string, error) {
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
