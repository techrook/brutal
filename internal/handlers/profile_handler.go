package handlers

import (
	"encoding/json"
	"net/http"

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

func (h *ProfileHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Handle      string `json:"handle"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	profile, err := h.profileService.CreateProfile(req.Handle, req.Title, req.Description)
	if err != nil {
		http.Error(w, "failed to create profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(profile)
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")
	profile, err := h.profileService.GetProfileByHandle(handle)
	if err != nil {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(profile)
}