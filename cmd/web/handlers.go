package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sgallaghe1541/epilogue/internal/auth"
	"github.com/sgallaghe1541/epilogue/views/layouts"
)

func (app *app) handleSignIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Info(err.Error())
	}
	email := strings.ToLower(r.FormValue("email"))
	switch email {
	case "":
		layouts.Signin().Render(context.Background(), w)
	default:
		if app.users.ValidEmail(email) {
			http.Redirect(w, r, "/auth/microsoft", http.StatusTemporaryRedirect)
		} else {
			layouts.Signin().Render(context.Background(), w)
		}
	}
}

const (
	setState        string = "oq^!yodE82F877#!^BpY7q#kG#"
	endpointProfile string = "https://graph.microsoft.com/v1.0/me"
)

func (app *app) microsoftLogin(w http.ResponseWriter, r *http.Request) {
	url := app.auth.AuthCodeURL(setState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *app) microsoftCallBack(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	state := vals.Get("state")
	if state != setState {
		app.logger.Info("States don't match")
	}

	code := vals.Get("code")

	token, err := app.auth.Exchange(r.Context(), code)
	if err != nil {
		serverError(app.logger, w, r, err)
		return
	}

	userRequest, err := http.NewRequest("GET", endpointProfile, nil)
	if err != nil {
		serverError(app.logger, w, r, err)
		return
	}

	userRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	client := http.DefaultClient

	userResponse, err := client.Do(userRequest)
	if err != nil {
		serverError(app.logger, w, r, err)
		return
	}

	defer userResponse.Body.Close()

	if userResponse.StatusCode != http.StatusOK {
		app.logger.Info("bad response", "httpcode", userResponse.StatusCode)
		return
	}

	email, err := readUserEmail(userResponse.Body)
	app.logger.Info(email)
	if err != nil {
		serverError(app.logger, w, r, err)
		return
	}

	eu, err := app.users.GetUser(email)
	if err != nil {
		app.logger.Info("no epilogue user found for", "email", email, "error", err.Error())
		clientError(w, http.StatusForbidden)
		return
	}

	jwt, err := auth.GenerateJWT(eu, os.Getenv("JWT_SECRET"), time.Hour)
	if err != nil {
		serverError(app.logger, w, r, err)
		return
	}
	refresh, err := auth.GenerateRefreshToken()
	if err != nil {
		serverError(app.logger, w, r, err)
		return
	}
	err = app.refreshTokens.SaveRefreshToken(eu.ID, refresh)
	if err != nil {
		serverError(app.logger, w, r, err)
		return
	}

	w.Header().Add("Set-Cookie", auth.NewJWTCookie(auth.JWTCookie, jwt))
	w.Header().Add("Set-Cookie", auth.NewRefreshCookie(auth.RefreshCookie, refresh))

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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

func (app *app) handleHome(w http.ResponseWriter, r *http.Request) {
	layouts.Base("failed").Render(context.Background(), w)
}
