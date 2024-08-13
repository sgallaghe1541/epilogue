package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.Handle("GET /{$}", app.authenticate(http.HandlerFunc(app.handleHome)))

	mux.HandleFunc("/signin", app.handleSignIn)
	mux.HandleFunc("GET /auth/microsoft", app.microsoftLogin)
	mux.HandleFunc("GET /auth/microsoft_callback", app.microsoftCallBack)

	mux.HandleFunc("GET /reports", app.handleReports)
	mux.HandleFunc("GET /reports/{reportname}", app.handleReportParams)

	mux.HandleFunc("GET /vp/alljobhours", app.handleAllJobHours)
	mux.HandleFunc("GET /vp/employeesforfringe/", app.HandleEmployeesForFringe)

	// reportRouter.Get("/{reportname}/downloads/{fname}", app.HandleDownloads)
	return commonHeaders(mux)
}
