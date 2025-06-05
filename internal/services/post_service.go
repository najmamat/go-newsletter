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

// GetPostsByNewsletterId retrieves a list of published posts for a newsletter
func (s *PostService) GetPostsByNewsletterId(
	ctx context.Context,
	newsletterID uuid.UUID,
	editorID string,
	published bool,
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

	posts, err := s.postRepo.GetPostsByNewsletterId(ctx, newsletterID, published)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list posts", "error", err)
		return nil, err
	}

	return posts, nil
}

func (s *PostService) GetPostById(ctx context.Context, newsletterID uuid.UUID, postId uuid.UUID, editorID string) (*generated.PublishedPost, error) {
	// validate newsletter ownership
	_, err := s.newsletterService.GetNewsletterByIDCheckOwnership(ctx, newsletterID.String(), editorID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		s.logger.ErrorContext(ctx, "Failed to get newsletter", "error", err)
		return nil, err
	}

	post, err := s.postRepo.GetPostById(ctx, postId)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list post", "error", err)
		return nil, err
	}

	return post, nil
}

func (s *PostService) DeletePostById(ctx context.Context, newsletterID uuid.UUID, postId uuid.UUID, editorID string) error {
	// validate newsletter ownership
	_, err := s.newsletterService.GetNewsletterByIDCheckOwnership(ctx, newsletterID.String(), editorID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		s.logger.ErrorContext(ctx, "Failed to get newsletter", "error", err)
		return err
	}

	err = s.postRepo.DeletePostById(ctx, postId)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete post", "error", err)
		return err
	}

	return nil
}
