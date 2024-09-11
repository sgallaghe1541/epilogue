package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
		app.serverError(w, r, err)
		return
	}

	userRequest, err := http.NewRequest("GET", endpointProfile, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	userRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	client := http.DefaultClient

	userResponse, err := client.Do(userRequest)
	if err != nil {
		app.serverError(w, r, err)
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
		app.serverError(w, r, err)
		return
	}

	eu, err := app.users.GetUser(email)
	if err != nil {
		app.logger.Info("no epilogue user found for", "email", email, "error", err.Error())
		app.clientError(w, r, http.StatusForbidden, err)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", eu.ID)
	app.sessionManager.Put(r.Context(), "userName", eu.Name)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *app) handleLogOut(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "Log Out Successful")
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}
