// main.go
package main

import (
	"fmt"
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
	"brutal/internal/services"
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

	//Services
	profileService := services.NewProfileService()

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
	if handle == "" {
		http.Error(w, "Handle required", http.StatusBadRequest)
		return
	}

	// Fetch profile to get title/prompt
	profile, err := profileService.GetProfileByHandle(handle)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// Read HTML template
	data, err := os.ReadFile("web/form.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	// Replace placeholders
	html := string(data)
	html = strings.Replace(html, "{{.Handle}}", profile.Handle, -1)
	html = strings.Replace(html, "{{.Title}}", profile.Title, -1)

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

r.Get("/inbox/{handle}", func(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")
	if handle == "" {
		http.Error(w, "Handle required", http.StatusBadRequest)
		return
	}

	// Fetch profile to get title
	profile, err := profileService.GetProfileByHandle(handle)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// Read inbox template
	data, err := os.ReadFile("web/inbox.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
// Replace placeholders
	html := string(data)
	html = strings.Replace(html, "{{.Handle}}", profile.Handle, -1)
	html = strings.Replace(html, "{{.Title}}", profile.Title, -1)
	// ✅ Inject origin (e.g., http://localhost:8080 or https://brutal.app)
	origin := fmt.Sprintf("%s://%s", r.URL.Scheme, r.Host)
	if origin == ":///" || origin == "://localhost" { // fallback for local dev
		origin = "http://localhost:8080"
	}
	html = strings.Replace(html, "{{.Origin}}", origin, -1)
	

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
})

port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}


	// Start server using config port
	addr := ":" + port
	log.Printf("🚀 %s server running on http://localhost%s (env: %s)", cfg.AppName, addr, cfg.AppEnv)
	log.Fatal(http.ListenAndServe(addr, r))
}