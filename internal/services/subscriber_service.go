package services

import (
	"context"
	"errors"
	"log/slog"

	"go-newsletter/internal/repository"
	"go-newsletter/pkg/generated"

	"github.com/google/uuid"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
)

type SubscriberService struct {
	newsletterRepo *repository.NewsletterRepository
	subscriberRepo *repository.SubscriberRepository
	logger         *slog.Logger
}

func NewSubscriberService(
	newsletterRepo *repository.NewsletterRepository,
	subscriberRepo *repository.SubscriberRepository,
	logger *slog.Logger,
) *SubscriberService {
	return &SubscriberService{
		newsletterRepo: newsletterRepo,
		subscriberRepo: subscriberRepo,
		logger:         logger,
	}
}

// ListSubscribers retrieves a list of subscribers for a newsletter
func (s *SubscriberService) ListSubscribers(
	ctx context.Context,
	newsletterID uuid.UUID,
	editorID string,
) ([]generated.Subscriber, error) {
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

	var result []generated.Subscriber
	for _, s := range subscribers {
		result = append(result, *s)
	}
	return result, nil
}