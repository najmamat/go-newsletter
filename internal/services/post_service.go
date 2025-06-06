package services

import (
	"context"
	"errors"
	"go-newsletter/internal/models"
	"go-newsletter/internal/models/enums"
	"go-newsletter/internal/repository"
	"log/slog"
	"strings"
	"time"

	"go-newsletter/pkg/generated"

	"github.com/google/uuid"
)

type PostService struct {
	postRepo          *repository.PostRepository
	newsletterService *NewsletterService
	subscriberService *SubscriberService
	mailingService    *MailingService
	logger            *slog.Logger
}

func NewPostService(
	postRepo *repository.PostRepository,
	newsletterService *NewsletterService,
	subscriberService *SubscriberService,
	mailingService *MailingService,
	logger *slog.Logger,
) *PostService {
	return &PostService{
		postRepo:          postRepo,
		newsletterService: newsletterService,
		subscriberService: subscriberService,
		mailingService:    mailingService,
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

	if *post.Status == enums.Posted.String() && post.PublishedAt != nil {
		if err := s.sendMailToSubscribers(ctx, post); err != nil {
			s.logger.ErrorContext(ctx, "Failed to send emails for new post", "error", err, "postId", post.Id)
		}
	}

	return post, nil
}

// sendMailToSubscribers sends a mail to all subscribers of a newsletter if the post is published
func (s *PostService) sendMailToSubscribers(ctx context.Context, post *generated.PublishedPost) error {
	if *post.Status != enums.Posted.String() || post.PublishedAt == nil {
		s.logger.InfoContext(ctx, "Skipping email sending for non-published post", "postId", post.Id, "status", post.Status)
		return nil
	}

	newsletter, err := s.newsletterService.GetNewsletterByID(ctx, post.NewsletterId.String())
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get newsletter for email", "error", err, "newsletterId", *post.NewsletterId)
		return err
	}

	subscribers, err := s.subscriberService.ListSubscribersWithouCheck(ctx, *post.NewsletterId)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get subscribers for newsletter", "error", err, "newsletterId", *post.NewsletterId)
		return err
	}

	if len(subscribers) == 0 {
		s.logger.InfoContext(ctx, "No subscribers for newsletter", "newsletterId", *post.NewsletterId)
		return nil
	}

	emailList := make([]string, 0, len(subscribers))
	for _, subscriber := range subscribers {
		if *subscriber.IsConfirmed {
			emailList = append(emailList, string(subscriber.Email))
		}
	}

	if len(emailList) == 0 {
		s.logger.InfoContext(ctx, "No confirmed subscribers for newsletter", "newsletterId", *post.NewsletterId)
		return nil
	}

	subject := post.Title
	if newsletter.Name != "" {
		subject = newsletter.Name + ": " + post.Title
	}

	err = s.mailingService.SendMail(emailList, subject, string(post.ContentHtml))
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to send newsletter email", "error", err, "postId", post.Id)
		return err
	}

	s.logger.InfoContext(ctx, "Newsletter email sent successfully", "postId", post.Id, "recipientCount", len(emailList))
	return nil
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

	if *post.Status == enums.Posted.String() && post.PublishedAt != nil {
		if err := s.sendMailToSubscribers(ctx, post); err != nil {
			s.logger.ErrorContext(ctx, "Failed to send emails for updated post", "error", err, "postId", post.Id)
		}
	}

	return post, nil
}

// GetPostsDueForPublication returns all scheduled posts that are due for publication
func (s *PostService) GetPostsDueForPublication(ctx context.Context, currentTime time.Time) ([]*generated.PublishedPost, error) {
	posts, err := s.postRepo.GetPostsDueForPublication(ctx, currentTime)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get posts due for publication", "error", err)
		return nil, err
	}
	return posts, nil
}

// PublishPost updates a post status to published and sends emails to subscribers
func (s *PostService) PublishPost(ctx context.Context, postId uuid.UUID) error {
	err := s.postRepo.PublishPost(ctx, postId)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to publish post", "postId", postId, "error", err)
		return err
	}

	post, err := s.postRepo.GetPostById(ctx, postId)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get published post for sending emails", "postId", postId, "error", err)
		return err
	}

	if err := s.sendMailToSubscribers(ctx, post); err != nil {
		s.logger.ErrorContext(ctx, "Failed to send emails for published post", "error", err, "postId", postId)
	}

	return nil
}
