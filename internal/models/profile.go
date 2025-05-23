package models

import "time"

// Profile represents a user profile in the system
type Profile struct {
	ID        string     `json:"id" db:"id"`
	FullName  *string    `json:"full_name,omitempty" db:"full_name"`
	AvatarURL *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	IsAdmin   bool       `json:"is_admin" db:"is_admin"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateProfileRequest represents the request payload for creating/updating profiles
type CreateProfileRequest struct {
	FullName  *string `json:"full_name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// UpdateProfileRequest represents the request payload for updating profiles
type UpdateProfileRequest struct {
	FullName  *string `json:"full_name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
} 