package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *app) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	standardMiddle := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	sessionMiddle := alice.New(app.sessionManager.LoadAndSave)
	protectedMiddle := sessionMiddle.Append(app.authenticate, app.requireAuthentication)
	htmxMiddle := protectedMiddle.Append(app.htmxOnly)

	mux.Handle("GET /{$}", protectedMiddle.ThenFunc(app.handleHome))

	mux.HandleFunc("/signin", app.handleSignIn)
	mux.Handle("GET /auth/microsoft", sessionMiddle.ThenFunc(app.microsoftLogin))
	mux.Handle("GET /auth/microsoft_callback", sessionMiddle.ThenFunc(app.microsoftCallBack))
	mux.Handle("GET /logout", sessionMiddle.ThenFunc(app.handleLogOut))

	mux.Handle("GET /reports", protectedMiddle.ThenFunc(app.handleReports))
	mux.Handle("GET /reports/{reportname}", protectedMiddle.ThenFunc(app.handleReportParams))

	mux.Handle("GET /timeentry", sessionMiddle.ThenFunc(app.handleTimeEntry))
	mux.Handle("GET /timeentry/newtime", sessionMiddle.ThenFunc(app.handleNewTimeEntry))
	mux.Handle("GET /vp/alljobhours", htmxMiddle.ThenFunc(app.handleAllJobHours))
	mux.Handle("GET /vp/jobselect", htmxMiddle.ThenFunc(app.handleJobsSelect))
	mux.Handle("GET /vp/employeesforfringe", htmxMiddle.ThenFunc(app.handleEmployeesForFringe))
	mux.Handle("GET /vp/fhwabygroup", htmxMiddle.ThenFunc(app.handleFHWAGroup))

	mux.Handle("GET /downloads/{fname}", protectedMiddle.ThenFunc(app.handleDownloads))
	return standardMiddle.Then(mux)
}
