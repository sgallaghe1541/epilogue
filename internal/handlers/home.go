package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/internal/epicontext"
	"github.com/sgallaghe1541/epilogue/internal/ui/layouts"
)

func HandleHome(epilogue *db.EpilogueConnection, sessionManager *scs.SessionManager) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			user := sessionManager.GetInt(r.Context(), string(epicontext.IsAuthenticatedContextKey))
			name := sessionManager.GetString(r.Context(), "userName")
			role, err := epilogue.GetRoleByUserID(user)
			if err != nil {
				fmt.Println(err.Error())
			}
			icons, err := epilogue.GetIconsByRoleID(role)
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println("icons don't work")
				return
			}

			layouts.Base(name, icons).Render(context.Background(), w)
		})
}
