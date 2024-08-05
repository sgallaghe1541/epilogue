package handlers

import (
	"context"
	"net/http"

	"github.com/sgallaghe1541/epilogue/views/layouts"
)

func HandleHome() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			layouts.Base("failed").Render(context.Background(), w)
		})
}
