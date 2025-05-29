package models

import "time"

// Profile represents a user profile in the system
type Profile struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	AvatarURL string    `json:"avatar_url"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} 