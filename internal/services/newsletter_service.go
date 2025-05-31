package services

import (
	"context"
	"go-newsletter/internal/models"
	"go-newsletter/internal/repository"
	"go-newsletter/pkg/generated"
	"log/slog"
)

type NewsletterService struct {
	repo   *repository.NewsletterRepository
	logger *slog.Logger
}

func NewNewsletterService(repo *repository.NewsletterRepository, logger *slog.Logger) *NewsletterService {
	return &NewsletterService{
		repo:   repo,
		logger: logger,
	}
}

func (s *NewsletterService) GetNewslettersOwnedByEditor(ctx context.Context, editorId string) ([]generated.Newsletter, error) {
	newsletters, err := s.repo.GetNewslettersOwnedByEditor(ctx, editorId)
	if err != nil {
		s.logger.ErrorContext(ctx, "SERVICE: failed to find newsletters of current editor", "error", err)
		return nil, err
	}

	var result []generated.Newsletter
	for _, n := range newsletters {
		result = append(result, n)
	}
	return result, nil
}

func (s *NewsletterService) GetNewsletterByID(ctx context.Context, newsletterId string, editorId string) (*generated.Newsletter, error) {
	newsletter, err := s.repo.GetByID(ctx, newsletterId)
	if err != nil {
		s.logger.ErrorContext(ctx, "SERVICE: failed to get newsletter by ID", "error", err)
		return nil, err
	}

	// Check if the requesting user is the editor of this newsletter
	if newsletter.EditorId.String() != editorId {
		s.logger.WarnContext(ctx, "SERVICE: unauthorized access attempt",
			"requested_editor_id", editorId,
			"newsletter_editor_id", newsletter.EditorId.String())
		return nil, models.NewForbiddenError("You don't have access to this newsletter")
	}

	return newsletter, nil
}

func (s *NewsletterService) CreateNewsletter(ctx context.Context, editorId string, newsletterCreate generated.NewsletterCreate) (*generated.Newsletter, error) {
	newsletter, err := s.repo.Create(ctx, editorId, &newsletterCreate)
	if err != nil {
		s.logger.ErrorContext(ctx, "SERVICE: failed to create newsletter", "error", err)
		return nil, err
	}
	return newsletter, nil
}
