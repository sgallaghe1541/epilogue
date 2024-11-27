package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/internal/ui/layouts"
	"github.com/sgallaghe1541/epilogue/internal/utils"
	"golang.org/x/oauth2"
)

func HandleSignIn(logger *slog.Logger, epilogue *db.EpilogueConnection) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				logger.Info(err.Error())
			}
			email := strings.ToLower(r.FormValue("email"))
			switch email {
			case "":
				layouts.Signin().Render(context.Background(), w)
			default:
				if epilogue.ValidEmail(email) {
					http.Redirect(w, r, "/auth/microsoft", http.StatusTemporaryRedirect)
				} else {
					layouts.Signin().Render(context.Background(), w)
				}
			}
		})
}

const (
	setState        string = "oq^!yodE82F877#!^BpY7q#kG#"
	endpointProfile string = "https://graph.microsoft.com/v1.0/me"
)

func HandleMicrosoftLogin(auth *oauth2.Config) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			url := auth.AuthCodeURL(setState)
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		})
}

func HandleMicrosoftCallBack(auth *oauth2.Config, logger *slog.Logger, sessionManager *scs.SessionManager, epilogue *db.EpilogueConnection) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vals := r.URL.Query()
			state := vals.Get("state")
			if state != setState {
				logger.Info("States don't match")
			}

			code := vals.Get("code")

			token, err := auth.Exchange(r.Context(), code)
			if err != nil {
				utils.ServerError(w, r, logger, err)
				return
			}

			userRequest, err := http.NewRequest("GET", endpointProfile, nil)
			if err != nil {
				utils.ServerError(w, r, logger, err)
				return
			}

			userRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

			client := http.DefaultClient

			userResponse, err := client.Do(userRequest)
			if err != nil {
				utils.ServerError(w, r, logger, err)
				return
			}

			defer userResponse.Body.Close()

			if userResponse.StatusCode != http.StatusOK {
				logger.Info("bad response", "httpcode", userResponse.StatusCode)
				return
			}

			email, err := utils.ReadUserEmail(userResponse.Body)
			logger.Info(email)
			if err != nil {
				utils.ServerError(w, r, logger, err)
				return
			}

			eu, err := epilogue.GetUser(email)
			if err != nil {
				logger.Info("no epilogue user found for", "email", email, "error", err.Error())
				utils.ClientError(w, http.StatusForbidden)
				return
			}

			err = sessionManager.RenewToken(r.Context())
			if err != nil {
				utils.ServerError(w, r, logger, err)
				return
			}

			sessionManager.Put(r.Context(), "authenticatedUserID", eu.ID)
			sessionManager.Put(r.Context(), "userName", eu.Name)

			http.Redirect(w, r, "/", http.StatusSeeOther)
		})
}

func HandleLogOut(logger *slog.Logger, sessionManager *scs.SessionManager) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := sessionManager.RenewToken(r.Context())
			if err != nil {
				utils.ServerError(w, r, logger, err)
				return
			}
			sessionManager.Remove(r.Context(), "authenticatedUserID")
			sessionManager.Put(r.Context(), "flash", "Log Out Successful")
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
		})
}
