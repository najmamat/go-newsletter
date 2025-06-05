package services

import (
	"context"
	"errors"
	"log/slog"

	"go-newsletter/internal/repository"
	"go-newsletter/pkg/generated"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/google/uuid"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrForbidden        = errors.New("forbidden")
	ErrAlreadySubscribed = errors.New("already subscribed")
)

type SubscriberService struct {
	subscriberRepo *repository.SubscriberRepository
	newsletterRepo *repository.NewsletterRepository
	logger         *slog.Logger
}

func NewSubscriberService(
	subscriberRepo *repository.SubscriberRepository,
	newsletterRepo *repository.NewsletterRepository,
	logger *slog.Logger,
) *SubscriberService {
	return &SubscriberService{
		subscriberRepo: subscriberRepo,
		newsletterRepo: newsletterRepo,
		logger:         logger,
	}
}

// ListSubscribers retrieves a list of subscribers for a newsletter
func (s *SubscriberService) ListSubscribers(
	ctx context.Context,
	newsletterID uuid.UUID,
	editorID string,
) ([]*generated.Subscriber, error) {
	// Verify newsletter ownership
	newsletter, err := s.newsletterRepo.GetByID(ctx, newsletterID.String())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		s.logger.ErrorContext(ctx, "Failed to get newsletter", "error", err)
		return nil, err
	}

	if newsletter.EditorId.String() != editorID {
		return nil, ErrForbidden
	}

	// Get subscribers
	subscribers, err := s.subscriberRepo.ListByNewsletterID(ctx, newsletterID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list subscribers", "error", err)
		return nil, err
	}

	return subscribers, nil
}

// Subscribe adds a new subscriber to a newsletter
func (s *SubscriberService) Subscribe(
	ctx context.Context,
	newsletterID uuid.UUID,
	email openapi_types.Email,
) (*generated.Subscriber, error) {
	// Check if newsletter exists
	_, err := s.newsletterRepo.GetByID(ctx, newsletterID.String())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		s.logger.ErrorContext(ctx, "Failed to get newsletter", "error", err)
		return nil, err
	}

	// Check if already subscribed
	exists, err := s.subscriberRepo.ExistsByEmail(ctx, newsletterID, string(email))
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check subscription", "error", err)
		return nil, err
	}
	if exists {
		return nil, ErrAlreadySubscribed
	}

	// Create subscriber
	subscriber, err := s.subscriberRepo.Create(ctx, newsletterID, string(email))
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create subscriber", "error", err)
		return nil, err
	}

	// TODO: Send confirmation email

	return subscriber, nil
}

// ConfirmSubscription confirms a subscription using a confirmation token
func (s *SubscriberService) ConfirmSubscription(ctx context.Context, token string) error {
	err := s.subscriberRepo.ConfirmByToken(ctx, token)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}