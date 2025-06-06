package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go-newsletter/internal/models"
	"go-newsletter/internal/services"
	"go-newsletter/internal/utils"
	"go-newsletter/pkg/generated"
	"net/http"
)

type PostHandler struct {
	postService *services.PostService
	responder   *utils.HTTPResponder
}

func NewPostHandler(postService *services.PostService, responder *utils.HTTPResponder) *PostHandler {
	return &PostHandler{
		postService: postService,
		responder:   responder,
	}
}

// GetPostsByNewsletterId handles GET /newsletters/{newsletterId}/posts and GET /newsletters/{newsletterId}/published-posts
// If published is true, only published posts are returned. Otherwise, only scheduled posts are returned.
func (h *PostHandler) GetPostsByNewsletterId(w http.ResponseWriter, r *http.Request, published bool) {
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

	// Get posts
	posts, err := h.postService.GetPostsByNewsletterId(r.Context(), newsletterID, user.UserID.String(), published)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.responder.RespondJSON(w, http.StatusOK, posts)
}

func (h *PostHandler) GetPostById(w http.ResponseWriter, r *http.Request) {
	newsletterID, err := uuid.Parse(chi.URLParam(r, "newsletterId"))
	if err != nil {
		http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
		return
	}

	postId, err := uuid.Parse(chi.URLParam(r, "postId"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Get user from context
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	post, err := h.postService.GetPostById(r.Context(), newsletterID, postId, user.UserID.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.responder.RespondJSON(w, http.StatusOK, post)
}

func (h *PostHandler) DeletePostById(w http.ResponseWriter, r *http.Request) {
	newsletterID, err := uuid.Parse(chi.URLParam(r, "newsletterId"))
	if err != nil {
		http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
		return
	}

	postId, err := uuid.Parse(chi.URLParam(r, "postId"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Get user from context
	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.postService.DeletePostById(r.Context(), newsletterID, postId, user.UserID.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.responder.RespondJSON(w, http.StatusOK, nil)
}

func (h *PostHandler) PostPost(w http.ResponseWriter, r *http.Request) {
	newsletterID, err := uuid.Parse(chi.URLParam(r, "newsletterId"))
	if err != nil {
		http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
		return
	}

	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		h.responder.HandleError(w, r, models.NewUnauthorizedError("User not authenticated"))
		return
	}

	var req generated.PublishPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.HandleError(w, r, models.NewBadRequestError("Invalid JSON payload"))
		return
	}

	newsletter, err := h.postService.CreatePost(r.Context(), user.UserID, req, newsletterID)
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	h.responder.RespondJSON(w, http.StatusCreated, newsletter)
}

func (h *PostHandler) PutPost(w http.ResponseWriter, r *http.Request) {
	newsletterID, err := uuid.Parse(chi.URLParam(r, "newsletterId"))
	if err != nil {
		http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
		return
	}

	postId, err := uuid.Parse(chi.URLParam(r, "postId"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	user, ok := services.GetUserFromContext(r.Context())
	if !ok {
		h.responder.HandleError(w, r, models.NewUnauthorizedError("User not authenticated"))
		return
	}

	var req generated.PublishPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.HandleError(w, r, models.NewBadRequestError("Invalid JSON payload"))
		return
	}

	post, err := h.postService.UpdatePost(r.Context(), user.UserID, postId, req, newsletterID)
	if err != nil {
		h.responder.HandleError(w, r, err)
		return
	}

	h.responder.RespondJSON(w, http.StatusOK, post)
}
