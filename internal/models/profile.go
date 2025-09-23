package models

import "time"

type Profile struct{
	ID          string    `db:"id"`
	UserID      string    `db:"user_id"` 
	Handle      string    `db:"handle"`      
	Title       string    `db:"title"`        
	Description string    `db:"description"` 
	IsActive    bool      `db:"is_active"`   
	CreatedAt   time.Time `db:"created_at"`
}