package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sgallaghe1541/epilogue/views/timeentry"
)

func (a *app) handleJobsSelect(w http.ResponseWriter, r *http.Request) {
	jobs, err := a.jobs.GetJobsByDivision("%.01")
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Header().Set("Content-Type", "text/html")
	timeentry.JobSelect(jobs).Render(context.Background(), w)
}
