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
	service        *services.NewsletterService
	profileService *services.ProfileService
	responder      *utils.HTTPResponder
}

func NewNewsletterHandler(service *services.NewsletterService, profileService *services.ProfileService, logger *slog.Logger) *NewsletterHandler {
	return &NewsletterHandler{
		service:        service,
		profileService: profileService,
		responder:      utils.NewHTTPResponder(logger),
	}
}

func (h *NewsletterHandler) GetNewslettersOwnedByEditor(w http.ResponseWriter, r *http.Request) {
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		h.responder.HandleError(w, r, models.NewUnauthorizedError("User not authenticated"))
		return
	}

	newsletters, err := h.service.GetNewslettersOwnedByEditor(r.Context(), user.UserID.String())
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	h.responder.RespondJSON(w, http.StatusOK, newsletters)
}

func (h *NewsletterHandler) GetNewsletterByID(w http.ResponseWriter, r *http.Request) {
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		h.responder.HandleError(w, r, models.NewUnauthorizedError("User not authenticated"))
		return
	}

	newsletterID := chi.URLParam(r, "newsletterId")
	if newsletterID == "" {
		h.responder.HandleError(w, r, models.NewBadRequestError("Newsletter ID is required"))
		return
	}

	newsletter, err := h.service.GetNewsletterByID(r.Context(), newsletterID, user.UserID.String())
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	h.responder.RespondJSON(w, http.StatusOK, newsletter)
}

func (h *NewsletterHandler) PostNewsletters(w http.ResponseWriter, r *http.Request) {
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		h.responder.HandleError(w, r, models.NewUnauthorizedError("User not authenticated"))
		return
	}

	var req generated.NewsletterCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.HandleError(w, r, models.NewBadRequestError("Invalid JSON payload"))
		return
	}

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
		h.responder.HandleError(w, r, models.NewUnauthorizedError("User not authenticated"))
		return
	}

	var req generated.NewsletterUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.HandleError(w, r, models.NewBadRequestError("Invalid JSON payload"))
		return
	}

	newsletterID := chi.URLParam(r, "newsletterId")
	if newsletterID == "" {
		h.responder.HandleError(w, r, models.NewBadRequestError("Newsletter ID is required"))
		return
	}

	newsletter, err := h.service.UpdateNewsletter(r.Context(), user.UserID.String(), newsletterID, req)
	if err != nil {
		if models.IsNotFoundError(err) {
			h.responder.HandleError(w, r, models.NewNotFoundError("Newsletter not found"))
			return
		}
		h.responder.HandleError(w, r, err)
		return
	}

	h.responder.RespondJSON(w, http.StatusOK, newsletter)
}

func (h *NewsletterHandler) DeleteNewsletter(w http.ResponseWriter, r *http.Request) {
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		h.responder.HandleError(w, r, models.NewUnauthorizedError("User not authenticated"))
		return
	}

	newsletterID := chi.URLParam(r, "newsletterId")
	if newsletterID == "" {
		h.responder.HandleError(w, r, models.NewBadRequestError("Newsletter ID is required"))
		return
	}

	if err := h.service.DeleteNewsletter(r.Context(), user.UserID.String(), newsletterID); err != nil {
		if models.IsNotFoundError(err) {
			h.responder.HandleError(w, r, models.NewNotFoundError("Newsletter not found"))
			return
		}
		h.responder.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *NewsletterHandler) GetAllNewsletters(w http.ResponseWriter, r *http.Request) {
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		h.responder.HandleError(w, r, models.NewUnauthorizedError("User not authenticated"))
		return
	}

	profile, err := h.profileService.GetProfileByID(r.Context(), user.UserID.String())
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	if profile.IsAdmin == nil || !*profile.IsAdmin {
		h.responder.HandleError(w, r, models.NewForbiddenError("Admin access required"))
		return
	}

	newsletters, err := h.service.AdminGetAllNewsletters(r.Context())
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	h.responder.RespondJSON(w, http.StatusOK, newsletters)
}

func (h *NewsletterHandler) DeleteNewsletterByID(w http.ResponseWriter, r *http.Request) {
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		h.responder.HandleError(w, r, models.NewUnauthorizedError("User not authenticated"))
		return
	}

	profile, err := h.profileService.GetProfileByID(r.Context(), user.UserID.String())
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	if profile.IsAdmin == nil || !*profile.IsAdmin {
		h.responder.HandleError(w, r, models.NewForbiddenError("Admin access required"))
		return
	}

	newsletterID := chi.URLParam(r, "newsletterId")
	if newsletterID == "" {
		h.responder.HandleError(w, r, models.NewBadRequestError("Newsletter ID is required"))
		return
	}

	if err := h.service.AdminDeleteNewsletterByID(r.Context(), newsletterID); err != nil {
		if models.IsNotFoundError(err) {
			h.responder.HandleError(w, r, models.NewNotFoundError("Newsletter not found"))
			return
		}
		h.responder.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
