package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sgallaghe1541/epilogue/package/viewpoint"
	"github.com/sgallaghe1541/epilogue/views/components"
	"github.com/sgallaghe1541/epilogue/views/reports"
)

func HandleHours(w http.ResponseWriter, r *http.Request) {

	var vpArgs viewpoint.QueryArgs

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
		//need to handle bad path
	}

	vpArgs.WEDate = v.Get("wedate")

	updatedURL := "/reports/employeehours/?" + v.Encode()
	fmt.Println(updatedURL)

	if v.Get("displayby") == "job" {
		jobHours := []viewpoint.JobHoursResult{}
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

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("HX-Push-Url", updatedURL)
		components.JobHoursTable(jobHours).Render(context.Background(), w)
	}

	if v.Get("displayby") == "employee" {
		empHours := []viewpoint.EmployeeHoursResult{}
		query, args, err := viewpoint.BuildInQuery(viewpoint.EmployeeHours, vpArgs)
		if err != nil {
			fmt.Print(err.Error())
		}

		rows, err := vpconn.Queryx(query, args...)

		if err != nil {
			fmt.Print(err.Error())
		}

		defer rows.Close()

		err = sqlx.StructScan(rows, &empHours)
		if err != nil {
			fmt.Print(err.Error())
		}

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("HX-Push-Url", updatedURL)
		components.EmployeeHoursTable(empHours).Render(context.Background(), w)
	}
}

func HandleHiddenParams(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	for name, value := range v {
		components.HiddenDisplay(name, value[0]).Render(context.Background(), w)
	}
}

func HandleJobHours(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	reports.Report("Job Hours").Render(context.Background(), w)
}

// func HandleJobHours(w http.ResponseWriter, r *http.Request) {

// 	var vpArgs viewpoint.QueryArgs

// 	vpconn := r.Context().Value("vp").(*sqlx.DB)

// 	v := r.URL.Query()

// 	div := v.Get("division")

// 	switch div {
// 	case "grading":
// 		vpArgs = viewpoint.Grading
// 	case "paving":
// 		vpArgs = viewpoint.Paving
// 	case "bridge":
// 		vpArgs = viewpoint.Bridge
// 	default:
// 		fmt.Println("division did not come through...")
// 		vpArgs = viewpoint.Grading
// 		//need to handle bad path
// 	}

// 	vpArgs.WEDate = v.Get("wedate")
// 	vpArgs.Job = v.Get("job")

// 	empHours := []viewpoint.EmployeeHoursResult{}
// 	query, args, err := viewpoint.BuildQuery(viewpoint.JobEmployeeHours, vpArgs)
// 	if err != nil {
// 		fmt.Print(err.Error())
// 	}

// 	rows, err := vpconn.Queryx(query, args...)

// 	if err != nil {
// 		fmt.Print(err.Error())
// 	}

// 	defer rows.Close()

// 	err = sqlx.StructScan(rows, &empHours)
// 	if err != nil {
// 		fmt.Print(err.Error())
// 	}

// 	w.Header().Set("Content-Type", "text/html")
// 	components.EmployeeHoursRows(vpArgs.Job, empHours).Render(context.Background(), w)
// }
