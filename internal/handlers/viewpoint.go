package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/internal/timeentry"
	"github.com/sgallaghe1541/epilogue/internal/viewpoint"
)

func HandleJobsSelect(logger *slog.Logger, epilogue *db.EpilogueConnection) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			jobs, err := epilogue.GetJobsByDivision("%.01")
			if err != nil {
				fmt.Println(err.Error())
			}
			w.Header().Set("Content-Type", "text/html")
			for _, job := range jobs {
				timeentry.SelectList(job).Render(context.Background(), w)
			}
		})
}

func HandlePhaseSelect(logger *slog.Logger, viewpoint *viewpoint.ViewpointConnection) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

		})
}
