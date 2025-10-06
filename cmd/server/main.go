// main.go
package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"brutal/internal/config"
	"brutal/internal/db"
	"brutal/internal/handlers"
)

func main() {
	// Load config from .env
	cfg := config.LoadConfig()
	//connect to database
	db.InitDB(cfg)
	//migration
	db.RunMigrations(db.DB) 
	// Create Chi router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Handlers
	profileHandler := handlers.NewProfileHandler()
	messageHandler := handlers.NewMessageHandler()

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/profiles", profileHandler.CreateProfile)
		r.Get("/profiles/{handle}", profileHandler.GetProfile)

		r.Post("/profiles/{handle}/messages", messageHandler.PostMessage)
		r.Get("/profiles/{handle}/messages", messageHandler.GetMessages)
	})

	// Route
// Serve static HTML files
r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
})

r.Get("/f/{handle}", func(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")
	log.Print(handle)
	// Read form.html and inject handle (simple template)
	data, _ := os.ReadFile("web/form.html")
	html := strings.Replace(string(data), "{{.Handle}}", handle, 1)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
})

r.Get("/f/{handle}/sent", func(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")
	log.Print(handle)
	data, _ := os.ReadFile("web/sent.html")
	html := strings.Replace(string(data), "handle=", "handle="+handle, 1)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
})

	// Start server using config port
	addr := ":" + cfg.ServerPort
	log.Printf("🚀 %s server running on http://localhost%s (env: %s)", cfg.AppName, addr, cfg.AppEnv)
	log.Fatal(http.ListenAndServe(addr, r))
}