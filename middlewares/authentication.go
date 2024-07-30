package middlewares

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

func Auth(session *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userid := session.GetInt(r.Context(), "authenticatedUserID")
			if userid == 0 {
				http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			}
			next.ServeHTTP(w, r)
		})
	}
}
