package handlers

import (
	"context"
	"net/http"

	"github.com/sgallaghe1541/epilogue/views/reports"
)

func HandleReports(w http.ResponseWriter, r *http.Request) {
	reportlist := map[string]string{
		"/jobhours/": "Job Hours",
	}
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Push-Url", r.URL.Path)
	reports.ListReports(reportlist).Render(context.Background(), w)
}

func HandleLoadHours(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()

	div := v.Get("division")
	if div == "" {
		div = "grading"
	}
	date := v.Get("wedate")
	updatedURL := "/reports/jobhours/?" + v.Encode()

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Push-Url", updatedURL)
	reports.JobHours(div, "job", date).Render(context.Background(), w)
}
