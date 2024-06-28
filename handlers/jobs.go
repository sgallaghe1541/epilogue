package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sgallaghe1541/epilogue/package/viewpoint"
	"github.com/sgallaghe1541/epilogue/views/layouts"
)

func HandleAllJobHours(w http.ResponseWriter, r *http.Request) {
	division := r.URL.Path
	vpconn := r.Context().Value("vp").(*sqlx.DB)

	var vpDiv viewpoint.QueryArgs

	switch division {
	case "/vp/grading":
		vpDiv = viewpoint.Grading
	case "/vp/paving":
		vpDiv = viewpoint.Paving
	case "/vp/bridge":
		vpDiv = viewpoint.Bridge
	default:
		//need to handle bad path
	}

	vpDiv.WEDate = r.URL.Query().Get("wedate")
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

	dates := getDates(vpDiv.WEDate)

	w.Header().Set("Content-Type", "text/html")
	layouts.DivisionLanding(dates, jobHours).Render(context.Background(), w)
}

func getDates(d string) [3]string {
	var dates [3]string

	datetime, _ := time.Parse("01/02/2006", d)
	prev := datetime.AddDate(0, 0, -7)
	next := datetime.AddDate(0, 0, 7)

	dates[0] = prev.Format("01/02/2006")
	dates[1] = d
	dates[2] = next.Format("01/02/2006")

	return dates
}
