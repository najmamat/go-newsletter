package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go-newsletter/internal/models"
	"go-newsletter/internal/services"
	"go-newsletter/internal/utils"
	"go-newsletter/pkg/generated"

	"github.com/go-chi/chi/v5"
)

// ProfileHandler handles HTTP requests for profiles
type ProfileHandler struct {
	service     *services.ProfileService
	authService *services.AuthService
	responder   *utils.HTTPResponder
}

// NewProfileHandler creates a new ProfileHandler
func NewProfileHandler(service *services.ProfileService, authService *services.AuthService, logger *slog.Logger) *ProfileHandler {
	return &ProfileHandler{
		service:     service,
		authService: authService,
		responder:   utils.NewHTTPResponder(logger),
	}
}

// GetMe handles GET /me endpoint
func (h *ProfileHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		h.responder.HandleError(w, r, models.NewUnauthorizedError("User not authenticated"))
		return
	}

	// Get the user's profile from database
	profile, err := h.service.GetProfileByID(r.Context(), user.UserID.String())
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	// Convert to API response format
	editorProfile := utils.ProfileToEditorProfile(*profile)
	h.responder.RespondJSON(w, http.StatusOK, editorProfile)
}

// GetAllProfiles handles GET /profiles
func (h *ProfileHandler) GetAllProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.service.GetAllProfiles(r.Context())
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	var editorProfiles []generated.EditorProfile
	for _, profile := range profiles {
		editorProfiles = append(editorProfiles, utils.ProfileToEditorProfile(profile))
	}

	h.responder.RespondJSON(w, http.StatusOK, editorProfiles)
}

// GetProfileByID handles GET /profiles/{id}
func (h *ProfileHandler) GetProfileByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	profile, err := h.service.GetProfileByID(r.Context(), id)
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	editorProfile := utils.ProfileToEditorProfile(*profile)
	h.responder.RespondJSON(w, http.StatusOK, editorProfile)
}

// UpdateProfile handles PUT /profiles/{id}
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	var req generated.PutMeJSONBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.HandleError(w, r, models.NewBadRequestError("Invalid JSON payload"))
		return
	}

	profile, err := h.service.UpdateProfile(r.Context(), id, req)
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	editorProfile := utils.ProfileToEditorProfile(*profile)
	h.responder.RespondJSON(w, http.StatusOK, editorProfile)
}

// PutMe handles PUT /me endpoint
func (h *ProfileHandler) PutMe(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		h.responder.HandleError(w, r, models.NewUnauthorizedError("User not authenticated"))
		return
	}

	var req generated.PutMeJSONBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.HandleError(w, r, models.NewBadRequestError("Invalid JSON payload"))
		return
	}

	// Convert and update the user's profile
	updateReq := utils.UpdateProfileRequestToInternal(req)
	updatedProfile, err := h.service.UpdateProfile(r.Context(), user.UserID.String(), updateReq)
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	profile := utils.ProfileToEditorProfile(*updatedProfile)
	h.responder.RespondJSON(w, http.StatusOK, profile)
} 