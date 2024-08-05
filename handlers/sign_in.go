package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sgallaghe1541/epilogue/auth"
	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/views/layouts"
	"golang.org/x/oauth2"
)

func HandleSignIn(u *db.UserModel, logger *slog.Logger) http.Handler {
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
				if u.ValidEmail(email) {
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

func MicrosoftLogin(config *oauth2.Config) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			url := config.AuthCodeURL(setState)
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		})
}

func MicrosoftCallBack(logger *slog.Logger, config *oauth2.Config, users *db.UserModel, tokenStore *db.RefreshTokenModel) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vals := r.URL.Query()
			state := vals.Get("state")
			if state != setState {
				logger.Info("States don't match")
			}

			code := vals.Get("code")

			token, err := config.Exchange(r.Context(), code)
			if err != nil {
				serverError(logger, w, r, err)
				return
			}

			userRequest, err := http.NewRequest("GET", endpointProfile, nil)
			if err != nil {
				serverError(logger, w, r, err)
				return
			}

			userRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

			client := http.DefaultClient

			userResponse, err := client.Do(userRequest)
			if err != nil {
				serverError(logger, w, r, err)
				return
			}

			defer userResponse.Body.Close()

			if userResponse.StatusCode != http.StatusOK {
				logger.Info("bad response", "httpcode", userResponse.StatusCode)
				return
			}

			email, err := readUserEmail(userResponse.Body)
			logger.Info(email)
			if err != nil {
				serverError(logger, w, r, err)
				return
			}

			eu, err := users.GetUser(email)
			if err != nil {
				logger.Info("no epilogue user found for", "email", email, "error", err.Error())
				clientError(w, http.StatusForbidden)
				return
			}

			jwt, err := auth.GenerateJWT(eu.ID, eu.Division, os.Getenv("JWT_SECRET"), time.Hour)
			if err != nil {
				serverError(logger, w, r, err)
				return
			}
			refresh, err := auth.GenerateRefreshToken()
			if err != nil {
				serverError(logger, w, r, err)
				return
			}
			err = tokenStore.SaveRefreshToken(eu.ID, refresh)
			if err != nil {
				serverError(logger, w, r, err)
				return
			}

			w.Header().Add("Set-Cookie", auth.NewJWTCookie(auth.JWTCookie, jwt))
			w.Header().Add("Set-Cookie", auth.NewRefreshCookie(auth.RefreshCookie, refresh))

			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		})
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
