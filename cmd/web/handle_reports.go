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
	"github.com/sgallaghe1541/epilogue/views/reports"
	"github.com/xuri/excelize/v2"
)

// import "net/http"

// func (a *app) handleReports(w http.ResponseWriter, r *http.Request) {
// 	reporturl := r.PathValue("reportname")
// 	v := r.URL.Query()

// 	if len(v) == 0 {
// 		updatedURL := r.URL.String()
// 		w.Header().Set("HX-Push-Url", updatedURL)
// 	}
// }

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
	divisionQuery := `SELECT report_permissions.divisionid AS divisionid, divisions.Description AS description 
	FROM report_permissions JOIN divisions 
	ON report_permissions.divisionid = divisions.divisionid 
	WHERE report_permissions.reportid = ?`

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

func (app *app) handleFHWAGroup(w http.ResponseWriter, r *http.Request) {
	vpargs := viewpoint.QueryArgs{}
	conn := app.viewpoint

	fhwa := &viewpoint.FHWAClassificationResult{}

	v := r.URL.Query()
	we := v.Get("wedate")
	if we == "" {
		app.serverError(w, r, fmt.Errorf("no wedate"))
	}
	vpargs.Date = we

	group := v.Get("prgroup")

	if group == "all" {
		vpargs.Groups = []string{"1", "2"}
	} else if group == "1" {
		vpargs.Groups = []string{"1"}
	} else if group == "2" {
		vpargs.Groups = []string{"2"}
	} else {
		vpargs.Groups = []string{"1", "2"}
		//app.serverError(w, r, fmt.Errorf("no groups"))
	}
	query, args, err := viewpoint.BuildInQuery(viewpoint.FHWAgroupQuery, vpargs)
	if err != nil {
		app.serverError(w, r, err)
	}
	rows, err := conn.Queryx(query, args...)

	if err != nil {
		fmt.Print(err.Error())
	}
	defer rows.Close()

	err = sqlx.StructScan(rows, &fhwa.Result)
	if err != nil {
		app.serverError(w, r, err)
	}

	w.Header().Set("Content-Type", "text/html")
	// w.Header().Set("HX-Push-Url", updatedURL)
	reports.FHWA(fhwa).Render(context.Background(), w)
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

	emps := &viewpoint.EmployeesForFringeResult{}
	query, args, err := viewpoint.BuildInQuery(viewpoint.EmployeesForFringeQuery, vpArgs)
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

func (app *app) handleFHWA(w http.ResponseWriter, r *http.Request) {
	args := viewpoint.QueryArgs{}
	conn := app.viewpoint

	fhwa := &viewpoint.FHWAClassificationResult{}

	v := r.URL.Query()
	we := v.Get("wedate")
	if we == "" {
		app.serverError(w, r, fmt.Errorf("no wedate"))
	}
	args.Date = we

	by := v.Get("displayby")

	if by == "job" {
		job := v.Get("job")
		if job == "" {
			app.serverError(w, r, fmt.Errorf("displayby=job but no job"))
		}
		args.Job = job

		rows, err := conn.NamedQuery(viewpoint.FHWAjobQuery, args)
		if err != nil {
			app.serverError(w, r, err)
		}
		defer rows.Close()

		err = sqlx.StructScan(rows, &fhwa.Result)
		if err != nil {
			app.serverError(w, r, err)
		}

	} else if by == "group" {
		group := v.Get("prgroup")

		if group == "all" {
			args.Groups = []string{"1", "2"}
		} else if group == "1" {
			args.Groups = []string{"1"}
		} else if group == "2" {
			args.Groups = []string{"2"}
		} else {
			app.serverError(w, r, fmt.Errorf("no groups"))
		}
		query, args, err := viewpoint.BuildInQuery(viewpoint.FHWAgroupQuery, args)
		if err != nil {
			app.serverError(w, r, err)
		}
		rows, err := conn.Queryx(query, args...)

		if err != nil {
			fmt.Print(err.Error())
		}
		defer rows.Close()

		err = sqlx.StructScan(rows, &fhwa.Result)
		if err != nil {
			app.serverError(w, r, err)
		}

	} else {
		app.serverError(w, r, fmt.Errorf("no displayby"))
	}

	w.Header().Set("Content-Type", "text/html")
	// w.Header().Set("HX-Push-Url", updatedURL)
	reports.FHWA(fhwa).Render(context.Background(), w)
}
