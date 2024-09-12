package handlers

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/sgallaghe1541/epilogue/views/layouts"
)

func HandleHome(sessionManager *scs.SessionManager) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			name := sessionManager.GetString(r.Context(), "userName")
			layouts.Base(name).Render(context.Background(), w)
		})
}
