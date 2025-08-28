// Package routes provides the router to set up handlers for the incoming HTTP requests.
package routes

import (
	"os"

	"github.com/angel-romero-f/rice-notes/internal/handlers"
	internal_middleware "github.com/angel-romero-f/rice-notes/internal/middleware"
	"github.com/angel-romero-f/rice-notes/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter sets up the routing and their handlers for incoming HTTP requests. Returns
// the router which main uses to start listsening for requests. 
func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(internal_middleware.CORSMiddleware)

	// Services and handlers to inject into respective methods
	noteService := services.NewNoteService()
	noteHandler := handlers.NewNoteHandler(noteService)

	// Auth setup with environment variables
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	// Create auth service and handler
	googleProvider := services.NewGoogleOAuth2Provider(googleClientID, googleClientSecret, redirectURL)
	authService := services.NewAuthService(googleProvider, jwtSecret)
	authHandler := handlers.NewAuthHandler(authService)

	// Routes setting up the handlers for incoming requests
	r.Get("/", noteHandler.CreateNote)

	// Auth routes
	r.Route("/api/auth", func(r chi.Router) {
		r.Get("/google", authHandler.GoogleLogin)
		r.Get("/google/callback", authHandler.GoogleCallback)
		r.Get("/me", authHandler.Me)
	})

	return r
}
