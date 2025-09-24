package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"brutal/internal/services"
)

type MessageHandler struct{
	messageService *services.MessageService
	profileService *services.ProfileService
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		messageService: services.NewMessageService(),
		profileService: services.NewProfileService(),
	}
}

func (h *MessageHandler) PostMessage(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")

	// 1. Validate profile exists
	profile, err := h.profileService.GetProfileByHandle(handle)
	if err != nil {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	// 2. Parse message
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}

	// 3. Get IP (simplified)
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	// 4. Save message
	msg, err := h.messageService.CreateMessage(profile.ID, req.Content, ip, r.UserAgent())
	if err != nil {
		http.Error(w, "failed to save message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")
	profile, err := h.profileService.GetProfileByHandle(handle)
	if err != nil {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}


	messages, err := h.messageService.GetMessagesByProfile(profile.ID)
	if err != nil {
		http.Error(w, "failed to fetch messages", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}