package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sgallaghe1541/epilogue/internal/viewpoint"
	"github.com/sgallaghe1541/epilogue/views/reports"
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
