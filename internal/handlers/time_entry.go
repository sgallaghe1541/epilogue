package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/internal/timeentry"
	"github.com/sgallaghe1541/epilogue/internal/utils"
	"github.com/sgallaghe1541/epilogue/internal/viewpoint"
)

func HandleTimeEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Push-Url", r.URL.Path)
	timeentry.TimeEntryLinks().Render(context.Background(), w)
}

func HandleNewTimeEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Push-Url", r.URL.Path)
	timeentry.ViewTimeCardHeader(timeentry.TimeCardHeaderCreateForm{}).Render(context.Background(), w)
}

func HandlePostNewTimeHeader(logger *slog.Logger, epilogue *db.EpilogueConnection, viewpoint *viewpoint.ViewpointConnection, sessionManager *scs.SessionManager) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				utils.ServerError(w, r, logger, err)
			}

			date, err := time.Parse("2006-01-02", r.PostForm.Get("workdate"))
			if err != nil {
				utils.ServerError(w, r, logger, err)
				return
			}

			currentWEDate := utils.GetWEDate(date)
			status := viewpoint.GetPayPeriodStatus(currentWEDate)

			form := timeentry.TimeCardHeaderCreateForm{
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
				timeentry.ViewTimeCardHeader(form).Render(context.Background(), w)
				return
			}

			id, err := epilogue.NewTimecard(form.Job, form.WorkDate, sessionManager.GetInt(r.Context(), "authenticatedUserID"))

			if err != nil {
				utils.ServerError(w, r, logger, err)
				return
			}
			redirectURL := fmt.Sprintf("/timeentry/timecard/%d", id)

			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		})
}

func HandleTimeCards(logger *slog.Logger, epilogue *db.EpilogueConnection) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			tcid, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				utils.ServerError(w, r, logger, err)
				return
			}
			tcHeader, err := epilogue.GetTimecardByID(tcid)
			if err != nil {
				utils.ServerError(w, r, logger, err)
				return
			}
			tcEmployees, err := epilogue.GetTimecardEmployees(tcid)
			if err != nil {
				utils.ServerError(w, r, logger, err)
			}
			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("HX-Push-Url", r.URL.Path)
			timeentry.EditTimeCard(tcHeader, tcEmployees, timeentry.TimeCardEmployeesForm{}).Render(context.Background(), w)
		})
}
