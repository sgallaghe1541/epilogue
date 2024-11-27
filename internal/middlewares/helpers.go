package middlewares

import (
	"net/http"

	"github.com/sgallaghe1541/epilogue/internal/epicontext"
)

func isAuthenticated(r *http.Request) bool {
	_, ok := r.Context().Value(epicontext.IsAuthenticatedContextKey).(int)
	return ok
}
