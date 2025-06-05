package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"go-newsletter/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// ListPosts handles GET /newsletters/{newsletterId}/posts
func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
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
	posts, err := h.postService.ListPosts(r.Context(), newsletterID, user.UserID.String())
	if err != nil {
		switch {
		case errors.Is(err, services.ErrNotFound):
			http.Error(w, "Newsletter not found", http.StatusNotFound)
		case errors.Is(err, services.ErrForbidden):
			http.Error(w, "You don't have permission to access this newsletter", http.StatusForbidden)
		default:
			http.Error(w, "Failed to list posts", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"posts": posts,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
