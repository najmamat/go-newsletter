package handlers

import (
	"encoding/json"
	"errors"
	"go-newsletter/internal/utils"
	"net/http"

	"go-newsletter/internal/services"
	"go-newsletter/pkg/generated"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SubscriberHandler struct {
	subscriberService *services.SubscriberService
	responder         *utils.HTTPResponder
}

func NewSubscriberHandler(subscriberService *services.SubscriberService, responder *utils.HTTPResponder) *SubscriberHandler {
	return &SubscriberHandler{
		subscriberService: subscriberService,
		responder:         responder,
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
		switch {
		case errors.Is(err, services.ErrNotFound):
			http.Error(w, "Newsletter not found", http.StatusNotFound)
		case errors.Is(err, services.ErrForbidden):
			http.Error(w, "You don't have permission to access this newsletter", http.StatusForbidden)
		default:
			http.Error(w, "Failed to list subscribers", http.StatusInternalServerError)
		}
		return
	}

	h.responder.RespondJSON(w, http.StatusOK, subscribers)
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

	_, err = h.subscriberService.Subscribe(r.Context(), newsletterID, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrNotFound):
			http.Error(w, "Newsletter not found", http.StatusNotFound)
		case errors.Is(err, services.ErrAlreadySubscribed):
			http.Error(w, "Already subscribed to this newsletter", http.StatusConflict)
		default:
			http.Error(w, "Failed to subscribe to newsletter", http.StatusInternalServerError)
		}
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Subscription successful. Please check your email to confirm your subscription.",
	}

	h.responder.RespondJSON(w, http.StatusOK, response)
}

// ConfirmSubscription handles the confirmation of a subscription using a token
func (h *SubscriberHandler) ConfirmSubscription(w http.ResponseWriter, r *http.Request, confirmationToken string) {
	err := h.subscriberService.ConfirmSubscription(r.Context(), confirmationToken)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			http.Error(w, "Invalid or expired confirmation token", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to confirm subscription", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Subscription confirmed successfully",
	}

	h.responder.RespondJSON(w, http.StatusOK, response)
}

// Unsubscribe handles the unsubscription using a token
func (h *SubscriberHandler) Unsubscribe(w http.ResponseWriter, r *http.Request, unsubscribeToken string) {
	err := h.subscriberService.Unsubscribe(r.Context(), unsubscribeToken)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			http.Error(w, "Invalid or expired unsubscribe token", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to unsubscribe", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Successfully unsubscribed from the newsletter",
	}

	h.responder.RespondJSON(w, http.StatusOK, response)
}
