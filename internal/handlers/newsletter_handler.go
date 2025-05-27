package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go-newsletter/internal/models"
	"go-newsletter/internal/services"
	"go-newsletter/pkg/generated"
)

// NewsletterHandler handles HTTP requests for newsletters
type NewsletterHandler struct {
	service *services.NewsletterService
	logger  *slog.Logger
}

// NewNewsletterHandler creates a new NewsletterHandler
func NewNewsletterHandler(service *services.NewsletterService, logger *slog.Logger) *NewsletterHandler {
	return &NewsletterHandler{
		service: service,
		logger:  logger,
	}
}

// GetAllNewsletters handles GET /newsletters
func (h *NewsletterHandler) GetAllNewsletters(w http.ResponseWriter, r *http.Request) {
	newsletters, err := h.service.GetAllNewsletters(r.Context())
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, newsletters)
}

// CreateNewsletter handles POST /newsletters
func (h *NewsletterHandler) CreateNewsletter(w http.ResponseWriter, r *http.Request) {
	var req generated.NewsletterCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err)
		return
	}

	// TODO: Replace with real auth once implemented
	editorID := "test-editor-id"

	newsletter, err := h.service.CreateNewsletter(r.Context(), editorID, req)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, newsletter)
}

// respondJSON sends a JSON response
func (h *NewsletterHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", "error", err)
	}
}

// handleError handles API errors and sends appropriate responses
func (h *NewsletterHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if apiErr, ok := err.(models.APIError); ok {
		h.logger.WarnContext(r.Context(), "API error", "code", apiErr.Code, "message", apiErr.Message)
		h.respondJSON(w, apiErr.Code, apiErr)
		return
	}

	// For unexpected errors, log them and return a generic 500
	h.logger.ErrorContext(r.Context(), "Unexpected error", "error", err)
	h.respondJSON(w, http.StatusInternalServerError, models.NewInternalServerError("An unexpected error occurred"))
}
