package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"go-newsletter/pkg/generated"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type contextKey string

const validatedUUIDsKey contextKey = "validatedUUIDs"

// UUIDParamValidationMiddleware creates a middleware that validates a UUID parameter from the URL.
// It parses the UUID and stores it in a map in the request context, keyed by its paramName.
// This allows multiple UUIDs to be validated and stored if a route requires it (e.g., /resource/{uuid1}/sub-resource/{uuid2}).
func UUIDParamValidationMiddleware(paramName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawID := chi.URLParam(r, paramName)
			parsedUUID, err := uuid.Parse(rawID)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				apiErr := generated.Error{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("Invalid %s: %s", paramName, err.Error()),
				}
				response := map[string]interface{}{
					"error": map[string]interface{}{
						"code":    apiErr.Code,
						"message": apiErr.Message,
					},
				}
				json.NewEncoder(w).Encode(response)
				return
			}

			// Get or initialize the map of validated UUIDs in the context.
			var validatedUUIDs map[string]uuid.UUID
			existingMap := r.Context().Value(validatedUUIDsKey)
			if existingMap != nil {
				validatedUUIDs, _ = existingMap.(map[string]uuid.UUID)
			}
			if validatedUUIDs == nil {
				validatedUUIDs = make(map[string]uuid.UUID)
			}

			validatedUUIDs[paramName] = parsedUUID
			ctx := context.WithValue(r.Context(), validatedUUIDsKey, validatedUUIDs)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUUIDFromContext retrieves a specific validated UUID from the request context using its original parameter name.
// It returns the UUID and true if found, otherwise uuid.Nil and false.
func GetUUIDFromContext(ctx context.Context, paramName string) (uuid.UUID, bool) {
	validatedUUIDs, ok := ctx.Value(validatedUUIDsKey).(map[string]uuid.UUID)
	if !ok {
		return uuid.Nil, false // Map not found or not the correct type
	}

	// Retrieve the specific UUID from the map.
	id, ok := validatedUUIDs[paramName]
	return id, ok
} 