package handlers

import (
	"context"
	"net/http"

	"github.com/sgallaghe1541/epilogue/internal/views/layouts"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	layouts.Landing().Render(context.Background(), w)
}
