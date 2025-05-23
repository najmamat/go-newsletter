package services

import (
	"context"
	"errors"
	"log/slog"

	"go-newsletter/internal/models"
	"go-newsletter/internal/repository"

	"github.com/jackc/pgx/v5"
)

// ProfileService handles business logic for profiles
type ProfileService struct {
	repo   *repository.ProfileRepository
	logger *slog.Logger
}

// NewProfileService creates a new ProfileService
func NewProfileService(repo *repository.ProfileRepository, logger *slog.Logger) *ProfileService {
	return &ProfileService{
		repo:   repo,
		logger: logger,
	}
}

// GetAllProfiles retrieves all profiles
func (s *ProfileService) GetAllProfiles(ctx context.Context) ([]models.Profile, error) {
	profiles, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get all profiles", "error", err)
		return nil, models.NewInternalServerError("Failed to retrieve profiles")
	}
	return profiles, nil
}

// GetProfileByID retrieves a profile by ID
func (s *ProfileService) GetProfileByID(ctx context.Context, id string) (*models.Profile, error) {
	if id == "" {
		return nil, models.NewBadRequestError("Profile ID is required")
	}

	profile, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.NewNotFoundError("Profile not found")
		}
		s.logger.ErrorContext(ctx, "Failed to get profile by ID", "id", id, "error", err)
		return nil, models.NewInternalServerError("Failed to retrieve profile")
	}

	return profile, nil
}

// UpdateProfile updates a profile
func (s *ProfileService) UpdateProfile(ctx context.Context, id string, req models.UpdateProfileRequest) (*models.Profile, error) {
	if id == "" {
		return nil, models.NewBadRequestError("Profile ID is required")
	}

	// First check if profile exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.NewNotFoundError("Profile not found")
		}
		s.logger.ErrorContext(ctx, "Failed to check profile existence", "id", id, "error", err)
		return nil, models.NewInternalServerError("Failed to update profile")
	}

	// Update the profile
	updatedProfile, err := s.repo.Update(ctx, id, req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to update profile", "id", id, "error", err)
		return nil, models.NewInternalServerError("Failed to update profile")
	}

	return updatedProfile, nil
} 