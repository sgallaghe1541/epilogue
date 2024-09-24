package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/internal/timeentry"
)

func HandleJobsSelect(logger *slog.Logger, epilogue *db.EpilogueConnection) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			jobs, err := epilogue.GetJobsByDepartment("01")
			if err != nil {
				fmt.Println(err.Error())
			}
			w.Header().Set("Content-Type", "text/html")
			for _, job := range jobs {
				timeentry.SelectList(job).Render(context.Background(), w)
			}
		})
}

func HandleJobInfo(logger *slog.Logger, epilogue *db.EpilogueConnection) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vals := r.URL.Query()
			job, err := epilogue.GetJobByNumber(vals.Get("job"))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("HX-Trigger", "jobSelected")
			timeentry.JobStateCertified(&job).Render(context.Background(), w)
		})
}

func HandlePhaseSelect(logger *slog.Logger, epilogue *db.EpilogueConnection) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vals := r.URL.Query()
			phases, err := epilogue.GetPhasesByJob(vals.Get("job"))
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			w.Header().Set("Content-Type", "text/html")
			for _, phase := range phases {
				timeentry.SelectList(phase).Render(context.Background(), w)
			}
		})
}

func HandleClearPhases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	timeentry.Phases().Render(context.Background(), w)
}
