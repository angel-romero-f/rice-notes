package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/angel-romero-f/rice-notes/internal/routes"
	"github.com/angel-romero-f/rice-notes/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .local.env file
	if err := godotenv.Load(".local.env"); err != nil {
		log.Printf("Warning: Could not load .local.env file: %v", err)
		log.Println("Continuing with system environment variables...")
	}

	// Initialize structured logger
	if err := logger.Init(); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	// Get required environment variables
	databaseURL := os.Getenv("DATABASE_URL")
	s3Bucket := os.Getenv("S3_BUCKET_NAME")
	s3Region := os.Getenv("S3_REGION")
	useMockS3 := os.Getenv("USE_MOCK_S3") == "true"

	// For development, allow running without real AWS S3
	if s3Bucket == "" && !useMockS3 {
		log.Println("Warning: S3_BUCKET_NAME not set, using mock S3 for development")
		useMockS3 = true
	}
	if s3Region == "" {
		s3Region = "us-east-1" // Default region
	}

	// Initialize database connection
	var db *pgxpool.Pool
	var err error

	if databaseURL != "" {
		db, err = pgxpool.New(context.Background(), databaseURL)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		// Test database connection
		if err := db.Ping(context.Background()); err != nil {
			log.Fatal("Failed to ping database:", err)
		}
		log.Println("Database connection established")
	} else {
		log.Println("Warning: DATABASE_URL not set, database operations will fail")
		// For development, you might want to create a mock database connection
	}

	// Create router configuration
	config := &routes.RouterConfig{
		DB:        db,
		S3Bucket:  s3Bucket,
		S3Region:  s3Region,
		UseMockS3: useMockS3,
	}

	// Initialize router
	r, err := routes.NewRouter(config)
	if err != nil {
		log.Fatal("Failed to initialize router:", err)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
