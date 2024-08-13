package main

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/sgallaghe1541/epilogue/internal/auth"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (app *app) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(auth.JWTCookie)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
				return
			default:
				app.logger.Error("server error", "err", err.Error())
				http.Redirect(w, r, "/signin", http.StatusInternalServerError)
				return
			}
		}

		username, userpermissions, err := auth.ValidateJWT(token.Value, os.Getenv("JWT_SECRET"))
		if err != nil {
			switch {
			case err.Error() == auth.JWTExpired:
				app.logger.Info(err.Error())
				//handle refresh
				http.Redirect(w, r, "/signin", http.StatusInternalServerError)
				return
			default:
				app.logger.Error("server error", "err", err.Error())
				http.Redirect(w, r, "/signin", http.StatusInternalServerError)
				return
			}
		}

		if userpermissions == nil {
			method := r.Method
			uri := r.URL.RequestURI()

			app.logger.Error("no user permissions found", "method", method, "uri", uri)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		r.Header.Set("user", username)
		for _, p := range userpermissions {
			r.Header.Add("perm", p)
		}
		next.ServeHTTP(w, r)
	})
}

func WEDate() func(http.Handler) http.Handler {
	weDate, _ := time.Parse("01/02/2006", "06/15/2024")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("wedate") == "" {
				q := r.URL.Query()
				q.Add("wedate", weDate.Format("01/02/2006"))
				r.URL.RawQuery = q.Encode()
				next.ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
