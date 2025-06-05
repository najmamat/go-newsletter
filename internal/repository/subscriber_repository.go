package repository

import (
	"context"
	"errors"
	"log/slog"

	"go-newsletter/pkg/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("not found")
)

type SubscriberRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewSubscriberRepository(db *pgxpool.Pool, logger *slog.Logger) *SubscriberRepository {
	return &SubscriberRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SubscriberRepository) ListByNewsletterID(ctx context.Context, newsletterID uuid.UUID) ([]*generated.Subscriber, error) {
	query := `
		SELECT id, newsletter_id, email, subscribed_at, is_confirmed
		FROM subscribers
		WHERE newsletter_id = $1
	`

	rows, err := r.db.Query(ctx, query, newsletterID)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to query subscribers", "error", err)
		return nil, err
	}
	defer rows.Close()

	var subscribers []*generated.Subscriber
	for rows.Next() {
		s := &generated.Subscriber{}
		err := rows.Scan(
			&s.Id,
			&s.NewsletterId,
			&s.Email,
			&s.SubscribedAt,
			&s.IsConfirmed,
		)
		if err != nil {
			r.logger.ErrorContext(ctx, "Failed to scan subscriber row", "error", err)
			return nil, err
		}
		subscribers = append(subscribers, s)
	}

	if err = rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "Error iterating subscriber rows", "error", err)
		return nil, err
	}

	return subscribers, nil
}