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

func (s *NewsletterService) GetNewslettersOwnedByEditor(ctx context.Context, editorID string) ([]generated.Newsletter, error) {
	newsletters, err := s.repo.GetNewslettersOwnedByEditor(ctx, editorID)
	if err != nil {
		s.logger.ErrorContext(ctx, "SERVICE: failed to find newsletters of current editor", "error", err)
		return nil, err
	}
	return newsletters, nil
}

func (s *NewsletterService) GetNewsletterByID(ctx context.Context, newsletterID string, editorID string) (*generated.Newsletter, error) {
	newsletter, err := s.repo.GetByID(ctx, newsletterID)
	if err != nil {
		s.logger.ErrorContext(ctx, "SERVICE: failed to get newsletter by ID", "error", err)
		return nil, err
	}

	// Check if the requesting user is the editor of this newsletter
	if newsletter.EditorId.String() != editorID {
		s.logger.WarnContext(ctx, "SERVICE: unauthorized access attempt",
			"requested_editor_id", editorID,
			"newsletter_editor_id", newsletter.EditorId.String())
		return nil, models.NewForbiddenError("You don't have access to this newsletter")
	}

	return newsletter, nil
}

func (s *NewsletterService) CreateNewsletter(ctx context.Context, editorID string, newsletterCreate generated.NewsletterCreate) (*generated.Newsletter, error) {
	newsletter, err := s.repo.Create(ctx, editorID, &newsletterCreate)
	if err != nil {
		s.logger.ErrorContext(ctx, "SERVICE: failed to create newsletter", "error", err)
		return nil, err
	}
	return newsletter, nil
}

func (s *NewsletterService) UpdateNewsletter(ctx context.Context, editorID string, newsletterID string, newsletterUpdate generated.NewsletterUpdate) (*generated.Newsletter, error) {
	// First check if the newsletter exists and user has access
	newsletter, err := s.repo.GetByID(ctx, newsletterID)
	if err != nil {
		if models.IsNotFoundError(err) {
			s.logger.ErrorContext(ctx, "SERVICE: Newsletter not found", "id", newsletterID)
			return nil, err
		}
		s.logger.ErrorContext(ctx, "SERVICE: failed to get newsletter", "error", err)
		return nil, err
	}

	// Check if the requesting user is the editor of this newsletter
	if newsletter.EditorId.String() != editorID {
		s.logger.WarnContext(ctx, "SERVICE: unauthorized access attempt",
			"requested_editor_id", editorID,
			"newsletter_editor_id", newsletter.EditorId.String())
		return nil, models.NewForbiddenError("You don't have access to this newsletter")
	}

	// Proceed with update
	updatedNewsletter, err := s.repo.Update(ctx, newsletterID, &newsletterUpdate)
	if err != nil {
		s.logger.ErrorContext(ctx, "SERVICE: failed to update newsletter", "error", err)
		return nil, err
	}

	return updatedNewsletter, nil
}

func (s *NewsletterService) DeleteNewsletter(ctx context.Context, editorID string, newsletterID string) error {
	// First check if the newsletter exists and user has access
	newsletter, err := s.repo.GetByID(ctx, newsletterID)
	if err != nil {
		if models.IsNotFoundError(err) {
			s.logger.ErrorContext(ctx, "SERVICE: Newsletter not found", "id", newsletterID)
			return err
		}
		s.logger.ErrorContext(ctx, "SERVICE: failed to get newsletter", "error", err)
		return err
	}

	// Check if the requesting user is the editor of this newsletter
	if newsletter.EditorId.String() != editorID {
		s.logger.WarnContext(ctx, "SERVICE: unauthorized deletion attempt",
			"requested_editor_id", editorID,
			"newsletter_editor_id", newsletter.EditorId.String())
		return models.NewForbiddenError("You don't have access to this newsletter")
	}

	// Proceed with deletion
	if err := s.repo.Delete(ctx, newsletterID); err != nil {
		s.logger.ErrorContext(ctx, "SERVICE: failed to delete newsletter", "error", err)
		return err
	}

	return nil
}

func (s *NewsletterService) AdminGetAllNewsletters(ctx context.Context) ([]generated.Newsletter, error) {
	newsletters, err := s.repo.AdminGetAll(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "SERVICE: failed to get all newsletters", "error", err)
		return nil, err
	}
	return newsletters, nil
}

func (s *NewsletterService) AdminDeleteNewsletterByID(ctx context.Context, newsletterID string) error {
	if err := s.repo.AdminDeleteByID(ctx, newsletterID); err != nil {
		s.logger.ErrorContext(ctx, "SERVICE: failed to delete newsletter", "error", err)
		return err
	}
	return nil
}
