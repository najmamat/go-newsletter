package services

import (
	"context"
	"go-newsletter/internal/config"
	"go-newsletter/internal/models"
	"go-newsletter/internal/repository"
	"go-newsletter/pkg/generated"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

// Validation errors
var (
	ErrNameRequired   = models.NewBadRequestError("Newsletter name is required")
	ErrNameTooLong    = models.NewBadRequestError("Newsletter name must be less than 100 characters")
	ErrDescTooLong    = models.NewBadRequestError("Description must be less than 500 characters")
	ErrInvalidID      = models.NewBadRequestError("Invalid newsletter ID format")
	ErrNoUpdateFields = models.NewBadRequestError("At least one field (name or description) must be provided for update")
	ErrEmptyName      = models.NewBadRequestError("Newsletter name cannot be empty")
)

type NewsletterService struct {
	repo   *repository.NewsletterRepository
	logger *slog.Logger
	config *config.NewsletterConfig
}

func NewNewsletterService(repo *repository.NewsletterRepository, logger *slog.Logger) *NewsletterService {
	return &NewsletterService{
		repo:   repo,
		logger: logger,
		config: config.DefaultNewsletterConfig(),
	}
}

// validateNewsletterCreate validates the newsletter creation request
func (s *NewsletterService) validateNewsletterCreate(ctx context.Context, editorID string, newsletter generated.NewsletterCreate) error {
	if strings.TrimSpace(newsletter.Name) == "" {
		return models.NewBadRequestError(s.config.RequiredNameMessage)
	}
	if len(newsletter.Name) > s.config.MaxNameLength {
		return models.NewBadRequestError(s.config.TooLongNameMessage)
	}
	if len(newsletter.Name) < s.config.MinNameLength {
		return models.NewBadRequestError(s.config.TooShortNameMessage)
	}
	if newsletter.Description != nil && len(*newsletter.Description) > s.config.MaxDescriptionLength {
		return models.NewBadRequestError(s.config.TooLongDescMessage)
	}

	// Check for duplicate name
	exists, err := s.repo.CheckDuplicateName(ctx, editorID, newsletter.Name, "")
	if err != nil {
		return err
	}
	if exists {
		return models.NewBadRequestError(s.config.DuplicateNameMessage)
	}

	return nil
}

// validateNewsletterUpdate validates the newsletter update request
func (s *NewsletterService) validateNewsletterUpdate(ctx context.Context, editorID string, newsletterID string, update generated.NewsletterUpdate) error {
	if update.Name != nil {
		if strings.TrimSpace(*update.Name) == "" {
			return models.NewBadRequestError(s.config.EmptyNameMessage)
		}
		if len(*update.Name) > s.config.MaxNameLength {
			return models.NewBadRequestError(s.config.TooLongNameMessage)
		}
		if len(*update.Name) < s.config.MinNameLength {
			return models.NewBadRequestError(s.config.TooShortNameMessage)
		}

		// Check for duplicate name
		exists, err := s.repo.CheckDuplicateName(ctx, editorID, *update.Name, newsletterID)
		if err != nil {
			return err
		}
		if exists {
			return models.NewBadRequestError(s.config.DuplicateNameMessage)
		}
	}
	if update.Description != nil && len(*update.Description) > s.config.MaxDescriptionLength {
		return models.NewBadRequestError(s.config.TooLongDescMessage)
	}
	return nil
}

// validateNewsletterID validates the newsletter ID format
func (s *NewsletterService) validateNewsletterID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return models.NewBadRequestError(s.config.InvalidIDMessage)
	}
	return nil
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
	// Validate input
	if err := s.validateNewsletterID(newsletterID); err != nil {
		return nil, err
	}

	newsletter, err := s.repo.GetByID(ctx, newsletterID)
	if err != nil {
		s.logger.ErrorContext(ctx, "SERVICE: failed to get newsletter by ID", "error", err)
		return nil, err
	}

	if err := s.checkNewsletterOwnership(ctx, newsletter, editorID); err != nil {
		return nil, err
	}

	return newsletter, nil
}

func (s *NewsletterService) CreateNewsletter(ctx context.Context, editorID string, newsletterCreate generated.NewsletterCreate) (*generated.Newsletter, error) {
	// Validate input
	if err := s.validateNewsletterCreate(ctx, editorID, newsletterCreate); err != nil {
		return nil, err
	}

	newsletter, err := s.repo.Create(ctx, editorID, &newsletterCreate)
	if err != nil {
		s.logger.ErrorContext(ctx, "SERVICE: failed to create newsletter", "error", err)
		return nil, err
	}
	return newsletter, nil
}

func (s *NewsletterService) UpdateNewsletter(ctx context.Context, editorID string, newsletterID string, newsletterUpdate generated.NewsletterUpdate) (*generated.Newsletter, error) {
	// Validate input
	if err := s.validateNewsletterID(newsletterID); err != nil {
		return nil, err
	}
	if err := s.validateNewsletterUpdate(ctx, editorID, newsletterID, newsletterUpdate); err != nil {
		return nil, err
	}

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

	if err := s.checkNewsletterOwnership(ctx, newsletter, editorID); err != nil {
		return nil, err
	}

	// Proceed with update
	updatedNewsletter, err := s.repo.Update(ctx, newsletterID, &newsletterUpdate)
	if err != nil {
		s.logger.ErrorContext(ctx, "SERVICE: failed to update newsletter", "error", err)
		return nil, err
	}

	return updatedNewsletter, nil
}

// Check if the requesting user is the editor of this newsletter
func (s *NewsletterService) checkNewsletterOwnership(ctx context.Context, newsletter *generated.Newsletter, editorId string) error {
	if newsletter.EditorId.String() != editorId {
		s.logger.WarnContext(ctx, "SERVICE: unauthorized access attempt",
			"requested_editor_id", editorId,
			"newsletter_editor_id", newsletter.EditorId.String())
		return models.NewForbiddenError("You don't have access to this newsletter")
	}

	return nil
}

func (s *NewsletterService) DeleteNewsletter(ctx context.Context, editorID string, newsletterID string) error {
	// Validate input
	if err := s.validateNewsletterID(newsletterID); err != nil {
		return err
	}

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

	if err := s.checkNewsletterOwnership(ctx, newsletter, editorID); err != nil {
		return err
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
