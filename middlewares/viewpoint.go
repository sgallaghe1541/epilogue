package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

func VPConnection(vp *sqlx.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "vp", vp)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
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
