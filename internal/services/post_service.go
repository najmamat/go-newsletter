package services

import (
	"context"
	"errors"
	"go-newsletter/internal/repository"
	"log/slog"

	"go-newsletter/pkg/generated"

	"github.com/google/uuid"
)

type PostService struct {
	postRepo          *repository.PostRepository
	newsletterService *NewsletterService
	subscriberService *SubscriberService
	logger            *slog.Logger
}

func NewPostService(
	postRepo *repository.PostRepository,
	newsletterService *NewsletterService,
	subscriberService *SubscriberService,
	logger *slog.Logger,
) *PostService {
	return &PostService{
		postRepo:          postRepo,
		newsletterService: newsletterService,
		subscriberService: subscriberService,
		logger:            logger,
	}
}

// ListPosts retrieves a list of published posts for a newsletter
func (s *PostService) ListPosts(
	ctx context.Context,
	newsletterID uuid.UUID,
	editorID string,
) ([]*generated.PublishedPost, error) {
	// validate newsletter ownership
	_, err := s.newsletterService.GetNewsletterByIDCheckOwnership(ctx, newsletterID.String(), editorID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		s.logger.ErrorContext(ctx, "Failed to get newsletter", "error", err)
		return nil, err
	}

	posts, err := s.postRepo.ListByNewsletterID(ctx, newsletterID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list posts", "error", err)
		return nil, err
	}

	return posts, nil
}
