package repository

import (
	"context"
	"log/slog"

	"go-newsletter/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ProfileRepository handles data access for profiles
type ProfileRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// NewProfileRepository creates a new ProfileRepository
func NewProfileRepository(db *pgxpool.Pool, logger *slog.Logger) *ProfileRepository {
	return &ProfileRepository{
		db:     db,
		logger: logger,
	}
}

// GetAll retrieves all profiles from the database
func (r *ProfileRepository) GetAll(ctx context.Context) ([]models.Profile, error) {
	query := `
		SELECT id, full_name, avatar_url, is_admin, created_at, updated_at 
		FROM public.profiles 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to query profiles", "error", err)
		return nil, err
	}
	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		var p models.Profile
		if err := rows.Scan(&p.ID, &p.FullName, &p.AvatarURL, &p.IsAdmin, &p.CreatedAt, &p.UpdatedAt); err != nil {
			r.logger.ErrorContext(ctx, "Failed to scan profile row", "error", err)
			return nil, err
		}
		profiles = append(profiles, p)
	}

	if err := rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "Error iterating profile rows", "error", err)
		return nil, err
	}

	return profiles, nil
}

// GetByID retrieves a single profile by ID
func (r *ProfileRepository) GetByID(ctx context.Context, id string) (*models.Profile, error) {
	query := `
		SELECT id, full_name, avatar_url, is_admin, created_at, updated_at 
		FROM public.profiles 
		WHERE id = $1
	`

	var p models.Profile
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.FullName, &p.AvatarURL, &p.IsAdmin, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to get profile by ID", "id", id, "error", err)
		return nil, err
	}

	return &p, nil
}

// Update updates a profile's editable fields
func (r *ProfileRepository) Update(ctx context.Context, id string, req models.UpdateProfileRequest) (*models.Profile, error) {
	query := `
		UPDATE public.profiles 
		SET full_name = $2, avatar_url = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, full_name, avatar_url, is_admin, created_at, updated_at
	`

	var p models.Profile
	err := r.db.QueryRow(ctx, query, id, req.FullName, req.AvatarURL).Scan(
		&p.ID, &p.FullName, &p.AvatarURL, &p.IsAdmin, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to update profile", "id", id, "error", err)
		return nil, err
	}

	return &p, nil
} 