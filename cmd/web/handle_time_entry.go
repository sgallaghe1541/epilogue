package main

import (
	"context"
	"net/http"

	"github.com/sgallaghe1541/epilogue/views/timeentry"
)

func (a *app) handleTimeEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Push-Url", r.URL.Path)
	timeentry.TimeEntryLinks().Render(context.Background(), w)
}

func (a *app) handleNewTimeEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Push-Url", r.URL.Path)
	timeentry.NewTimeCard().Render(context.Background(), w)
}
