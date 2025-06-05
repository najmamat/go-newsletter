package repository

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"go-newsletter/pkg/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	openapi_types "github.com/oapi-codegen/runtime/types"
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

// ExistsByEmail checks if a subscriber with the given email already exists for a newsletter
func (r *SubscriberRepository) ExistsByEmail(ctx context.Context, newsletterID uuid.UUID, email string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM subscribers
			WHERE newsletter_id = $1 AND email = $2
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, newsletterID, email).Scan(&exists)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to check subscriber existence", "error", err)
		return false, err
	}

	return exists, nil
}

// Create adds a new subscriber to a newsletter
func (r *SubscriberRepository) Create(ctx context.Context, newsletterID uuid.UUID, email string) (*generated.Subscriber, error) {
	query := `
		INSERT INTO subscribers (id, newsletter_id, email, subscribed_at, is_confirmed, unsubscribe_token)
		VALUES ($1, $2, $3, $4, false, $5)
		RETURNING id, newsletter_id, email, subscribed_at, is_confirmed, unsubscribe_token
	`

	unsubscribeToken := uuid.New().String()

	subscriber := &generated.Subscriber{
		Id:            &uuid.UUID{},
		NewsletterId:  &newsletterID,
		Email:         openapi_types.Email(email),
		SubscribedAt:  &time.Time{},
		IsConfirmed:   new(bool),
	}

	err := r.db.QueryRow(
		ctx,
		query,
		uuid.New(),
		newsletterID,
		email,
		time.Now().UTC(),
		unsubscribeToken,
	).Scan(
		&subscriber.Id,
		&subscriber.NewsletterId,
		&subscriber.Email,
		&subscriber.SubscribedAt,
		&subscriber.IsConfirmed,
		&subscriber.UnsubscribeToken,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to create subscriber", "error", err)
		return nil, err
	}

	return subscriber, nil
}