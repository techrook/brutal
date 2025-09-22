// main.go
package main

import (
	"net/http"
	"log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Create a new Chi router
	r := chi.NewRouter()

	// Add basic middleware (logging + panic recovery)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Define a simple route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Brutal 💥 — Truth Hurts."))
	})


	log.Println("🚀 Brutal server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}