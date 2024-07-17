package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/views/reports"
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
	reports.ListReports(reportlist).Render(context.Background(), w)
}

func HandleReportParams(w http.ResponseWriter, r *http.Request) {
	reportQuery := "SELECT * FROM reports WHERE reporturl = ?"
	paramsQuery := "SELECT * FROM parameters WHERE reportid = ?"

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

	date := getRecentWEDate()

	w.Header().Set("Content-Type", "text/html")
	// w.Header().Set("HX-Push-Url", r.URL.Path)
	reports.Report(date, report, reportparams).Render(context.Background(), w)
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
