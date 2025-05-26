package models

import "time"

// Newsletter represents a newsletter created by an editor
type Newsletter struct {
	ID          string    `json:"id" db:"id"`
	EditorID    string    `json:"editor_id" db:"editor_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// NewsletterCreateRequest is used when creating a new newsletter
type NewsletterCreateRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}
