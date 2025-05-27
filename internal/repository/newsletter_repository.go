package repository

import (
	"context"
	"go-newsletter/pkg/generated"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewsletterRepository
type NewsletterRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// NewNewsletterRepository creates a new NewsletterRepository
func NewNewsletterRepository(db *pgxpool.Pool, logger *slog.Logger) *NewsletterRepository {
	return &NewsletterRepository{
		db:     db,
		logger: logger,
	}
}

func (r *NewsletterRepository) GetAll(ctx context.Context) ([]generated.Newsletter, error) {
	query := `
		SELECT id, name, description, created_at, updated_at, editor_id 
		FROM newsletters
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to query newsletters", "error", err)
		return nil, err
	}
	defer rows.Close()

	var newsletters []generated.Newsletter
	for rows.Next() {
		var n generated.Newsletter
		if err := rows.Scan(&n.Id, &n.Name, &n.Description, &n.CreatedAt, &n.UpdatedAt, &n.EditorId); err != nil {
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

func (r *NewsletterRepository) Insert(ctx context.Context, editorID string, req generated.NewsletterCreate) (*generated.Newsletter, error) {
	query := `
		INSERT INTO newsletters (id, name, description, created_at, updated_at, editor_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, created_at, updated_at, editor_id;
	`

	now := time.Now().UTC()
	id := uuid.New().String()

	var result generated.Newsletter
	err := r.db.QueryRow(ctx, query,
		id, req.Name, req.Description, now, now, editorID,
	).Scan(
		&result.Id, &result.Name, &result.Description, &result.CreatedAt, &result.UpdatedAt, &result.EditorId,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to insert newsletter", "error", err)
		return nil, err
	}

	return &result, nil
}
