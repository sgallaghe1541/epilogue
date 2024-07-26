package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/markbates/goth/gothic"
	"github.com/sgallaghe1541/epilogue/handlers"
	"github.com/sgallaghe1541/epilogue/middlewares"
)

func (app *App) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Use(middlewares.StartSession(gothic.Store))
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	fileDir := http.Dir("./static/")
	FileServer(r, "/static", fileDir)

	r.Get("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")
		r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
		if user, err := gothic.CompleteUserAuth(w, r); err == nil {
			fmt.Println(user)
			http.Redirect(w, r, "http://localhost:3000/", http.StatusFound)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
	})
	r.Get("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Made it")
		provider := chi.URLParam(r, "provider")
		r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(user)
		http.Redirect(w, r, "http://localhost:3000/", http.StatusFound)
	})
	r.Get("/", handlers.HandleHome)

	reportRouter := chi.NewRouter()
	reportRouter.Use(middlewares.EpilogueConnection(app.Epilogue))
	reportRouter.Get("/", handlers.HandleReports)
	reportRouter.Get("/{reportname}", handlers.HandleReportParams)
	reportRouter.Get("/{reportname}/downloads/{fname}", handlers.HandleDownloads)

	vprouter := chi.NewRouter()
	vprouter.Use(middlewares.VPConnection(app.Viewpoint))
	//vprouter.Get("/hoursbyjob/", handlers.HandleJobHours)
	vprouter.Get("/alljobhours/", handlers.HandleAllJobHours)
	vprouter.Get("/employeesforfringe/", handlers.HandleEmployeesForFringe)

	r.Mount("/reports", reportRouter)
	r.Mount("/vp", vprouter)

	// r.Get("/downloads/{fname}", handlers.HandleDownloads)

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
