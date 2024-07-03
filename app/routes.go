package app

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/sgallaghe1541/epilogue/handlers"
	"github.com/sgallaghe1541/epilogue/middleware"
)

func (app *Epilogue) Routes() *chi.Mux {
	r := chi.NewRouter()

	fileDir := http.Dir("./static/")
	FileServer(r, "/static", fileDir)

	r.Get("/", handlers.HandleHome)
	r.Get("/reports/employeehours/*", handlers.HandleEmployeeHours)

	r.Get("/cmp/params/", handlers.HandleHiddenParams)

	r.Group(func(r chi.Router) {
		r.Use(middleware.VPConnection(app.Viewpoint))
		r.Get("/vp/hoursbyjob/", handlers.HandleJobHours)
		r.Get("/vp/hours/", handlers.HandleHours)
	})
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
