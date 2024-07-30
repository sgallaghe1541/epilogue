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
	"github.com/sgallaghe1541/epilogue/middlewares"
)

func (app *App) Routes() *chi.Mux {
	mux := chi.NewRouter()
	fileDir := http.Dir("./static/")
	FileServer(mux, "/static", fileDir)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(app.SessionManager.LoadAndSave)
	r.Use(middleware.Logger)
	// fileDir := http.Dir("./static/")
	// FileServer(r, "/static", fileDir)

	//************AUTH**************
	r.Post("/signin", app.HandleSignin)
	r.Get("/signin", app.HandleSignin)
	r.Get("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")
		r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
		if _, err := gothic.CompleteUserAuth(w, r); err == nil {

			http.Redirect(w, r, "http://localhost:3000/", http.StatusFound)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
	})
	r.Get("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")
		r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
		gothUser, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(gothUser.RefreshToken)
		user, _ := app.Users.GetUser(gothUser.Email)
		// app.SessionManager.Put(r.Context(), "authenticatedUserID", user.ID)
		app.SessionManager.Put(r.Context(), "authenticatedUserName", user.Name)
		app.SessionManager.Put(r.Context(), "authenticatedUserDiv", user.Division)
		app.SessionManager.Put(r.Context(), "authenticatedUserPermission", user.PermissionLevel)
		http.Redirect(w, r, "http://localhost:3000/", http.StatusFound)
	})
	//************AUTH**************

	r.Group(func(r chi.Router) {
		r.Use(middlewares.Auth(app.SessionManager))
		r.Get("/", app.HandleHome)
	})

	reportRouter := chi.NewRouter()
	reportRouter.Get("/", app.HandleReports)
	reportRouter.Get("/{reportname}", app.HandleReportParams)
	reportRouter.Get("/{reportname}/downloads/{fname}", app.HandleDownloads)

	vprouter := chi.NewRouter()
	vprouter.Get("/alljobhours/", app.HandleAllJobHours)
	vprouter.Get("/employeesforfringe/", app.HandleEmployeesForFringe)

	r.Mount("/reports", reportRouter)
	r.Mount("/vp", vprouter)

	mux.Mount("/", r)
	return mux
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
