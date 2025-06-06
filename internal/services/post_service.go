package services

import (
	"context"
	"errors"
	"go-newsletter/internal/models"
	"go-newsletter/internal/repository"
	"log/slog"
	"strings"

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

func (s *PostService) CreatePost(ctx context.Context, editorID uuid.UUID, createPost generated.PublishPostRequest, newsletterId uuid.UUID) (*generated.PublishedPost, error) {
	// validate newsletter ownership
	_, err := s.newsletterService.GetNewsletterByIDCheckOwnership(ctx, newsletterId.String(), editorID.String())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		s.logger.ErrorContext(ctx, "Failed to get newsletter", "error", err)
		return nil, err
	}

	// Validate input
	if err := s.validatePublishPostRequest(createPost); err != nil {
		return nil, err
	}

	post, err := s.postRepo.CreatePost(ctx, editorID, &createPost, newsletterId)
	if err != nil {
		s.logger.ErrorContext(ctx, "SERVICE: failed to publish post", "error", err)
		return nil, err
	}
	return post, nil
}

// validatePublishPostRequest validates the post creation request
func (s *PostService) validatePublishPostRequest(post generated.PublishPostRequest) error {
	if strings.TrimSpace(post.Title) == "" {
		return models.NewBadRequestError("Title is required")
	}
	if post.ScheduledAt == nil {
		return models.NewBadRequestError("ScheduledAt is required")
	}

	return nil
}

func (s *PostService) UpdatePost(ctx context.Context, editorID uuid.UUID, postId uuid.UUID, updatePost generated.PublishPostRequest, newsletterId uuid.UUID) (*generated.PublishedPost, error) {
	_, err := s.newsletterService.GetNewsletterByIDCheckOwnership(ctx, newsletterId.String(), editorID.String())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		s.logger.ErrorContext(ctx, "Failed to get newsletter", "error", err)
		return nil, err
	}

	existingPost, err := s.postRepo.GetPostById(ctx, postId)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get post for update", "error", err)
		return nil, err
	}

	if existingPost.NewsletterId.String() != newsletterId.String() {
		return nil, models.NewForbiddenError("Post does not belong to the specified newsletter")
	}

	post, err := s.postRepo.UpdatePost(ctx, postId, &updatePost)
	if err != nil {
		s.logger.ErrorContext(ctx, "SERVICE: failed to update post", "error", err)
		return nil, err
	}
	return post, nil
}
