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
	"github.com/sgallaghe1541/epilogue/package/viewpoint"
	"github.com/sgallaghe1541/epilogue/views/reports"
	"github.com/xuri/excelize/v2"
)

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
