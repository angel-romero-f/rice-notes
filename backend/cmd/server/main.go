package main

import (
	"log"
	"net/http"

	"github.com/angel-romero-f/rice-notes/internal/routes"
)

func main() {
	r := routes.NewRouter()
	http.ListenAndServe(":3000", r)

	log.Println("Server running on :3000")

}
