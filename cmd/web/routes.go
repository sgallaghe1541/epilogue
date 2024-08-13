package main

import (
	"net/http"
)

func (app *app) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/signin", app.handleSignIn)
	mux.HandleFunc("GET /auth/microsoft", app.microsoftLogin)
	mux.HandleFunc("GET /auth/microsoft_callback", app.microsoftCallBack)

	mux.Handle("GET /{$}", app.authenticate(http.HandlerFunc(app.handleHome)))

	// reportRouter := chi.NewRouter()
	// reportRouter.Get("/", app.HandleReports)
	// reportRouter.Get("/{reportname}", app.HandleReportParams)
	// reportRouter.Get("/{reportname}/downloads/{fname}", app.HandleDownloads)

	// vprouter := chi.NewRouter()
	// vprouter.Get("/alljobhours/", app.HandleAllJobHours)
	// vprouter.Get("/employeesforfringe/", app.HandleEmployeesForFringe)
	return mux
}
