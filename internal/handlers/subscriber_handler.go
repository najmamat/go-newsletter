package handlers

import (
	"encoding/json"
	"go-newsletter/internal/models"
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
		h.responder.HandleError(w, r, models.NewBadRequestError("Invalid newsletter ID"))
		return
	}

	// Get user from context
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		h.responder.HandleError(w, r, models.NewUnauthorizedError("User not authenticated"))
		return
	}

	// Get subscribers
	subscribers, err := h.subscriberService.ListSubscribers(r.Context(), newsletterID, user.UserID.String())
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	h.responder.RespondJSON(w, http.StatusOK, subscribers)
}

// Subscribe handles POST /newsletters/{newsletterId}/subscribe
func (h *SubscriberHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	newsletterID, err := uuid.Parse(chi.URLParam(r, "newsletterId"))
	if err != nil {
		h.responder.HandleError(w, r, models.NewBadRequestError("Invalid newsletter ID"))
		return
	}

	var req generated.SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.HandleError(w, r, models.NewBadRequestError("Invalid request body"))
		return
	}

	_, err = h.subscriberService.Subscribe(r.Context(), newsletterID, req.Email)
	if err != nil {
		h.responder.HandleError(w, r, err)
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
		h.responder.HandleError(w, r, err)
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
		h.responder.HandleError(w, r, err)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Successfully unsubscribed from the newsletter",
	}

	h.responder.RespondJSON(w, http.StatusOK, response)
}
