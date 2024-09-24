package middlewares

import "net/http"

func isAuthenticated(r *http.Request) bool {
	_, ok := r.Context().Value(IsAuthenticatedContextKey).(int)
	return ok
}
