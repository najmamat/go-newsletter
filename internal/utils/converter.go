package utils

import (
	"go-newsletter/pkg/generated"
)

func ProfileToEditorProfile(p generated.EditorProfile) generated.EditorProfile {
	return generated.EditorProfile{
		Id:        p.Id,
		FullName:  p.FullName,
		AvatarUrl: p.AvatarUrl,
		Email:     p.Email,
		IsAdmin:   p.IsAdmin,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// UpdateProfileRequestToInternal converts generated request to our internal type.
func UpdateProfileRequestToInternal(req generated.PutMeJSONBody) generated.PutMeJSONBody {
	return generated.PutMeJSONBody{
		FullName:  req.FullName,
		AvatarUrl: req.AvatarUrl,
	}
} 