package main

import (
	"log"
	"net/http"

	"github.com/angel-romero-f/rice-notes/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .local.env file
	if err := godotenv.Load(".local.env"); err != nil {
		log.Printf("Warning: Could not load .local.env file: %v", err)
		log.Println("Continuing with system environment variables...")
	}
	
	r := routes.NewRouter()
	
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
