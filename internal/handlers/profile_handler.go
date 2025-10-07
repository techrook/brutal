package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"brutal/internal/services"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{
		profileService: services.NewProfileService(),
	}
}

// ErrorResponse standardizes API errors
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}

func (h *ProfileHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Handle      string `json:"handle"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Validate input
	handle := strings.TrimSpace(req.Handle)
	title := strings.TrimSpace(req.Title)

	if handle == "" {
		writeError(w, http.StatusBadRequest, "Handle is required")
		return
	}
	if len(handle) < 3 {
		writeError(w, http.StatusBadRequest, "Handle must be at least 3 characters")
		return
	}
	if title == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}

	profile, err := h.profileService.CreateProfile(handle, title, req.Description)
	if err != nil {
		// Handle DB constraint violations (e.g., duplicate handle)
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			writeError(w, http.StatusConflict, "This username is already taken. Try another!")
			return
		}

		// Log unexpected errors (in production, use structured logger)
		// log.Printf("CreateProfile error: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to create profile. Please try again.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(profile)
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")
	if handle == "" {
		writeError(w, http.StatusBadRequest, "Handle is required")
		return
	}

	profile, err := h.profileService.GetProfileByHandle(handle)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "Profile not found. Did you type the username correctly?")
			return
		}

		// Other DB errors (connection, scan, etc.)
		// log.Printf("GetProfile error: %v", err)
		writeError(w, http.StatusInternalServerError, "Could not fetch profile. Please try again.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}