package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sgallaghe1541/epilogue/internal/db"
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
	timeentry.NewTimeCardHeader(db.TimeCardHeaderCreateForm{}).Render(context.Background(), w)
}

func (a *app) handlePostNewTimeHeader(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		a.serverError(w, r, err)
	}

	date, err := time.Parse("2006-01-02", r.PostForm.Get("workdate"))
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	currentWEDate := getWEDate(date)
	status := a.jobs.GetPayPeriodStatus(currentWEDate)

	form := db.TimeCardHeaderCreateForm{
		WorkDate:    date,
		Job:         r.PostForm.Get("job"),
		FieldErrors: map[string]string{},
	}

	if !form.WorkDate.Before(time.Now()) {
		form.FieldErrors["workdate"] = "Time card date cannot be in the future."
	} else if status == 1 {
		form.FieldErrors["workdate"] = fmt.Sprintf("The pay period containing %s is closed.", date.Format("01/02/2006"))
	}
	// validate job

	if len(form.FieldErrors) > 0 {
		for k, v := range form.FieldErrors {
			fmt.Printf("%s: %s\n", k, v)
		}
		timeentry.NewTimeCardHeader(form).Render(context.Background(), w)
		return
	}

	id, err := a.timecards.NewTimecard(form.Job, form.WorkDate, a.sessionManager.GetInt(r.Context(), "authenticatedUserID"))

	if err != nil {
		a.serverError(w, r, err)
		return
	}
	redirectURL := fmt.Sprintf("/timeentry/timecard/%d", id)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (a *app) handleTimeCard(w http.ResponseWriter, r *http.Request) {
	tcid, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		a.serverError(w, r, err)
		return
	}
	tcHeader, err := a.timecards.GetTimecardByID(tcid)
	if err != nil {
		a.serverError(w, r, err)
		return
	}
	tcEmployees, err := a.timecards.GetTimecardEmployees(tcid)
	if err != nil {
		a.serverError(w, r, err)
	}
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Push-Url", r.URL.Path)
	timeentry.EditTimeCard(tcHeader, tcEmployees, db.TimeCardEmployeesForm{}).Render(context.Background(), w)
}
