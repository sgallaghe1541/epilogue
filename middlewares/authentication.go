package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

func StartSession(store sessions.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "session-name")
			if err != nil {
				fmt.Println(err.Error())
			}
			err = session.Save(r, w)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(session)

			next.ServeHTTP(w, r)
		})
	}
}
