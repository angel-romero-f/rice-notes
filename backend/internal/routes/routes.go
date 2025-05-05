// Package routes provides the router to set up handlers for the incoming HTTP requests.
package routes

import (
	"net/http"

	internal_middleware "github.com/angel-romero-f/rice-notes/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(internal_middleware.CORSMiddleware)

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Response"))
	})

	return r
}
