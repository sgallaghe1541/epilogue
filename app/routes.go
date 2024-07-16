package app

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sgallaghe1541/epilogue/handlers"
	"github.com/sgallaghe1541/epilogue/middlewares"
)

func (app *App) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fileDir := http.Dir("./static/")
	FileServer(r, "/static", fileDir)

	r.Get("/", handlers.HandleHome)
	r.Get("/cmp/params/", handlers.HandleHiddenParams)

	reportRouter := chi.NewRouter()
	reportRouter.Use(middlewares.EpilogueConnection(app.Epilogue))
	reportRouter.Get("/", handlers.HandleReports)
	reportRouter.Get("/{reportname}", handlers.HandleReportParams)

	vprouter := chi.NewRouter()
	vprouter.Use(middlewares.VPConnection(app.Viewpoint))
	//vprouter.Get("/hoursbyjob/", handlers.HandleJobHours)
	vprouter.Get("/jobhours/", handlers.HandleHours)

	r.Mount("/reports", reportRouter)
	r.Mount("/vp", vprouter)

	return r
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
