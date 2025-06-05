package repository

import (
	"context"
	"log/slog"

	"go-newsletter/pkg/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewPostRepository(db *pgxpool.Pool, logger *slog.Logger) *PostRepository {
	return &PostRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostRepository) GetPostsByNewsletterId(ctx context.Context, newsletterID uuid.UUID, published bool) ([]*generated.PublishedPost, error) {
	query := `
		SELECT id, newsletter_id, editor_id, title, content_html, content_text, status, scheduled_at, published_at, created_at
		FROM published_posts
		WHERE newsletter_id = $1`

	if published {
		query += ` AND published_at IS NOT NULL`
	} else {
		query += ` AND published_at IS NULL`
	}

	rows, err := r.db.Query(ctx, query, newsletterID)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to query posts", "error", err)
		return nil, err
	}
	defer rows.Close()

	var posts []*generated.PublishedPost
	for rows.Next() {
		s := &generated.PublishedPost{}
		err := rows.Scan(
			&s.Id,
			&s.NewsletterId,
			&s.EditorId,
			&s.Title,
			&s.ContentText,
			&s.ContentHtml,
			&s.Status,
			&s.ScheduledAt,
			&s.PublishedAt,
			&s.CreatedAt,
		)
		if err != nil {
			r.logger.ErrorContext(ctx, "Failed to scan post row", "error", err)
			return nil, err
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "Error iterating post rows", "error", err)
		return nil, err
	}

	return posts, nil
}
