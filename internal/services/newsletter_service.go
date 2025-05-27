package services

import (
	"context"
	"errors"
	"go-newsletter/internal/repository"
	"go-newsletter/pkg/generated"
	"log/slog"
	"strings"
)

var (
	ErrInvalidNewsletterName = errors.New("newsletter name is required")
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
func (s *NewsletterService) GetAllNewsletters(ctx context.Context) ([]generated.Newsletter, error) {
	newsletters, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get newsletters", "error", err)
		return nil, err
	}
	return newsletters, nil
}

func (s *NewsletterService) CreateNewsletter(ctx context.Context, editorID string, req generated.NewsletterCreate) (*generated.Newsletter, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrInvalidNewsletterName
	}

	newsletter, err := s.repo.Insert(ctx, editorID, req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create newsletter", "error", err)
		return nil, err
	}

	return newsletter, nil
}
