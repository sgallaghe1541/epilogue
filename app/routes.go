package app

import (
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sgallaghe1541/epilogue/handlers"
	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/middleware"
	"golang.org/x/oauth2"
)

func addRoutes(
	mux *http.ServeMux,
	auth *oauth2.Config,
	logger *slog.Logger,
	viewpoint *sqlx.DB,
	users *db.UserModel,
	refreshTokens *db.RefreshTokenModel,
) {
	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/signin", handlers.HandleSignIn(users, logger))
	mux.Handle("GET /auth/microsoft", handlers.MicrosoftLogin(auth))
	mux.Handle("GET /auth/microsoft_callback", handlers.MicrosoftCallBack(logger, auth, users, refreshTokens))

	mux.Handle("GET /{$}", middleware.Authenticate(logger, handlers.HandleHome()))

	// reportRouter := chi.NewRouter()
	// reportRouter.Get("/", app.HandleReports)
	// reportRouter.Get("/{reportname}", app.HandleReportParams)
	// reportRouter.Get("/{reportname}/downloads/{fname}", app.HandleDownloads)

	// vprouter := chi.NewRouter()
	// vprouter.Get("/alljobhours/", app.HandleAllJobHours)
	// vprouter.Get("/employeesforfringe/", app.HandleEmployeesForFringe)
}
