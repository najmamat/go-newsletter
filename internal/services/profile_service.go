package services

import (
	"context"
	"log/slog"

	"go-newsletter/internal/models"
	"go-newsletter/internal/repository"
	"go-newsletter/internal/utils"
	"go-newsletter/pkg/generated"

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
func (s *ProfileService) GetAllProfiles(ctx context.Context) ([]generated.EditorProfile, error) {
	profiles, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []generated.EditorProfile
	for _, p := range profiles {
		result = append(result, utils.ProfileToEditorProfile(p))
	}
	return result, nil
}

// GetProfileByID retrieves a profile by ID
func (s *ProfileService) GetProfileByID(ctx context.Context, id string) (*generated.EditorProfile, error) {
	profile, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.NewNotFoundError("Profile not found")
		}
		return nil, err
	}
	result := utils.ProfileToEditorProfile(*profile)
	return &result, nil
}

// UpdateProfile updates a profile
func (s *ProfileService) UpdateProfile(ctx context.Context, id string, req generated.PutMeJSONBody) (*generated.EditorProfile, error) {
	// Check if profile exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.NewNotFoundError("Profile not found")
		}
		return nil, err
	}

	// Update profile
	updatedProfile, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	result := utils.ProfileToEditorProfile(*updatedProfile)
	return &result, nil
}

// GrantAdmin grants admin privileges to a user
func (s *ProfileService) GrantAdmin(ctx context.Context, id string) (*generated.EditorProfile, error) {
	// Check if profile exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.NewNotFoundError("Profile not found")
		}
		return nil, err
	}

	profile, err := s.repo.GrantAdmin(ctx, id)
	if err != nil {
		return nil, err
	}
	result := utils.ProfileToEditorProfile(*profile)
	return &result, nil
}

// RevokeAdmin revokes admin privileges from a user
func (s *ProfileService) RevokeAdmin(ctx context.Context, id string) (*generated.EditorProfile, error) {
	// Check if profile exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.NewNotFoundError("Profile not found")
		}
		return nil, err
	}

	profile, err := s.repo.RevokeAdmin(ctx, id)
	if err != nil {
		return nil, err
	}
	result := utils.ProfileToEditorProfile(*profile)
	return &result, nil
} 