package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go-newsletter/internal/models"
	"go-newsletter/internal/services"

	"github.com/go-chi/chi/v5"
)

// ProfileHandler handles HTTP requests for profiles
type ProfileHandler struct {
	service *services.ProfileService
	logger  *slog.Logger
}

// NewProfileHandler creates a new ProfileHandler
func NewProfileHandler(service *services.ProfileService, logger *slog.Logger) *ProfileHandler {
	return &ProfileHandler{
		service: service,
		logger:  logger,
	}
}

// GetAllProfiles handles GET /profiles
func (h *ProfileHandler) GetAllProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.service.GetAllProfiles(r.Context())
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, profiles)
}

// GetProfileByID handles GET /profiles/{id}
func (h *ProfileHandler) GetProfileByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	profile, err := h.service.GetProfileByID(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, profile)
}

// UpdateProfile handles PUT /profiles/{id}
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, models.NewBadRequestError("Invalid JSON payload"))
		return
	}

	profile, err := h.service.UpdateProfile(r.Context(), id, req)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, profile)
}

// respondJSON sends a JSON response
func (h *ProfileHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", "error", err)
	}
}

// handleError handles API errors and sends appropriate responses
func (h *ProfileHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if apiErr, ok := err.(models.APIError); ok {
		h.logger.WarnContext(r.Context(), "API error", "code", apiErr.Code, "message", apiErr.Message)
		h.respondJSON(w, apiErr.Code, apiErr)
		return
	}

	// For unexpected errors, log them and return a generic 500
	h.logger.ErrorContext(r.Context(), "Unexpected error", "error", err)
	h.respondJSON(w, http.StatusInternalServerError, models.NewInternalServerError("An unexpected error occurred"))
} 