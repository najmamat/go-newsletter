package handlers

import (
	"log/slog"
	"net/http"

	"go-newsletter/internal/services"
	"go-newsletter/internal/utils"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authService *services.AuthService
	responder   *utils.HTTPResponder
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *services.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		responder:   utils.NewHTTPResponder(logger),
	}
}

// PostAuthPasswordResetRequest handles POST /auth/password-reset endpoint
func (h *AuthHandler) PostAuthPasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Password reset is handled by Supabase Auth on the frontend",
		"instructions": map[string]interface{}{
			"frontend":     "Use supabase.auth.resetPasswordForEmail(email)",
			"documentation": "https://supabase.com/docs/guides/auth/passwords#reset-a-password",
		},
	}
	h.responder.RespondJSON(w, http.StatusOK, response)
}

// PostAuthSignin handles POST /auth/signin endpoint
func (h *AuthHandler) PostAuthSignin(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Sign-in is handled by Supabase Auth on the frontend",
		"instructions": map[string]interface{}{
			"frontend":     "Use supabase.auth.signInWithPassword({ email, password })",
			"documentation": "https://supabase.com/docs/guides/auth/passwords#sign-in-with-password",
			"note":         "After successful sign-in, include the JWT token in the Authorization header for API requests",
		},
	}
	h.responder.RespondJSON(w, http.StatusOK, response)
}

// PostAuthSignup handles POST /auth/signup endpoint
func (h *AuthHandler) PostAuthSignup(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Sign-up is handled by Supabase Auth on the frontend",
		"instructions": map[string]interface{}{
			"frontend":     "Use supabase.auth.signUp({ email, password })",
			"documentation": "https://supabase.com/docs/guides/auth/passwords#sign-up-with-password",
			"note":         "A profile will be automatically created in the profiles table upon successful registration",
		},
	}
	h.responder.RespondJSON(w, http.StatusOK, response)
} 