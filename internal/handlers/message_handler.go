package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"brutal/internal/models"
	"brutal/internal/services"

	"github.com/go-chi/chi/v5"
)

type MessageHandler struct {
	messageService *services.MessageService
	profileService *services.ProfileService
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		messageService: services.NewMessageService(),
		profileService: services.NewProfileService(),
	}
}

// ErrorResponse matches the one in profile_handler.go


func (h *MessageHandler) PostMessage(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")
	if handle == "" {
		writeError(w, http.StatusBadRequest, "Handle is required")
		return
	}

	// 1. Validate profile exists
	profile, err := h.profileService.GetProfileByHandle(handle)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "Profile not found. Did you type the username correctly?")
			return
		}
		// Log real error in production
		writeError(w, http.StatusInternalServerError, "Could not verify profile. Please try again.")
		return
	}

	// 2. Parse message
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		writeError(w, http.StatusBadRequest, "Message cannot be empty")
		return
	}
	if len(content) > 1000 {
		writeError(w, http.StatusBadRequest, "Message is too long (max 1000 characters)")
		return
	}

	// 3. Get IP (simplified)
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	// 4. Save message
	msg, err := h.messageService.CreateMessage(profile.ID, content, ip, r.UserAgent())
	if err != nil {
		// Log real error in production
		writeError(w, http.StatusInternalServerError, "Failed to send message. Please try again.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")
	if handle == "" {
		writeError(w, http.StatusBadRequest, "Handle is required")
		return
	}

	profile, err := h.profileService.GetProfileByHandle(handle)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "Profile not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Could not load messages. Please try again.")
		return
	}

	messages, err := h.messageService.GetMessagesByProfile(profile.ID)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "Failed to fetch messages")
        return
    }

    // Ensure it's never nil
    if messages == nil {
        messages = []*models.Message{}
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(messages) // always []
}