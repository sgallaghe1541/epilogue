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

func HandleNewTimeCard(logger *slog.Logger, epilogue *db.EpilogueConnection, sessionManager *scs.SessionManager) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.PathValue("id") == "" {
				uid := r.Context().Value(middlewares.IsAuthenticatedContextKey).(int)
				tcid, err := epilogue.NewTimecard(uid)
				if err != nil {
					utils.ServerError(w, r, logger, err)
					return
				}
				updatedURL, _ := url.JoinPath(r.URL.Path, "timecard", fmt.Sprintf("%d", tcid))
				w.Header().Set("Content-Type", "text/html")
				w.Header().Set("HX-Push-Url", updatedURL)

				timeentry.NewTimeCardForm().Render(context.Background(), w)
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

func HandlePutTimeCard(logger *slog.Logger, epilogue *db.EpilogueConnection, sessionManager *scs.SessionManager) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vals := r.URL.Query()
			fmt.Println("made it")
			for k, v := range vals {
				fmt.Printf("%s -- %s\n", k, v)
			}
			r.ParseForm()
			for k, v := range r.Form {
				fmt.Printf("%s -- %s\n", k, v)
			}
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
				w.Header().Set("Content-Type", "text/html")
				timeentry.JobStateCertified(nil).Render(context.Background(), w)
				return
			}

			job, err := epilogue.GetJobByNumber(jobNum)
			if err != nil {
				fmt.Println(err.Error())
				w.Header().Set("Content-Type", "text/html")
				timeentry.JobStateCertified(nil).Render(context.Background(), w)
				return
			}
			w.Header().Set("Content-Type", "text/html")
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
			phaseCountString := vals.Get("phasecount")
			if phaseCountString == "" {
				logger.Error("no phasecount received...")
				return
			}
			phaseCount, err := strconv.Atoi(phaseCountString)
			if err != nil {
				logger.Error("malformed phasecount received...")
				return
			}
			certified := vals.Get("certified")
			if certified != "N" && certified != "Y" {
				logger.Error("invalid certified value received...")
			}
			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("HX-Trigger", "employeeAdded")
			timeentry.EmployeeRow([]string{}, empCount+1, phaseCount, certified).Render(context.Background(), w)
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
			phaseCountString := vals.Get("phasecount")
			if equipCountString == "" {
				logger.Error("no phasecount received...")
				return
			}
			phaseCount, err := strconv.Atoi(phaseCountString)
			if err != nil {
				logger.Error("malformed phasecount received...")
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("HX-Trigger", "equipmentAdded")
			timeentry.EquipmentRow([]string{}, equipCount+1, phaseCount).Render(context.Background(), w)
		})
}

func HandleClassesSelect(logger *slog.Logger, epilogue *db.EpilogueConnection) http.Handler {
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
				logger.Error("invalid or closed job number...", "job", jobNum)
				return
			}
			classes, err := epilogue.GetCraftClassByJob(job)
			if err != nil {
				fmt.Println(err.Error())
			}
			w.Header().Set("Content-Type", "text/html")
			timeentry.BlankOption().Render(context.Background(), w)
			for _, class := range classes {
				timeentry.SelectList(class).Render(context.Background(), w)
			}
		})
}

func HandleEarnCodeSelect(logger *slog.Logger, epilogue *db.EpilogueConnection) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			earnCodes, err := epilogue.GetEarnCodes()
			if err != nil {
				logger.Error("could not get earncodes...", "err", err)
			}
			w.Header().Set("Content-Type", "text/html")
			timeentry.BlankOption().Render(context.Background(), w)
			for _, ec := range earnCodes {
				timeentry.SelectList(ec).Render(context.Background(), w)
			}
		})
}
