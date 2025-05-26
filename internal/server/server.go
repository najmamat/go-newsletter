package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go-newsletter/internal/middleware"
	"go-newsletter/internal/models"
	"go-newsletter/internal/services"
	"go-newsletter/internal/utils"
	"go-newsletter/pkg/generated"
)

// Server implements the generated ServerInterface
type Server struct {
	profileService *services.ProfileService
	authService    *services.AuthService
	logger         *slog.Logger
}

// NewServer creates a new server instance
func NewServer(profileService *services.ProfileService, authService *services.AuthService, logger *slog.Logger) *Server {
	return &Server{
		profileService: profileService,
		authService:    authService,
		logger:         logger,
	}
}

// GetAuthService returns the auth service instance
func (s *Server) GetAuthService() *services.AuthService {
	return s.authService
}

// GetMe implements the GET /me endpoint
func (s *Server) GetMe(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		s.handleError(w, r, models.NewUnauthorizedError("User not authenticated"))
		return
	}

	// Get the user's profile from database
	profile, err := s.profileService.GetProfileByID(r.Context(), user.UserID.String())
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	// Convert to API response format
	editorProfile := utils.ProfileToEditorProfile(*profile)
	s.respondJSON(w, http.StatusOK, editorProfile)
}

// PutMe implements the PUT /me endpoint
func (s *Server) PutMe(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		s.handleError(w, r, models.NewUnauthorizedError("User not authenticated"))
		return
	}

	var req generated.PutMeJSONBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.handleError(w, r, models.NewBadRequestError("Invalid JSON payload"))
		return
	}

	// Convert and update the user's profile
	updateReq := utils.UpdateProfileRequestToInternal(req)
	updatedProfile, err := s.profileService.UpdateProfile(r.Context(), user.UserID.String(), updateReq)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	profile := utils.ProfileToEditorProfile(*updatedProfile)
	s.respondJSON(w, http.StatusOK, profile)
}

// Placeholder implementations for other endpoints (will implement as needed)
func (s *Server) GetAdminNewsletters(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) DeleteAdminNewslettersNewsletterId(w http.ResponseWriter, r *http.Request) {
	newsletterId, ok := middleware.GetUUIDFromContext(r.Context(), "newsletterId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("newsletterId not found in context"))
		return
	}
	s.logger.Info("DeleteAdminNewslettersNewsletterId called", "newsletterId", newsletterId)
	s.notImplemented(w, r)
}

func (s *Server) GetAdminUsers(w http.ResponseWriter, r *http.Request) {
	// This one we can implement since we have profile service
	profiles, err := s.profileService.GetAllProfiles(r.Context())
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	var editorProfiles []generated.EditorProfile
	for _, profile := range profiles {
		editorProfiles = append(editorProfiles, utils.ProfileToEditorProfile(profile))
	}

	s.respondJSON(w, http.StatusOK, editorProfiles)
}

func (s *Server) PutAdminUsersUserIdGrantAdmin(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUUIDFromContext(r.Context(), "userId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("userId not found in context"))
		return
	}
	s.logger.Info("PutAdminUsersUserIdGrantAdmin called", "userId", userId)
	s.notImplemented(w, r)
}

func (s *Server) PutAdminUsersUserIdRevokeAdmin(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUUIDFromContext(r.Context(), "userId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("userId not found in context"))
		return
	}
	s.logger.Info("PutAdminUsersUserIdRevokeAdmin called", "userId", userId)
	s.notImplemented(w, r)
}

func (s *Server) DeleteAdminUsersUserId(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUUIDFromContext(r.Context(), "userId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("userId not found in context"))
		return
	}
	s.logger.Info("DeleteAdminUsersUserId called", "userId", userId)
	s.notImplemented(w, r)
}

func (s *Server) PostAuthPasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Password reset is handled by Supabase Auth on the frontend",
		"instructions": map[string]interface{}{
			"frontend": "Use supabase.auth.resetPasswordForEmail(email)",
			"documentation": "https://supabase.com/docs/guides/auth/passwords#reset-a-password",
		},
	}
	s.respondJSON(w, http.StatusOK, response)
}

