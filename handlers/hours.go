package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sgallaghe1541/epilogue/package/viewpoint"
	"github.com/sgallaghe1541/epilogue/views/components"
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

	updatedURL := "/landing/?" + v.Encode()
	fmt.Println(updatedURL)

	if v.Get("displayby") == "job" {
		jobHours := []viewpoint.JobTotalHours{}
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
		empHours := []viewpoint.EmployeeTotalHours{}
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
