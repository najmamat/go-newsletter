package handlers

import (
	"encoding/json"
	"go-newsletter/internal/models"
	"go-newsletter/internal/services"
	"go-newsletter/internal/utils"
	"go-newsletter/pkg/generated"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type NewsletterHandler struct {
	service   *services.NewsletterService
	responder *utils.HTTPResponder
}

func NewNewsletterHandler(service *services.NewsletterService, logger *slog.Logger) *NewsletterHandler {
	return &NewsletterHandler{
		service:   service,
		responder: utils.NewHTTPResponder(logger),
	}
}

func (h *NewsletterHandler) GetNewslettersOwnedByEditor(w http.ResponseWriter, r *http.Request) {
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		h.responder.HandleError(w, r, models.NewUnauthorizedError("HANDLER: HANDLER: User not authenticated"))
		return
	}

	newsletters, err := h.service.GetNewslettersOwnedByEditor(r.Context(), user.UserID.String())
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	var newslettersResponse []generated.Newsletter
	for _, newsletter := range newsletters {
		newslettersResponse = append(newslettersResponse, newsletter)
	}

	h.responder.RespondJSON(w, http.StatusOK, newslettersResponse)
}

func (h *NewsletterHandler) GetNewsletterByID(w http.ResponseWriter, r *http.Request) {
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		h.responder.HandleError(w, r, models.NewUnauthorizedError("HANDLER: User not authenticated"))
		return
	}

	newsletterId := chi.URLParam(r, "newsletterId")
	if newsletterId == "" {
		h.responder.HandleError(w, r, models.NewBadRequestError("HANDLER: Newsletter ID is required"))
		return
	}

	newsletter, err := h.service.GetNewsletterByID(r.Context(), newsletterId, user.UserID.String())
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	h.responder.RespondJSON(w, http.StatusOK, newsletter)
}

func (h *NewsletterHandler) PostNewsletters(w http.ResponseWriter, r *http.Request) {
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		h.responder.HandleError(w, r, models.NewUnauthorizedError("HANDLER: User not authenticated"))
		return
	}

	// Decode request body
	var req generated.NewsletterCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.HandleError(w, r, models.NewBadRequestError("HANDLER: Invalid JSON payload"))
		return
	}

	// Create newsletter
	newsletter, err := h.service.CreateNewsletter(r.Context(), user.UserID.String(), req)
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	h.responder.RespondJSON(w, http.StatusCreated, newsletter)
}

func (h *NewsletterHandler) PutNewsletters(w http.ResponseWriter, r *http.Request) {
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		h.responder.HandleError(w, r, models.NewUnauthorizedError("HANDLER: User not authenticated"))
		return
	}

	var req generated.NewsletterUpdate
	// TODO partial update
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.HandleError(w, r, models.NewBadRequestError("HANDLER: Invalid JSON payload"))
		return
	}

	newsletterId := chi.URLParam(r, "newsletterId")
	if newsletterId == "" {
		h.responder.HandleError(w, r, models.NewBadRequestError("HANDLER: Newsletter ID is required"))
		return
	}

	newsletter, err := h.service.UpdateNewsletter(r.Context(), user.UserID.String(), newsletterId, req)
	if err != nil {
		// TODO return 404 when not found (curr. 500)
		h.responder.HandleError(w, r, err)
		return
	}

	h.responder.RespondJSON(w, http.StatusCreated, newsletter)
}
