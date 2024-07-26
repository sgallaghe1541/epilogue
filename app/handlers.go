package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/package/viewpoint"
	"github.com/sgallaghe1541/epilogue/views/layouts"
	"github.com/sgallaghe1541/epilogue/views/reports"
	"github.com/sgallaghe1541/epilogue/views/timeentry"
	"github.com/xuri/excelize/v2"
)

func (a *App) HandleHome(w http.ResponseWriter, r *http.Request) {
	layouts.Base().Render(context.Background(), w)
}

func (a *App) HandleAllJobHours(w http.ResponseWriter, r *http.Request) {

	var vpArgs viewpoint.QueryArgs

	vpconn := a.Viewpoint

	v := r.URL.Query()

	div := v.Get("division")

	switch div {
	case "grading":
		vpArgs = viewpoint.Grading
	case "paving":
		vpArgs = viewpoint.Paving
	case "bridge":
		vpArgs = viewpoint.Bridge
	default:
		fmt.Println("division did not come through...")
		vpArgs = viewpoint.Grading
		//need to handle bad path
	}

	vpArgs.StartWEDate = v.Get("startwedate")
	vpArgs.EndWEDate = v.Get("endwedate")

	updatedURL := "/reports/alljobhours/?" + v.Encode()
	fmt.Println(updatedURL)

	jobHours := []*viewpoint.JobHoursResult{}
	query, args, err := viewpoint.BuildInQuery(viewpoint.JobHours, vpArgs)
	if err != nil {
		fmt.Print(err.Error())
	}

	rows, err := vpconn.Queryx(query, args...)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer rows.Close()

	err = sqlx.StructScan(rows, &jobHours)
	if err != nil {
		fmt.Print(err.Error())
	}

	for _, job := range jobHours {
		eerows, err := vpconn.Queryx(viewpoint.JobEmployeeHours, job.Job.String, vpArgs.StartWEDate, vpArgs.EndWEDate)
		if err != nil {
			fmt.Print(err.Error())
		}
		err = sqlx.StructScan(eerows, &job.EEHoursDetail)
		if err != nil {
			fmt.Print(err.Error())
		}
		eerows.Close()

		eqrows, err := vpconn.Queryx(viewpoint.JobEquipmentHours, job.Job.String, vpArgs.StartWEDate, vpArgs.EndWEDate)
		if err != nil {
			fmt.Print(err.Error())
		}
		err = sqlx.StructScan(eqrows, &job.EQHoursDetail)
		if err != nil {
			fmt.Print(err.Error())
		}
		eqrows.Close()
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Push-Url", updatedURL)
	reports.AllJobHours(jobHours).Render(context.Background(), w)
}

func (a *App) HandleReports(w http.ResponseWriter, r *http.Request) {
	conn := a.Epilogue
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

func (a *App) HandleReportParams(w http.ResponseWriter, r *http.Request) {
	reportQuery := "SELECT * FROM reports WHERE reporturl = ?"
	paramsQuery := "SELECT * FROM parameters WHERE reportid = ?"
	divisionQuery := "SELECT * FROM divisions"

	conn := a.Epilogue
	reporturl := chi.URLParam(r, "reportname")

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
			rows, err := conn.Queryx(divisionQuery)

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

func (a *App) HandleTimeCardLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Push-Url", r.URL.Path)
	timeentry.TimeEntryLinks().Render(context.Background(), w)
}

func (a *App) HandleEmployeesForFringe(w http.ResponseWriter, r *http.Request) {

	var vpArgs viewpoint.QueryArgs

	h := r.Header
	excel := h.Get("excel")

	vpconn := a.Viewpoint

	v := r.URL.Query()
	div := v.Get("division")

	switch div {
	case "grading":
		vpArgs = viewpoint.Grading
	case "paving":
		vpArgs = viewpoint.Paving
	case "bridge":
		vpArgs = viewpoint.Bridge
	default:
		fmt.Println("division did not come through...")
		vpArgs = viewpoint.Grading
		//need to handle bad path
	}

	updatedURL := "/reports/employeesforfringe/?" + v.Encode()

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
		w.Header().Set("HX-Redirect", filepath.Join(r.URL.Host, "downloads", filepath.Base(dir)))
		fmt.Println(filepath.Join(r.URL.Host, "downloads", filepath.Base(dir)))
	}
}

func (a *App) HandleDownloads(w http.ResponseWriter, r *http.Request) {
	fname := chi.URLParam(r, "fname")

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
