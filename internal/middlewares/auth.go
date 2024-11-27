package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/internal/epicontext"
	"github.com/sgallaghe1541/epilogue/internal/utils"
)

func RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAuthenticated(r) {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func NewAuthMiddleware(epilogue *db.EpilogueConnection, sessionManager *scs.SessionManager, logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := sessionManager.GetInt(r.Context(), "authenticatedUserID")
			if id == 0 {
				next.ServeHTTP(w, r)
			}
			active, err := epilogue.ActiveUser(id)
			if err != nil {
				utils.ServerError(w, r, logger, err)
				return
			}
			if active {
				ctx := context.WithValue(r.Context(), epicontext.IsAuthenticatedContextKey, id)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}
