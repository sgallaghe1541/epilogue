package auth

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/azuread"
)

const (
	MaxAge = 86400 * 30
	IsProd = false
)

func NewAuth() {

	store := sessions.NewFilesystemStore("", []byte(os.Getenv("SESSION_SECRET")))
	store.MaxLength(8192)
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = true

	gothic.Store = store

	goth.UseProviders(
		azuread.New(os.Getenv("MICROSOFTONLINE_KEY"), os.Getenv("MICROSOFTONLINE_SECRET"), "http://localhost:3000/auth/azuread/callback", nil),
	)
}
