package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/sgallaghe1541/epilogue/auth"
)

func Authenticate(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(auth.JWTCookie)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				http.Redirect(w, r, "/signin", http.StatusForbidden)
				return
			default:
				logger.Error("server error", "err", err.Error())
				http.Redirect(w, r, "/signin", http.StatusInternalServerError)
				return
			}
		}

		userid, userdiv, err := auth.ValidateJWT(token.Value, os.Getenv("JWT_SECRET"))
		if err != nil {
			switch {
			case err.Error() == auth.JWTExpired:
				logger.Info(err.Error())
				//handle refresh
				http.Redirect(w, r, "/signin", http.StatusInternalServerError)
				return
			default:
				logger.Error("server error", "err", err.Error())
				http.Redirect(w, r, "/signin", http.StatusInternalServerError)
				return
			}
		}

		r.Header.Set("userid", userid)
		r.Header.Set("userdiv", userdiv)
		next.ServeHTTP(w, r)
	})
}
