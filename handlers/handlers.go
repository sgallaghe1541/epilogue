package handlers

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sgallaghe1541/epilogue/package/viewpoint"
	"github.com/sgallaghe1541/epilogue/views/layouts"
)

func DivisionJobs(w http.ResponseWriter, r *http.Request) {
	division := r.URL.Path
	vpconn := r.Context().Value("vp").(*sqlx.DB)

	var vpDiv viewpoint.Division

	switch division {
	case "/grading":
		vpDiv = viewpoint.Grading
	default:
		//need to handle bad path
	}

	divisionJobs := []viewpoint.Job{}
	query, args, err := viewpoint.BuildInQuery(viewpoint.JobListQuery, vpDiv)
	if err != nil {
	}

	rows, err := vpconn.Queryx(query, args...)

	if err != nil {
	}

	defer rows.Close()

	err = sqlx.StructScan(rows, &divisionJobs)
	if err != nil {
	}

	layouts.JobList(divisionJobs).Render(context.Background(), w)
}
