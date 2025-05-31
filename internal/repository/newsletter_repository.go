package repository

import (
	"context"
	"go-newsletter/internal/models"
	"go-newsletter/pkg/generated"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsletterRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewNewsletterRepository(db *pgxpool.Pool, logger *slog.Logger) *NewsletterRepository {
	return &NewsletterRepository{
		db:     db,
		logger: logger,
	}
}

// Retrieves a list of newsletters owned by the authenticated editor.
func (r *NewsletterRepository) GetNewslettersOwnedByEditor(ctx context.Context, editorID string) ([]generated.Newsletter, error) {
	query := `
		SELECT id, name, description, editor_id, created_at, updated_at
		FROM public.newsletters
		WHERE editor_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, editorID)
	if err != nil {
		r.logger.ErrorContext(ctx, "REPO: Failed to get all newsletters", "error", err)
		return nil, err
	}
	defer rows.Close()
	var newsletters []generated.Newsletter
	for rows.Next() {
		var n generated.Newsletter
		if err := rows.Scan(&n.Id,
			&n.Name,
			&n.Description,
			&n.EditorId,
			&n.CreatedAt,
			&n.UpdatedAt); err != nil {
			r.logger.ErrorContext(ctx, "Failed to scan newsletter row", "error", err)
			return nil, err
		}
		newsletters = append(newsletters, n)
	}

	if err := rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "Error iterating newsletter rows", "error", err)
		return nil, err
	}

	return newsletters, nil

}

func (r *NewsletterRepository) GetByID(ctx context.Context, newsletterID string) (*generated.Newsletter, error) {
	query := `
		SELECT id, name, description, editor_id, created_at, updated_at
		FROM public.newsletters
		WHERE id = $1
	`
	var n generated.Newsletter
	err := r.db.QueryRow(ctx, query, newsletterID).Scan(
		&n.Id,
		&n.Name,
		&n.Description,
		&n.EditorId,
		&n.CreatedAt,
		&n.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.ErrorContext(ctx, "REPO: Newsletter not found", "id", newsletterID)
			return nil, models.NewNotFoundError("Newsletter not found")
		}
		r.logger.ErrorContext(ctx, "REPO: Failed to get newsletter by ID", "id", newsletterID, "error", err)
		return nil, err
	}
	return &n, nil

}

func (r *NewsletterRepository) Create(ctx context.Context, editorID string, newsletterCreate *generated.NewsletterCreate) (*generated.Newsletter, error) {
	query := `
	INSERT INTO public.newsletters (id, name, description, editor_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, editor_id, created_at, updated_at
	`

	// ProfileRepo uses SQL NOW() func for this part.
	id := uuid.New()
	now := time.Now()

	var n generated.Newsletter
	err := r.db.QueryRow(ctx, query,
		id,
		newsletterCreate.Name,
		newsletterCreate.Description,
		editorID,
		now,
		now,
	).Scan(
		&n.Id,
		&n.Name,
		&n.Description,
		&n.EditorId,
		&n.CreatedAt,
		&n.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("REPO: failed to create newsletter", "error", err)
		return nil, err
	}

	return &n, nil
}

func (r *NewsletterRepository) Update(ctx context.Context, newsletterID string, newsletterUpdate *generated.NewsletterUpdate) (*generated.Newsletter, error) {
	// First get the current newsletter to handle partial updates
	current, err := r.GetByID(ctx, newsletterID)
	if err != nil {
		return nil, err
	}

	// Use existing values if update fields are not provided
	name := current.Name
	if newsletterUpdate.Name != nil {
		name = *newsletterUpdate.Name
	}

	description := current.Description
	if newsletterUpdate.Description != nil {
		description = newsletterUpdate.Description
	}

	query := `
		UPDATE public.newsletters
		SET name = $2, description = $3, updated_at = $4
		WHERE id = $1
		RETURNING id, name, description, editor_id, created_at, updated_at
	`
	now := time.Now()
	var n generated.Newsletter
	err = r.db.QueryRow(ctx, query, newsletterID, name, description, now).Scan(
		&n.Id,
		&n.Name,
		&n.Description,
		&n.EditorId,
		&n.CreatedAt,
		&n.UpdatedAt,
	)
	if err != nil {
		r.logger.Error("REPO: failed to update newsletter", "error", err)
		return nil, err
	}

	return &n, nil
}

func (r *NewsletterRepository) Delete(ctx context.Context, newsletterID string) error {
	query := `
		DELETE FROM public.newsletters
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, newsletterID)
	if err != nil {
		r.logger.ErrorContext(ctx, "REPO: failed to delete newsletter", "id", newsletterID, "error", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.logger.ErrorContext(ctx, "REPO: Newsletter not found for deletion", "id", newsletterID)
		return models.NewNotFoundError("Newsletter not found")
	}

	return nil
}

func (r *NewsletterRepository) AdminGetAll(ctx context.Context) ([]generated.Newsletter, error) {
	query := `
		SELECT id, name, description, editor_id, created_at, updated_at
		FROM public.newsletters
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.ErrorContext(ctx, "REPO: Failed to get all newsletters", "error", err)
		return nil, err
	}
	defer rows.Close()

	var newsletters []generated.Newsletter
	for rows.Next() {
		var n generated.Newsletter
		if err := rows.Scan(&n.Id,
			&n.Name,
			&n.Description,
			&n.EditorId,
			&n.CreatedAt,
			&n.UpdatedAt); err != nil {
			r.logger.ErrorContext(ctx, "Failed to scan newsletter row", "error", err)
			return nil, err
		}
		newsletters = append(newsletters, n)
	}

	if err := rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "Error iterating newsletter rows", "error", err)
		return nil, err
	}

	return newsletters, nil
}

func (r *NewsletterRepository) AdminDeleteByID(ctx context.Context, newsletterID string) error {
	query := `
		DELETE FROM public.newsletters
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, newsletterID)
	if err != nil {
		r.logger.ErrorContext(ctx, "REPO: failed to delete newsletter", "id", newsletterID, "error", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.logger.ErrorContext(ctx, "REPO: Newsletter not found for deletion", "id", newsletterID)
		return models.NewNotFoundError("Newsletter not found")
	}

	return nil
}
