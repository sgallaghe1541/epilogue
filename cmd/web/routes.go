package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/sgallaghe1541/epilogue/internal/handlers"
	"github.com/sgallaghe1541/epilogue/internal/middlewares"
)

func (a *app) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	authenticate := middlewares.NewAuthMiddleware(a.epilogue, a.sessionManager, a.logger)
	log := middlewares.NewLoggerMiddleware(a.logger)
	recoverPanic := middlewares.NewRecoverPanicMiddleware(a.logger)

	standardMiddle := alice.New(recoverPanic, log, middlewares.CommonHeaders)
	sessionMiddle := alice.New(a.sessionManager.LoadAndSave)
	protectedMiddle := sessionMiddle.Append(authenticate, middlewares.RequireAuthentication)
	htmxMiddle := protectedMiddle.Append(middlewares.HtmxOnly)

	mux.Handle("GET /{$}", protectedMiddle.Then(handlers.HandleHome(a.sessionManager)))

	mux.Handle("/signin", handlers.HandleSignIn(a.logger, a.epilogue))
	mux.Handle("GET /auth/microsoft", sessionMiddle.Then(handlers.HandleMicrosoftLogin(a.auth)))
	mux.Handle("GET /auth/microsoft_callback", sessionMiddle.Then(handlers.HandleMicrosoftCallBack(a.auth, a.logger, a.sessionManager, a.epilogue)))
	mux.Handle("GET /logout", sessionMiddle.Then(handlers.HandleLogOut(a.logger, a.sessionManager)))

	// mux.Handle("GET /reports", protectedMiddle.ThenFunc(app.handleReports))
	// mux.Handle("GET /reports/{reportname}", protectedMiddle.ThenFunc(app.handleReportParams))

	// mux.Handle("GET /timeentry", sessionMiddle.ThenFunc(handlers.HandleTimeEntry))
	mux.Handle("GET /timeentry", htmxMiddle.Then(handlers.HandleTimeCards(a.logger, a.epilogue, a.sessionManager)))
	// mux.Handle("GET /timeentry/newtime", htmxMiddle.ThenFunc(handlers.HandleNewTimeCard))
	// mux.Handle("POST /timeentry/newtime", htmxMiddle.Then(handlers.HandlePostNewTimeHeader(a.logger, a.epilogue, a.viewpoint, a.sessionManager)))
	// mux.Handle("GET /timeentry/timecard/{id}", protectedMiddle.Then(handlers.HandleTimeCards(a.logger, a.epilogue)))
	// mux.Handle("GET /vp/alljobhours", htmxMiddle.ThenFunc(app.handleAllJobHours))
	mux.Handle("GET /tc/jobselect", htmxMiddle.Then(handlers.HandleJobsSelect(a.logger, a.epilogue)))
	mux.Handle("GET /tc/jobinfo", htmxMiddle.Then(handlers.HandleJobInfo(a.logger, a.epilogue)))
	mux.Handle("GET /tc/clearphases", htmxMiddle.ThenFunc(handlers.HandleClearPhases))
	mux.Handle("GET /tc/phaseselect", htmxMiddle.Then(handlers.HandlePhaseSelect(a.logger, a.epilogue)))
	mux.Handle("GET /tc/addemployeerow", htmxMiddle.Then(handlers.HandleAddEmployeeRow(a.logger)))
	mux.Handle("GET /tc/addequipmentrow", htmxMiddle.Then(handlers.HandleAddEquipmentRow(a.logger)))
	// mux.Handle("GET /vp/employeesforfringe", htmxMiddle.ThenFunc(app.handleEmployeesForFringe))
	// mux.Handle("GET /vp/fhwabygroup", htmxMiddle.ThenFunc(app.handleFHWAGroup))

	mux.Handle("GET /utils/weekday", htmxMiddle.Then(handlers.HandleWeekDay(a.logger)))

	// mux.Handle("GET /downloads/{fname}", protectedMiddle.ThenFunc(app.handleDownloads))
	return standardMiddle.Then(mux)
}
