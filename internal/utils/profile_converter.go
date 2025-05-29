package utils

import (
	"go-newsletter/internal/models"
	"go-newsletter/pkg/generated"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// EditorProfileToProfile converts a generated.EditorProfile to models.Profile
func EditorProfileToProfile(ep generated.EditorProfile) models.Profile {
	var fullName, avatarURL string
	if ep.FullName != nil {
		fullName = *ep.FullName
	}
	if ep.AvatarUrl != nil {
		avatarURL = *ep.AvatarUrl
	}

	return models.Profile{
		ID:        ep.Id.String(),
		FullName:  fullName,
		AvatarURL: avatarURL,
		IsAdmin:   *ep.IsAdmin,
		CreatedAt: *ep.CreatedAt,
		UpdatedAt: *ep.UpdatedAt,
	}
}

// InternalProfileToEditorProfile converts a models.Profile to generated.EditorProfile
func InternalProfileToEditorProfile(p *models.Profile) generated.EditorProfile {
	id := openapi_types.UUID{}
	if err := id.UnmarshalText([]byte(p.ID)); err != nil {
		// Handle error if needed
	}
	fullName := p.FullName
	avatarUrl := p.AvatarURL
	isAdmin := p.IsAdmin
	createdAt := p.CreatedAt
	updatedAt := p.UpdatedAt

	return generated.EditorProfile{
		Id:        &id,
		FullName:  &fullName,
		AvatarUrl: &avatarUrl,
		IsAdmin:   &isAdmin,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}
} 