package services

import (
	"context"
	"go-newsletter/internal/models"
	"go-newsletter/internal/repository"
	"log/slog"
	"strings"
)

// NewsletterServce handles business logic for newsletters
type NewsletterService struct {
	repo   *repository.NewsletterRepository
	logger *slog.Logger
}

// NewNewsletterService creates a new NewsletterService
func NewNewsletterService(repo *repository.NewsletterRepository, logger *slog.Logger) *NewsletterService {
	return &NewsletterService{
		repo:   repo,
		logger: logger,
	}
}

// GetAllNewsletters retrieves all newsletters
func (s *NewsletterService) GetAllNewsletters(ctx context.Context) ([]models.Newsletter, error) {
	newsletters, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get all newsletters", "error", err)
		return nil, models.NewInternalServerError("Failed to get all newsletters")
	}
	return newsletters, nil
}

func (s *NewsletterService) CreateNewsletter(ctx context.Context, editorID string, req models.NewsletterCreateRequest) (*models.Newsletter, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, models.NewBadRequestError("Name is required")
	}
	return s.repo.Insert(ctx, editorID, req)
}
