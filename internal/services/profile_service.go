package services

import (
	"brutal/internal/db"
	"brutal/internal/models"
	"fmt"
	"strings"
	"time"
)

type ProfileService struct{}

func NewProfileService() *ProfileService {
	return &ProfileService{}
}

func (s *ProfileService) CreateProfile(handle, title, description string) (*models.Profile, error) {

	handle = strings.ToLower(strings.TrimSpace(handle))
	if len(handle) < 3 {
		return nil, fmt.Errorf("handle must be at least 3 characters")
	}

	profile := &models.Profile{
		ID:          "",
		Handle:      handle, // now guaranteed lowercase
		Title:       title,
		Description: description,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	query := `
		INSERT INTO profiles (handle, title, description, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := db.DB.QueryRow(
		query,
		profile.Handle,
		profile.Title,
		profile.Description,
		profile.IsActive,
		profile.CreatedAt,
	).Scan(&profile.ID)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *ProfileService) GetProfileByHandle(handle string) (*models.Profile, error) {
	
	handle = strings.ToLower(strings.TrimSpace(handle))

	var profile models.Profile
	query := "SELECT * FROM profiles WHERE handle = $1 AND is_active = true"
	err := db.DB.Get(&profile, query, handle)
	return &profile, err
}