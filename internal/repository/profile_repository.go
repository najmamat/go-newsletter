package repository

import (
	"context"
	"log/slog"

	"go-newsletter/pkg/generated"

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
func (r *ProfileRepository) GetAll(ctx context.Context) ([]generated.EditorProfile, error) {
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

	var profiles []generated.EditorProfile
	for rows.Next() {
		var p generated.EditorProfile
		if err := rows.Scan(&p.Id, &p.FullName, &p.AvatarUrl, &p.IsAdmin, &p.CreatedAt, &p.UpdatedAt); err != nil {
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
func (r *ProfileRepository) GetByID(ctx context.Context, id string) (*generated.EditorProfile, error) {
	query := `
		SELECT id, full_name, avatar_url, is_admin, created_at, updated_at 
		FROM public.profiles 
		WHERE id = $1
	`

	var p generated.EditorProfile
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.Id, &p.FullName, &p.AvatarUrl, &p.IsAdmin, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to get profile by ID", "id", id, "error", err)
		return nil, err
	}

	return &p, nil
}

// Update updates a profile's editable fields
func (r *ProfileRepository) Update(ctx context.Context, id string, req generated.PutMeJSONBody) (*generated.EditorProfile, error) {
	query := `
		UPDATE public.profiles 
		SET full_name = $2, avatar_url = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, full_name, avatar_url, is_admin, created_at, updated_at
	`

	var p generated.EditorProfile
	err := r.db.QueryRow(ctx, query, id, req.FullName, req.AvatarUrl).Scan(
		&p.Id, &p.FullName, &p.AvatarUrl, &p.IsAdmin, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to update profile", "id", id, "error", err)
		return nil, err
	}

	return &p, nil
}

// Create creates a new profile for a user
func (r *ProfileRepository) Create(ctx context.Context, id string) (*generated.EditorProfile, error) {
	query := `
		INSERT INTO public.profiles (id, full_name, avatar_url, is_admin, created_at, updated_at)
		VALUES ($1, '', '', false, NOW(), NOW())
		RETURNING id, full_name, avatar_url, is_admin, created_at, updated_at
	`

	var p generated.EditorProfile
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.Id, &p.FullName, &p.AvatarUrl, &p.IsAdmin, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to create profile", "id", id, "error", err)
		return nil, err
	}

	return &p, nil
}

// GrantAdmin grants admin privileges to a user
func (r *ProfileRepository) GrantAdmin(ctx context.Context, id string) (*generated.EditorProfile, error) {
	query := `
		UPDATE public.profiles 
		SET is_admin = true, updated_at = NOW()
		WHERE id = $1
		RETURNING id, full_name, avatar_url, is_admin, created_at, updated_at
	`

	var p generated.EditorProfile
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.Id, &p.FullName, &p.AvatarUrl, &p.IsAdmin, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to grant admin privileges", "id", id, "error", err)
		return nil, err
	}

	return &p, nil
}

// RevokeAdmin revokes admin privileges from a user
func (r *ProfileRepository) RevokeAdmin(ctx context.Context, id string) (*generated.EditorProfile, error) {
	query := `
		UPDATE public.profiles 
		SET is_admin = false, updated_at = NOW()
		WHERE id = $1
		RETURNING id, full_name, avatar_url, is_admin, created_at, updated_at
	`

	var p generated.EditorProfile
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.Id, &p.FullName, &p.AvatarUrl, &p.IsAdmin, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to revoke admin privileges", "id", id, "error", err)
		return nil, err
	}

	return &p, nil
} 