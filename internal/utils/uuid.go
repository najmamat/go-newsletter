package utils

import (
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// StringToUUIDPtr converts a string to a UUID pointer
func StringToUUIDPtr(s string) *openapi_types.UUID {
	if s == "" {
		return nil
	}
	
	parsed, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	
	result := openapi_types.UUID(parsed)
	return &result
}

// UUIDPtrToString converts a UUID pointer to string
func UUIDPtrToString(u *openapi_types.UUID) string {
	if u == nil {
		return ""
	}
	return uuid.UUID(*u).String()
} 