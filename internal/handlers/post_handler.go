package handlers

import (
	"errors"
	"go-newsletter/internal/utils"
	"net/http"

	"go-newsletter/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
		switch {
		case errors.Is(err, services.ErrNotFound):
			http.Error(w, "Posts not found", http.StatusNotFound)
		case errors.Is(err, services.ErrForbidden):
			http.Error(w, "You don't have permission to access these posts", http.StatusForbidden)
		default:
			http.Error(w, "Failed to list posts", http.StatusInternalServerError)
		}
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
		switch {
		case errors.Is(err, services.ErrNotFound):
			http.Error(w, "Post not found", http.StatusNotFound)
		case errors.Is(err, services.ErrForbidden):
			http.Error(w, "You don't have permission to access this post", http.StatusForbidden)
		default:
			http.Error(w, "Failed to list post", http.StatusInternalServerError)
		}
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
		switch {
		case errors.Is(err, services.ErrNotFound):
			http.Error(w, "Posts not found", http.StatusNotFound)
		case errors.Is(err, services.ErrForbidden):
			http.Error(w, "You don't have permission to access these posts", http.StatusForbidden)
		default:
			http.Error(w, "Failed to list posts", http.StatusInternalServerError)
		}
		return
	}

	h.responder.RespondJSON(w, http.StatusOK, nil)
}
