package handlers

import (
	"context"
	"net/http"

	"github.com/sgallaghe1541/epilogue/views/components"
	"github.com/sgallaghe1541/epilogue/views/reports"
)

func HandleEmployeeHours(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	if v.Get("wedate") == "" {
		date := getRecentWEDate()
		v.Set("wedate", date)
	}
	if v.Get("displayby") == "" {
		v.Set("displayby", "job")
	}
	vpDiv := v.Get("division")
	date := v.Get("wedate")
	displayBy := v.Get("displayby")

	w.Header().Set("Content-Type", "text/html")
	reports.EmployeeHours(vpDiv, displayBy, date).Render(context.Background(), w)
}

func HandleDisplayIn(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	components.HiddenDisplay(v.Get("display")).Render(context.Background(), w)
}
