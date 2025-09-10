// Package routes provides the router to set up handlers for the incoming HTTP requests.
package routes

import (
	"context"
	"log/slog"
	"os"

	"github.com/angel-romero-f/rice-notes/internal/handlers"
	"github.com/angel-romero-f/rice-notes/internal/infra/storage"
	internal_middleware "github.com/angel-romero-f/rice-notes/internal/middleware"
	"github.com/angel-romero-f/rice-notes/internal/repository"
	"github.com/angel-romero-f/rice-notes/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RouterConfig contains configuration for setting up the router
type RouterConfig struct {
	DB           *pgxpool.Pool
	S3Bucket     string
	S3Region     string
	UseMockS3    bool // For development/testing
}

// NewRouter sets up the routing and their handlers for incoming HTTP requests. Returns
// the router which main uses to start listening for requests. 
func NewRouter(config *RouterConfig) (*chi.Mux, error) {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(internal_middleware.CORSMiddleware)

	// Auth setup with environment variables
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	// Create auth service and handler
	googleProvider := services.NewGoogleOAuth2Provider(googleClientID, googleClientSecret, redirectURL)
	authService := services.NewAuthService(googleProvider, jwtSecret)
	authHandler := handlers.NewAuthHandler(authService)

	// Create S3 uploader (or mock for development)
	var uploader storage.Uploader
	var err error

	if config.UseMockS3 {
		slog.Info("Using mock S3 uploader for development")
		uploader = storage.NewMockUploader()
	} else {
		slog.Info("Initializing S3 uploader", "bucket", config.S3Bucket, "region", config.S3Region)
		uploader, err = storage.NewS3Uploader(context.Background(), config.S3Bucket, config.S3Region)
		if err != nil {
			slog.Error("Failed to initialize S3 uploader", "error", err)
			return nil, err
		}
	}

	// Create repository layer
	noteRepo := repository.NewPostgresNoteRepository(config.DB)

	// Create services
	noteService := services.NewNoteService(noteRepo, uploader)

	// Create handlers  
	noteHandler := handlers.NewNoteHandler(noteService)

	// Public routes
	r.Get("/", noteHandler.Welcome)

	// Auth routes (public)
	r.Route("/api/auth", func(r chi.Router) {
		r.Get("/google", authHandler.GoogleLogin)
		r.Get("/google/callback", authHandler.GoogleCallback)
		r.Get("/me", authHandler.Me)
	})

	// Protected note routes (require JWT authentication)
	r.Route("/api/notes", func(r chi.Router) {
		// Apply JWT middleware to all routes in this group
		r.Use(internal_middleware.JWTMiddleware(authService))

		// Note endpoints
		r.Post("/", noteHandler.CreateNote)           // POST /api/notes - upload PDF
		r.Get("/", noteHandler.GetNotes)              // GET /api/notes - list user's notes
		r.Get("/{id}", noteHandler.GetNote)           // GET /api/notes/{id} - get specific note
		r.Delete("/{id}", noteHandler.DeleteNote)     // DELETE /api/notes/{id} - delete note
	})

	slog.Info("Router initialized successfully")
	return r, nil
}
