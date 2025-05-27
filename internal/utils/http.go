package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go-newsletter/internal/models"
	"go-newsletter/pkg/generated"
)

// HTTPResponder provides common HTTP response functionality
type HTTPResponder struct {
	Logger *slog.Logger
}

// NewHTTPResponder creates a new HTTPResponder
func NewHTTPResponder(logger *slog.Logger) *HTTPResponder {
	return &HTTPResponder{
		Logger: logger,
	}
}

// RespondJSON sends a JSON response
func (h *HTTPResponder) RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.Logger.Error("Failed to encode JSON response", "error", err)
	}
}

// HandleError handles API errors and sends appropriate responses
func (h *HTTPResponder) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if apiErr, ok := err.(models.APIError); ok {
		h.Logger.WarnContext(r.Context(), "API error", "code", apiErr.Code, "message", apiErr.Message)
		errorResponse := generated.Error{
			Code:    int32(apiErr.Code),
			Message: apiErr.Message,
		}
		h.RespondJSON(w, apiErr.Code, errorResponse)
		return
	}

	// For unexpected errors, log them and return a generic 500
	h.Logger.ErrorContext(r.Context(), "Unexpected error", "error", err)
	errorResponse := generated.Error{
		Code:    500,
		Message: "An unexpected error occurred",
	}
	h.RespondJSON(w, http.StatusInternalServerError, errorResponse)
} 