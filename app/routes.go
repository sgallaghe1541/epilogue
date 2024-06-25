package app

import (
	"github.com/go-chi/chi"
	"github.com/sgallaghe1541/epilogue/handlers"
	"github.com/sgallaghe1541/epilogue/middleware"
)

func (app *Epilogue) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", handlers.HandleHome)

	r.Group(func(r chi.Router) {
		r.Use(middleware.VPConnection(app.Viewpoint))
		r.Get("/grading", handlers.DivisionJobs)
	})
	return r
}
