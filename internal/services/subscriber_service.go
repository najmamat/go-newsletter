package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go-newsletter/internal/config"
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
	mailingService *MailingService
	logger         *slog.Logger
	config         *config.Config
}

func NewSubscriberService(
	subscriberRepo *repository.SubscriberRepository,
	newsletterRepo *repository.NewsletterRepository,
	mailingService *MailingService,
	config *config.Config,
	logger *slog.Logger,
) *SubscriberService {
	return &SubscriberService{
		subscriberRepo: subscriberRepo,
		newsletterRepo: newsletterRepo,
		mailingService: mailingService,
		config:         config,
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
	newsletter, err := s.newsletterRepo.GetByID(ctx, newsletterID.String())
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

	// Send confirmation email
	confirmationLink := fmt.Sprintf("%s/api/v1/subscribe/confirm/%s", s.config.Server.APIBaseURL, *subscriber.ConfirmationToken)
	htmlContent := fmt.Sprintf(`
		<h1>Confirm Your Subscription to %s</h1>
		<p>Thank you for subscribing to our newsletter! Please click the link below to confirm your subscription:</p>
		<p><a href="%s">Confirm Subscription</a></p>
		<p>If you did not request this subscription, you can safely ignore this email.</p>
	`, newsletter.Name, confirmationLink)

	err = s.mailingService.SendMail([]string{string(email)}, "Confirm Your Newsletter Subscription", htmlContent)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to send confirmation email", "error", err)
	}

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

// Unsubscribe handles unsubscription using a token
func (s *SubscriberService) Unsubscribe(ctx context.Context, token string) error {
	err := s.subscriberRepo.UnsubscribeByToken(ctx, token)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}