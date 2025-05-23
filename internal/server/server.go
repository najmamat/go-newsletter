package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go-newsletter/internal/models"
	"go-newsletter/internal/services"
	"go-newsletter/internal/utils"
	"go-newsletter/pkg/generated"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Server implements the generated ServerInterface
type Server struct {
	profileService *services.ProfileService
	logger         *slog.Logger
}

// NewServer creates a new server instance
func NewServer(profileService *services.ProfileService, logger *slog.Logger) *Server {
	return &Server{
		profileService: profileService,
		logger:         logger,
	}
}

// GetMe implements the GET /me endpoint
func (s *Server) GetMe(w http.ResponseWriter, r *http.Request) {
	// For now, we'll use a dummy profile ID since we don't have auth yet
	// TODO: Extract user ID from JWT token when auth is implemented
	profiles, err := s.profileService.GetAllProfiles(r.Context())
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	if len(profiles) == 0 {
		s.handleError(w, r, models.NewNotFoundError("No profiles found"))
		return
	}

	// For demo purposes, return the first profile
	profile := utils.ProfileToEditorProfile(profiles[0])
	s.respondJSON(w, http.StatusOK, profile)
}

// PutMe implements the PUT /me endpoint
func (s *Server) PutMe(w http.ResponseWriter, r *http.Request) {
	var req generated.PutMeJSONBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.handleError(w, r, models.NewBadRequestError("Invalid JSON payload"))
		return
	}

	// For now, we'll use a dummy profile ID since we don't have auth yet
	// TODO: Extract user ID from JWT token when auth is implemented
	profiles, err := s.profileService.GetAllProfiles(r.Context())
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	if len(profiles) == 0 {
		s.handleError(w, r, models.NewNotFoundError("No profiles found"))
		return
	}

	// Convert and update the first profile for demo
	updateReq := utils.UpdateProfileRequestToInternal(req)
	updatedProfile, err := s.profileService.UpdateProfile(r.Context(), profiles[0].ID, updateReq)
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

func (s *Server) DeleteAdminNewslettersNewsletterId(w http.ResponseWriter, r *http.Request, newsletterId openapi_types.UUID) {
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

func (s *Server) PutAdminUsersUserIdGrantAdmin(w http.ResponseWriter, r *http.Request, userId openapi_types.UUID) {
	s.notImplemented(w, r)
}

func (s *Server) PutAdminUsersUserIdRevokeAdmin(w http.ResponseWriter, r *http.Request, userId openapi_types.UUID) {
	s.notImplemented(w, r)
}

func (s *Server) PostAuthPasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) PostAuthSignin(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) PostAuthSignup(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) GetNewsletters(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) PostNewsletters(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) DeleteNewslettersNewsletterId(w http.ResponseWriter, r *http.Request, newsletterId openapi_types.UUID) {
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterId(w http.ResponseWriter, r *http.Request, newsletterId openapi_types.UUID) {
	s.notImplemented(w, r)
}

func (s *Server) PutNewslettersNewsletterId(w http.ResponseWriter, r *http.Request, newsletterId openapi_types.UUID) {
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterIdPosts(w http.ResponseWriter, r *http.Request, newsletterId openapi_types.UUID) {
	s.notImplemented(w, r)
}

func (s *Server) PostNewslettersNewsletterIdPosts(w http.ResponseWriter, r *http.Request, newsletterId openapi_types.UUID) {
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterIdScheduledPosts(w http.ResponseWriter, r *http.Request, newsletterId openapi_types.UUID) {
	s.notImplemented(w, r)
}

func (s *Server) DeleteNewslettersNewsletterIdScheduledPostsPostId(w http.ResponseWriter, r *http.Request, newsletterId openapi_types.UUID, postId openapi_types.UUID) {
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterIdScheduledPostsPostId(w http.ResponseWriter, r *http.Request, newsletterId openapi_types.UUID, postId openapi_types.UUID) {
	s.notImplemented(w, r)
}

func (s *Server) PutNewslettersNewsletterIdScheduledPostsPostId(w http.ResponseWriter, r *http.Request, newsletterId openapi_types.UUID, postId openapi_types.UUID) {
	s.notImplemented(w, r)
}

func (s *Server) PostNewslettersNewsletterIdSubscribe(w http.ResponseWriter, r *http.Request, newsletterId openapi_types.UUID) {
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterIdSubscribers(w http.ResponseWriter, r *http.Request, newsletterId openapi_types.UUID) {
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