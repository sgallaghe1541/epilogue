package main

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sgallaghe1541/epilogue/internal/viewpoint"
	"github.com/sgallaghe1541/epilogue/views/layouts"
	"github.com/sgallaghe1541/epilogue/views/reports"
	"github.com/sgallaghe1541/epilogue/views/timeentry"
)

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

func (app *app) handleTimeCardLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Push-Url", r.URL.Path)
	timeentry.TimeEntryLinks().Render(context.Background(), w)
}
