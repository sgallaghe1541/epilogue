package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/internal/middlewares"
	"github.com/sgallaghe1541/epilogue/internal/timeentry"
	"github.com/sgallaghe1541/epilogue/internal/utils"
)

// func HandleTimeEntry(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	w.Header().Set("HX-Push-Url", r.URL.Path)
// 	timeentry.TimeEntryLinks().Render(context.Background(), w)
// }

// func HandleNewTimeCard(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	w.Header().Set("HX-Push-Url", r.URL.Path)
// 	timeentry.ViewTimeCardHeader(timeentry.TimeCardHeaderForm{}).Render(context.Background(), w)
// }

// func HandlePostNewTimeHeader(logger *slog.Logger, epilogue *db.EpilogueConnection, viewpoint *viewpoint.ViewpointConnection, sessionManager *scs.SessionManager) http.Handler {
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			err := r.ParseForm()
// 			if err != nil {
// 				utils.ServerError(w, r, logger, err)
// 			}

// 			date, err := time.Parse("2006-01-02", r.PostForm.Get("workdate"))
// 			if err != nil {
// 				utils.ServerError(w, r, logger, err)
// 				return
// 			}

// 			currentWEDate := utils.GetWEDate(date)
// 			status := viewpoint.GetPayPeriodStatus(currentWEDate)

// 			jobNumber := r.PostForm.Get("job")
// 			job, err := epilogue.GetJobByNumber(jobNumber)

// 			form := timeentry.TimeCardHeaderForm{
// 				WorkDate:    date,
// 				Job:         &job,
// 				FieldErrors: map[string]string{},
// 			}

// 			if !form.WorkDate.Before(time.Now()) {
// 				form.FieldErrors["workdate"] = "Time card date cannot be in the future."
// 			} else if status == 1 {
// 				form.FieldErrors["workdate"] = fmt.Sprintf("The pay period containing %s is closed.", date.Format("01/02/2006"))
// 			}
// 			// validate job

// 			if len(form.FieldErrors) > 0 {
// 				for k, v := range form.FieldErrors {
// 					fmt.Printf("%s: %s\n", k, v)
// 				}
// 				timeentry.ViewTimeCardHeader(form).Render(context.Background(), w)
// 				return
// 			}

// 			id, err := epilogue.NewTimecard(form.Job, form.WorkDate, sessionManager.GetInt(r.Context(), "authenticatedUserID"))

// 			if err != nil {
// 				utils.ServerError(w, r, logger, err)
// 				return
// 			}
// 			redirectURL := fmt.Sprintf("/timeentry/timecard/%d", id)

// 			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
// 		})
// }

func HandleTimeCards(logger *slog.Logger, epilogue *db.EpilogueConnection, sessionManager *scs.SessionManager) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.PathValue("id") == "" {
				uid := r.Context().Value(middlewares.IsAuthenticatedContextKey).(int)
				tcid, err := epilogue.NewTimecard(uid)
				if err != nil {
					utils.ServerError(w, r, logger, err)
					return
				}
				updatedURL, _ := url.JoinPath(r.URL.Path, fmt.Sprintf("%d", tcid))
				w.Header().Set("Content-Type", "text/html")
				w.Header().Set("HX-Push-Url", updatedURL)

				timeentry.TimeCardForm(&timeentry.TimeCardHeaderForm{}, &timeentry.TimeCardDetailForm{}).Render(context.Background(), w)
				return
			}

			// tcid, err := strconv.Atoi(r.PathValue("id"))
			// if err != nil {
			// 	utils.ServerError(w, r, logger, err)
			// 	return
			// }
			// tcHeader, err := epilogue.GetTimecardByID(tcid)
			// if err != nil {
			// 	utils.ServerError(w, r, logger, err)
			// 	return
			// }
			// tcEmployees, err := epilogue.GetTimecardEmployees(tcid)
			// if err != nil {
			// 	utils.ServerError(w, r, logger, err)
			// }
			// w.Header().Set("Content-Type", "text/html")
			// w.Header().Set("HX-Push-Url", r.URL.Path)
			// timeentry.EditTimeCard(tcHeader, tcEmployees, timeentry.TimeCardDetailForm{}).Render(context.Background(), w)
		})
}

func HandleClearPhases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	timeentry.Phases(nil).Render(context.Background(), w)
}

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
			jobNum, jobDescription := utils.GetSelectNumDescription(vals.Get("job"))
			if jobNum == "" || jobDescription == "" {
				fmt.Println("failed to get job number")
				return
			}

			job, err := epilogue.GetJobByNumber(jobNum)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			w.Header().Set("Content-Type", "text/html")
			// w.Header().Set("HX-Trigger", "jobSelected")
			timeentry.JobStateCertified(&job).Render(context.Background(), w)
		})
}

func HandlePhaseSelect(logger *slog.Logger, epilogue *db.EpilogueConnection) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vals := r.URL.Query()
			jobNum, jobDescription := utils.GetSelectNumDescription(vals.Get("job"))
			if jobNum == "" || jobDescription == "" {
				fmt.Println("failed to get job number")
				return
			}
			phases, err := epilogue.GetPhasesByJob(jobNum)
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

func HandleEmployeeSelect(logger *slog.Logger, epilogue *db.EpilogueConnection) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			emps, err := epilogue.GetEmployeesByDepartment("01")
			if err != nil {
				logger.Error(err.Error())
				return
			}

			w.Header().Set("Content-Type", "text/html")
			for _, emp := range emps {
				timeentry.SelectList(emp).Render(context.Background(), w)
			}
		})
}

func HandleAddEmployeeRow(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vals := r.URL.Query()
			empCountString := vals.Get("employeerownumber")
			if empCountString == "" {
				logger.Error("no employeerownumber received...")
				return
			}
			empCount, err := strconv.Atoi(empCountString)
			if err != nil {
				logger.Error("malformed employeerownumber received...")
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("HX-Trigger", "employeeAdded")
			timeentry.EmployeeRow([]string{}, empCount+1).Render(context.Background(), w)
		})
}

func HandleEquipmentSelect(logger *slog.Logger, epilogue *db.EpilogueConnection) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			equipments, err := epilogue.GetEquipmentByDepartment("01")
			if err != nil {
				logger.Error(err.Error())
				return
			}
			w.Header().Set("Content-Type", "text/html")
			for _, equip := range equipments {
				timeentry.SelectList(equip).Render(context.Background(), w)
			}
		})
}

func HandleAddEquipmentRow(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vals := r.URL.Query()
			equipCountString := vals.Get("equipmentrownumber")
			if equipCountString == "" {
				logger.Error("no equipmentrownumber received...")
				return
			}
			equipCount, err := strconv.Atoi(equipCountString)
			if err != nil {
				logger.Error("malformed equipmentrownumber received...")
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("HX-Trigger", "equipmentAdded")
			timeentry.EquipmentRow([]string{}, equipCount+1).Render(context.Background(), w)
		})
}
