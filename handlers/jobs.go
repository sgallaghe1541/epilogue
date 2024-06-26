package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sgallaghe1541/epilogue/package/viewpoint"
	"github.com/sgallaghe1541/epilogue/views/layouts"
)

func HandleAllJobHours(w http.ResponseWriter, r *http.Request) {
	division := r.URL.Path
	vpconn := r.Context().Value("vp").(*sqlx.DB)

	var vpDiv viewpoint.Division

	switch division {
	case "/grading":
		vpDiv = viewpoint.Grading
	default:
		//need to handle bad path
	}

	jobHours := []viewpoint.JobTotalHours{}
	query, args, err := viewpoint.BuildInQuery(viewpoint.JobHours, vpDiv)
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

	layouts.DivisionLanding(jobHours).Render(context.Background(), w)
}
