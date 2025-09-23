package models

import "time"

type Message struct {
	ID          string    `db:"id"`
	ProfileID   string    `db:"profile_id"`
	Content     string    `db:"content"`    
	IPAddress   string    `db:"ip_address"`   
	UserAgent   string    `db:"user_agent"`   
	IsHidden    bool      `db:"is_hidden"`    
	IsFlagged   bool      `db:"is_flagged"`   
	CreatedAt   time.Time `db:"created_at"`
}