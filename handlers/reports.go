package handlers

import (
	"context"
	"net/http"

	"github.com/sgallaghe1541/epilogue/views/reports"
)

func HandleReports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Push-Url", r.URL.Path)
	reports.ListReports().Render(context.Background(), w)
}