func (s *Server) PostAuthSignin(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Sign-in is handled by Supabase Auth on the frontend",
		"instructions": map[string]interface{}{
			"frontend": "Use supabase.auth.signInWithPassword({ email, password })",
			"documentation": "https://supabase.com/docs/guides/auth/passwords#sign-in-with-password",
			"note": "After successful sign-in, include the JWT token in the Authorization header for API requests",
		},
	}
	s.respondJSON(w, http.StatusOK, response)
}

func (s *Server) PostAuthSignup(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Sign-up is handled by Supabase Auth on the frontend",
		"instructions": map[string]interface{}{
			"frontend": "Use supabase.auth.signUp({ email, password })",
			"documentation": "https://supabase.com/docs/guides/auth/passwords#sign-up-with-password",
			"note": "A profile will be automatically created in the profiles table upon successful registration",
		},
	}
	s.respondJSON(w, http.StatusOK, response)
}

func (s *Server) GetNewsletters(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) PostNewsletters(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) DeleteNewslettersNewsletterId(w http.ResponseWriter, r *http.Request) {
	newsletterId, ok := middleware.GetUUIDFromContext(r.Context(), "newsletterId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("newsletterId not found in context"))
		return
	}
	s.logger.Info("DeleteNewslettersNewsletterId called", "newsletterId", newsletterId)
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterId(w http.ResponseWriter, r *http.Request) {
	newsletterId, ok := middleware.GetUUIDFromContext(r.Context(), "newsletterId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("newsletterId not found in context"))
		return
	}
	s.logger.Info("GetNewslettersNewsletterId called", "newsletterId", newsletterId)
	s.notImplemented(w, r)
}

func (s *Server) PutNewslettersNewsletterId(w http.ResponseWriter, r *http.Request) {
	newsletterId, ok := middleware.GetUUIDFromContext(r.Context(), "newsletterId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("newsletterId not found in context"))
		return
	}
	s.logger.Info("PutNewslettersNewsletterId called", "newsletterId", newsletterId)
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterIdPosts(w http.ResponseWriter, r *http.Request) {
	newsletterId, ok := middleware.GetUUIDFromContext(r.Context(), "newsletterId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("newsletterId not found in context"))
		return
	}
	s.logger.Info("GetNewslettersNewsletterIdPosts called", "newsletterId", newsletterId)
	s.notImplemented(w, r)
}

func (s *Server) PostNewslettersNewsletterIdPosts(w http.ResponseWriter, r *http.Request) {
	newsletterId, ok := middleware.GetUUIDFromContext(r.Context(), "newsletterId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("newsletterId not found in context"))
		return
	}
	s.logger.Info("PostNewslettersNewsletterIdPosts called", "newsletterId", newsletterId)
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterIdScheduledPosts(w http.ResponseWriter, r *http.Request) {
	newsletterId, ok := middleware.GetUUIDFromContext(r.Context(), "newsletterId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("newsletterId not found in context"))
		return
	}
	s.logger.Info("GetNewslettersNewsletterIdScheduledPosts called", "newsletterId", newsletterId)
	s.notImplemented(w, r)
}

func (s *Server) DeleteNewslettersNewsletterIdScheduledPostsPostId(w http.ResponseWriter, r *http.Request) {
	newsletterId, ok := middleware.GetUUIDFromContext(r.Context(), "newsletterId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("newsletterId not found in context"))
		return
	}
	postId, ok := middleware.GetUUIDFromContext(r.Context(), "postId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("postId not found in context"))
		return
	}
	s.logger.Info("DeleteNewslettersNewsletterIdScheduledPostsPostId called", "newsletterId", newsletterId, "postId", postId)
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterIdScheduledPostsPostId(w http.ResponseWriter, r *http.Request) {
	newsletterId, ok := middleware.GetUUIDFromContext(r.Context(), "newsletterId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("newsletterId not found in context"))
		return
	}
	postId, ok := middleware.GetUUIDFromContext(r.Context(), "postId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("postId not found in context"))
		return
	}
	s.logger.Info("GetNewslettersNewsletterIdScheduledPostsPostId called", "newsletterId", newsletterId, "postId", postId)
	s.notImplemented(w, r)
}

func (s *Server) PutNewslettersNewsletterIdScheduledPostsPostId(w http.ResponseWriter, r *http.Request) {
	newsletterId, ok := middleware.GetUUIDFromContext(r.Context(), "newsletterId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("newsletterId not found in context"))
		return
	}
	postId, ok := middleware.GetUUIDFromContext(r.Context(), "postId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("postId not found in context"))
		return
	}
	s.logger.Info("PutNewslettersNewsletterIdScheduledPostsPostId called", "newsletterId", newsletterId, "postId", postId)
	s.notImplemented(w, r)
}

func (s *Server) PostNewslettersNewsletterIdSubscribe(w http.ResponseWriter, r *http.Request) {
	newsletterId, ok := middleware.GetUUIDFromContext(r.Context(), "newsletterId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("newsletterId not found in context"))
		return
	}
	s.logger.Info("PostNewslettersNewsletterIdSubscribe called", "newsletterId", newsletterId)
	s.notImplemented(w, r)
}

func (s *Server) PostNewslettersNewsletterIdUnsubscribe(w http.ResponseWriter, r *http.Request) {
	newsletterId, ok := middleware.GetUUIDFromContext(r.Context(), "newsletterId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("newsletterId not found in context"))
		return
	}
	s.logger.Info("PostNewslettersNewsletterIdUnsubscribe called", "newsletterId", newsletterId)
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterIdConfirmSubscription(w http.ResponseWriter, r *http.Request) {
	newsletterId, ok := middleware.GetUUIDFromContext(r.Context(), "newsletterId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("newsletterId not found in context"))
		return
	}
	s.logger.Info("GetNewslettersNewsletterIdConfirmSubscription called", "newsletterId", newsletterId)
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterIdSubscribers(w http.ResponseWriter, r *http.Request) {
	newsletterId, ok := middleware.GetUUIDFromContext(r.Context(), "newsletterId")
	if !ok {
		s.handleError(w, r, models.NewBadRequestError("newsletterId not found in context"))
		return
	}
	s.logger.Info("GetNewslettersNewsletterIdSubscribers called", "newsletterId", newsletterId)
	s.notImplemented(w, r)
}

func (s *Server) GetSubscribeConfirmConfirmationToken(w http.ResponseWriter, r *http.Request, confirmationToken string) {
	s.notImplemented(w, r)
}

func (s *Server) GetUnsubscribeUnsubscribeToken(w http.ResponseWriter, r *http.Request, unsubscribeToken string) {
	s.notImplemented(w, r)
}

// Helper methods
func (s *Server) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Error("Failed to encode JSON response", "error", err)
	}
}

func (s *Server) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if apiErr, ok := err.(models.APIError); ok {
		s.logger.WarnContext(r.Context(), "API error", "code", apiErr.Code, "message", apiErr.Message)
		errorResponse := generated.Error{
			Code:    int32(apiErr.Code),
			Message: apiErr.Message,
		}
		s.respondJSON(w, apiErr.Code, errorResponse)
		return
	}

	// For unexpected errors, log them and return a generic 500
	s.logger.ErrorContext(r.Context(), "Unexpected error", "error", err)
	errorResponse := generated.Error{
		Code:    500,
		Message: "An unexpected error occurred",
	}
	s.respondJSON(w, http.StatusInternalServerError, errorResponse)
}

func (s *Server) notImplemented(w http.ResponseWriter, r *http.Request) {
	errorResponse := generated.Error{
		Code:    501,
		Message: "Endpoint not yet implemented",
	}
	s.respondJSON(w, http.StatusNotImplemented, errorResponse)
} 