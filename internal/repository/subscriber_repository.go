package repository

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"go-newsletter/pkg/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		INSERT INTO subscribers (id, newsletter_id, email, subscribed_at, is_confirmed, unsubscribe_token, confirmation_token)
		VALUES ($1, $2, $3, $4, false, $5, $6)
		RETURNING id, newsletter_id, email, subscribed_at, is_confirmed, unsubscribe_token, confirmation_token
	`

	unsubscribeToken := uuid.New().String()
	confirmationToken := uuid.New().String()
	
	subscriber := &generated.Subscriber{
		Id:            &uuid.UUID{},
		NewsletterId:  &newsletterID,
		Email:         openapi_types.Email(email),
		SubscribedAt:  &time.Time{},
		IsConfirmed:   new(bool),
		UnsubscribeToken: &unsubscribeToken,
		ConfirmationToken: &confirmationToken,
	}

	err := r.db.QueryRow(
		ctx,
		query,
		uuid.New(),
		newsletterID,
		email,
		time.Now().UTC(),
		unsubscribeToken,
		confirmationToken,
	).Scan(
		&subscriber.Id,
		&subscriber.NewsletterId,
		&subscriber.Email,
		&subscriber.SubscribedAt,
		&subscriber.IsConfirmed,
		&subscriber.UnsubscribeToken,
		&subscriber.ConfirmationToken,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to create subscriber", "error", err)
		return nil, err
	}

	return subscriber, nil
}

// ConfirmByToken confirms a subscription using a confirmation token
func (r *SubscriberRepository) ConfirmByToken(ctx context.Context, token string) error {
	query := `
		UPDATE subscribers
		SET is_confirmed = true
		WHERE confirmation_token = $1
		RETURNING id
	`

	var id uuid.UUID
	err := r.db.QueryRow(ctx, query, token).Scan(&id)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to confirm subscription", "error", err)
		return err
	}

	return nil
}

// UnsubscribeByToken unsubscribes a user using their unsubscribe token
func (r *SubscriberRepository) UnsubscribeByToken(ctx context.Context, token string) error {
	query := `
		UPDATE subscribers
		SET unsubscribed_at = NOW()
		WHERE unsubscribe_token = $1 AND unsubscribed_at IS NULL
		RETURNING id
	`

	var id uuid.UUID
	err := r.db.QueryRow(ctx, query, token).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrNotFound
		}
		r.logger.ErrorContext(ctx, "Failed to unsubscribe", "error", err)
		return err
	}

	return nil
}