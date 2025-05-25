package utils

import (
	"time"

	"go-newsletter/internal/models"
	"go-newsletter/pkg/generated"
)

// ProfileToEditorProfile converts our internal Profile model to the generated EditorProfile
func ProfileToEditorProfile(p models.Profile) generated.EditorProfile {
	return generated.EditorProfile{
		Id:        StringToUUIDPtr(p.ID),
		FullName:  p.FullName,
		AvatarUrl: p.AvatarURL,
		IsAdmin:   &p.IsAdmin,
		CreatedAt: &p.CreatedAt,
		UpdatedAt: &p.UpdatedAt,
	}
}

// EditorProfileToProfile converts generated EditorProfile to our internal Profile model
func EditorProfileToProfile(ep generated.EditorProfile) models.Profile {
	var profile models.Profile
	
	if ep.Id != nil {
		profile.ID = ep.Id.String()
	}
	
	if ep.FullName != nil {
		profile.FullName = ep.FullName
	}
	
	if ep.AvatarUrl != nil {
		profile.AvatarURL = ep.AvatarUrl
	}
	
	if ep.IsAdmin != nil {
		profile.IsAdmin = *ep.IsAdmin
	}
	
	if ep.CreatedAt != nil {
		profile.CreatedAt = *ep.CreatedAt
	} else {
		profile.CreatedAt = time.Now()
	}
	
	if ep.UpdatedAt != nil {
		profile.UpdatedAt = *ep.UpdatedAt
	} else {
		profile.UpdatedAt = time.Now()
	}
	
	return profile
}

// UpdateProfileRequestToInternal converts generated request to our internal type
func UpdateProfileRequestToInternal(req generated.PutMeJSONBody) models.UpdateProfileRequest {
	return models.UpdateProfileRequest{
		FullName:  req.FullName,
		AvatarURL: req.AvatarUrl,
	}
} 