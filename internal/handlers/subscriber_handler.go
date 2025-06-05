package handlers

import (
	"encoding/json"
	"net/http"

	"go-newsletter/internal/services"
	"go-newsletter/pkg/generated"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SubscriberHandler struct {
	subscriberService *services.SubscriberService
}

func NewSubscriberHandler(subscriberService *services.SubscriberService) *SubscriberHandler {
	return &SubscriberHandler{
		subscriberService: subscriberService,
	}
}

// ListSubscribers handles GET /newsletters/{newsletterId}/subscribers
func (h *SubscriberHandler) ListSubscribers(w http.ResponseWriter, r *http.Request) {
	newsletterID, err := uuid.Parse(chi.URLParam(r, "newsletterId"))
	if err != nil {
		http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
		return
	}

	// Get user from context
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get subscribers
	subscribers, err := h.subscriberService.ListSubscribers(r.Context(), newsletterID, user.UserID.String())
	if err != nil {
		switch err {
		case services.ErrNotFound:
			http.Error(w, "Newsletter not found", http.StatusNotFound)
		case services.ErrForbidden:
			http.Error(w, "You don't have permission to access this newsletter", http.StatusForbidden)
		default:
			http.Error(w, "Failed to list subscribers", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"subscribers": subscribers,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Subscribe handles POST /newsletters/{newsletterId}/subscribe
func (h *SubscriberHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	newsletterID, err := uuid.Parse(chi.URLParam(r, "newsletterId"))
	if err != nil {
		http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
		return
	}

	var req generated.SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	subscriber, err := h.subscriberService.Subscribe(r.Context(), newsletterID, req.Email)
	if err != nil {
		switch err {
		case services.ErrNotFound:
			http.Error(w, "Newsletter not found", http.StatusNotFound)
		case services.ErrAlreadySubscribed:
			http.Error(w, "Already subscribed to this newsletter", http.StatusConflict)
		default:
			http.Error(w, "Failed to subscribe to newsletter", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"message":    "Subscription successful. Please check your email to confirm your subscription.",
		"subscriber": subscriber,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
} 