package handlers

import (
	"context"
	"net/http"

	"github.com/sgallaghe1541/epilogue/views/layouts"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	layouts.Home().Render(context.Background(), w)
}
