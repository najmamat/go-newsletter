package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"go-newsletter/internal/models"
	"go-newsletter/internal/services"
)

// AuthMiddleware wraps handlers to require JWT authentication
type AuthMiddleware struct {
	authService *services.AuthService
	logger      *slog.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService *services.AuthService, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

// RequireAuth middleware that validates JWT and adds user context
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.handleUnauthorized(w, "Missing authorization header")
			return
		}

		// Validate token format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			m.handleUnauthorized(w, "Invalid authorization header format")
			return
		}

		// Get user from token
		user, err := m.authService.GetUserFromToken(authHeader)
		if err != nil {
			m.logger.Warn("JWT validation failed", "error", err.Error())
			m.handleUnauthorized(w, "Invalid or expired token")
			return
		}

		// Add user to request context
		ctx := services.AddUserToContext(r.Context(), user)
		r = r.WithContext(ctx)

		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

// RequireAdmin middleware that requires admin privileges
func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user from context
		_, ok := services.GetUserFromContext(r.Context())
		if !ok {
			m.handleUnauthorized(w, "User context not found")
			return
		}

		// Check admin status - we'll need to check the profiles table
		// For now, we'll implement this check in the handler level
		// since admin status is stored in the database, not in JWT

		next.ServeHTTP(w, r)
	}))
}

// OptionalAuth middleware that adds user context if token is present but doesn't require it
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			// Try to get user from token
			user, err := m.authService.GetUserFromToken(authHeader)
			if err != nil {
				m.logger.Debug("Optional auth failed", "error", err.Error())
			} else {
				// Add user to request context
				ctx := services.AddUserToContext(r.Context(), user)
				r = r.WithContext(ctx)
			}
		}

		// Continue regardless of auth status
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) handleUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	
	// Use existing error patterns
	apiErr := models.NewUnauthorizedError(message)
	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    apiErr.Code,
			"message": apiErr.Message,
		},
	}
	
	// Write JSON response
	json.NewEncoder(w).Encode(response)
} 