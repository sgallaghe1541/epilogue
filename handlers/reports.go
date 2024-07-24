package handlers

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
	"github.com/sgallaghe1541/epilogue/views/reports"
	"github.com/xuri/excelize/v2"
)

func HandleReports(w http.ResponseWriter, r *http.Request) {
	conn := r.Context().Value("epilogue").(*sqlx.DB)

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

func HandleReportParams(w http.ResponseWriter, r *http.Request) {
	reportQuery := "SELECT * FROM reports WHERE reporturl = ?"
	paramsQuery := "SELECT * FROM parameters WHERE reportid = ?"
	divisionQuery := "SELECT * FROM divisions"

	conn := r.Context().Value("epilogue").(*sqlx.DB)
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
	// w.Header().Set("HX-Push-Url", r.URL.Path)
	reports.Report(date, report, reportparams, divisions).Render(context.Background(), w)
}

// func HandleLoadHours(w http.ResponseWriter, r *http.Request) {
// 	v := r.URL.Query()

// 	div := v.Get("division")
// 	if div == "" {
// 		div = "grading"
// 	}
// 	date := v.Get("wedate")
// 	updatedURL := "/reports/jobhours/?" + v.Encode()

// 	w.Header().Set("Content-Type", "text/html")
// 	w.Header().Set("HX-Push-Url", updatedURL)
// 	reports.JobHours(div, "job", date).Render(context.Background(), w)
// }

func HandleEmployeesForFringe(w http.ResponseWriter, r *http.Request) {

	var vpArgs viewpoint.QueryArgs

	h := r.Header
	excel := h.Get("excel")

	vpconn := r.Context().Value("vp").(*sqlx.DB)

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
		//components.DownloadLink(filepath.Join(r.URL.Host, "downloads", filepath.Base(dir))).Render(context.Background(), w)
	}
}

func HandleDownloads(w http.ResponseWriter, r *http.Request) {
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
