package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/internal/viewpoint"
	"github.com/sgallaghe1541/epilogue/views/layouts"
	"github.com/sgallaghe1541/epilogue/views/reports"
	"github.com/sgallaghe1541/epilogue/views/timeentry"
	"github.com/xuri/excelize/v2"
)

func (app *app) handleSignIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Info(err.Error())
	}
	email := strings.ToLower(r.FormValue("email"))
	switch email {
	case "":
		layouts.Signin().Render(context.Background(), w)
	default:
		if app.users.ValidEmail(email) {
			http.Redirect(w, r, "/auth/microsoft", http.StatusTemporaryRedirect)
		} else {
			layouts.Signin().Render(context.Background(), w)
		}
	}
}

const (
	setState        string = "oq^!yodE82F877#!^BpY7q#kG#"
	endpointProfile string = "https://graph.microsoft.com/v1.0/me"
)

func (app *app) microsoftLogin(w http.ResponseWriter, r *http.Request) {
	url := app.auth.AuthCodeURL(setState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *app) microsoftCallBack(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	state := vals.Get("state")
	if state != setState {
		app.logger.Info("States don't match")
	}

	code := vals.Get("code")

	token, err := app.auth.Exchange(r.Context(), code)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	userRequest, err := http.NewRequest("GET", endpointProfile, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	userRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	client := http.DefaultClient

	userResponse, err := client.Do(userRequest)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	defer userResponse.Body.Close()

	if userResponse.StatusCode != http.StatusOK {
		app.logger.Info("bad response", "httpcode", userResponse.StatusCode)
		return
	}

	email, err := readUserEmail(userResponse.Body)
	app.logger.Info(email)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	eu, err := app.users.GetUser(email)
	if err != nil {
		app.logger.Info("no epilogue user found for", "email", email, "error", err.Error())
		clientError(w, http.StatusForbidden)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", eu.ID)
	app.sessionManager.Put(r.Context(), "userName", eu.Name)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *app) handleLogOut(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "Log Out Successful")
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

func (app *app) handleHome(w http.ResponseWriter, r *http.Request) {
	name := app.sessionManager.GetString(r.Context(), "userName")
	layouts.Base(name).Render(context.Background(), w)
}

func (app *app) handleAllJobHours(w http.ResponseWriter, r *http.Request) {

	var vpArgs viewpoint.QueryArgs

	vpconn := app.viewpoint
	excel := r.Header.Get("Excel")

	v := r.URL.Query()

	div := v.Get("division")

	switch div {
	case "01":
		vpArgs = viewpoint.Grading
	case "02":
		vpArgs = viewpoint.Paving
	case "06":
		vpArgs = viewpoint.Bridge
	case "30":
		vpArgs = viewpoint.Plants
	default:
		fmt.Println("division did not come through...")
		vpArgs = viewpoint.Grading
		//need to handle bad path
	}

	vpArgs.StartWEDate = v.Get("startwedate")
	vpArgs.EndWEDate = v.Get("endwedate")

	updatedURL := "/reports/alljobhours?" + v.Encode()

	jobHours := viewpoint.JobHoursResult{
		Result: []*viewpoint.JobHours{},
	}
	query, args, err := viewpoint.BuildInQuery(viewpoint.JobHoursQuery, vpArgs)
	if err != nil {
		fmt.Print(err.Error())
	}

	rows, err := vpconn.Queryx(query, args...)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer rows.Close()

	err = sqlx.StructScan(rows, &jobHours.Result)
	if err != nil {
		fmt.Print(err.Error())
	}

	for _, job := range jobHours.Result {
		job.GetDetail(vpArgs, vpconn)
	}

	if excel == "" {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("HX-Push-Url", updatedURL)
		reports.AllJobHours(jobHours).Render(context.Background(), w)
	} else {
		fName := fmt.Sprintf("AllJobHours-%s.xlsx", strings.Title(div))
		dir := filepath.Join(r.URL.Host, "tempfiles", fName)

		err := jobHours.ToExcel(dir)
		if err != nil {
			fmt.Println(err.Error())
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("HX-Redirect", fmt.Sprintf("/downloads/%s", filepath.Base(dir)))

	}
}

func (app *app) handleReports(w http.ResponseWriter, r *http.Request) {
	conn := app.epilogue
	reportlist := []db.EpilogueReport{}
	rows, err := conn.Queryx("SELECT * FROM reports")

	if err != nil {
		fmt.Print(err.Error())
	}

	defer rows.Close()

	err = sqlx.StructScan(rows, &reportlist)
	if err != nil {
		fmt.Print(err.Error())
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Push-Url", r.URL.Path)
	reports.ListReports(getAbsURL(*r.URL), reportlist).Render(context.Background(), w)
}

func (app *app) handleReportParams(w http.ResponseWriter, r *http.Request) {
	reportQuery := "SELECT * FROM reports WHERE reporturl = ?"
	paramsQuery := "SELECT * FROM parameters WHERE reportid = ?"
	divisionQuery := `SELECT report_divisions.divisionid AS divisionid, divisions.Description AS description 
	FROM report_divisions JOIN divisions 
	ON report_divisions.divisionid = divisions.divisionid 
	WHERE report_divisions.reportid = ?`

	conn := app.epilogue
	reporturl := r.PathValue("reportname")

	report := db.EpilogueReport{}
	err := conn.Get(&report, reportQuery, reporturl)
	if err != nil {
		fmt.Print(err.Error())
	}

	reportparams := []db.ReportParameter{}

	rows, err := conn.Queryx(paramsQuery, report.ID)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer rows.Close()

	err = sqlx.StructScan(rows, &reportparams)
	if err != nil {
		fmt.Print(err.Error())
	}

	date := ""
	for _, param := range reportparams {
		if param.Type == "date" {
			date = getRecentWEDate()
			break
		}
	}
	divisions := []db.Division{}
	for _, param := range reportparams {
		if param.Type == "division" {
			rows, err := conn.Queryx(divisionQuery, report.ID)

			if err != nil {
				fmt.Print(err.Error())
			}

			defer rows.Close()

			err = sqlx.StructScan(rows, &divisions)
			if err != nil {
				fmt.Print(err.Error())
			}
			break
		}
	}

	w.Header().Set("Content-Type", "text/html")
	reports.Report(date, report, reportparams, divisions).Render(context.Background(), w)
}

func (app *app) handleTimeCardLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Push-Url", r.URL.Path)
	timeentry.TimeEntryLinks().Render(context.Background(), w)
}

func (app *app) handleEmployeesForFringe(w http.ResponseWriter, r *http.Request) {

	var vpArgs viewpoint.QueryArgs

	h := r.Header
	excel := h.Get("excel")

	vpconn := app.viewpoint

	v := r.URL.Query()
	div := v.Get("division")

	switch div {
	case "01":
		vpArgs = viewpoint.Grading
	case "02":
		vpArgs = viewpoint.Paving
	case "06":
		vpArgs = viewpoint.Bridge
	default:
		fmt.Println("division did not come through...")
		vpArgs = viewpoint.Grading
		//need to handle bad path
	}

	updatedURL := "/reports/employeesforfringe?" + v.Encode()

	emps := viewpoint.EmployeesForFringeResult{}
	query, args, err := viewpoint.BuildInQuery(viewpoint.EmployeesForFringe, vpArgs)
	if err != nil {
		fmt.Print(err.Error())
	}

	rows, err := vpconn.Queryx(query, args...)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer rows.Close()

	err = sqlx.StructScan(rows, &emps.Result)
	if err != nil {
		fmt.Print(err.Error())
	}

	if excel == "" {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("HX-Push-Url", updatedURL)
		reports.EmpsForFringe(emps).Render(context.Background(), w)
	} else {
		fName := fmt.Sprintf("EmployeesForFringe-%s.xlsx", strings.Title(div))
		dir := filepath.Join(r.URL.Host, "tempfiles", fName)

		err := emps.ToExcel(dir)
		if err != nil {
			fmt.Println(err.Error())
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("HX-Redirect", fmt.Sprintf("/downloads/%s", filepath.Base(dir)))

	}
}

func (app *app) handleDownloads(w http.ResponseWriter, r *http.Request) {
	fname := r.PathValue("fname")

	downloadFile := filepath.Join(r.URL.Host, "tempfiles", fname)

	f, err := excelize.OpenFile(downloadFile)
	if err != nil {
		fmt.Println(err.Error())
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		fmt.Println(err.Error())
	}

	http.ServeContent(w, r, fname, time.Time{}, strings.NewReader(buf.String()))
	defer func() {
		err := os.Remove(downloadFile)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
}
