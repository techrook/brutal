// main.go
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"brutal/internal/config"
)

func main() {
	// Load config from .env
	cfg := config.LoadConfig()

	// Create Chi router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from " + cfg.AppName + " 💥 — Truth Hurts. Running in " + cfg.AppEnv + " mode."))
	})

	// Start server using config port
	addr := ":" + cfg.ServerPort
	log.Printf("🚀 %s server running on http://localhost%s (env: %s)", cfg.AppName, addr, cfg.AppEnv)
	log.Fatal(http.ListenAndServe(addr, r))
}