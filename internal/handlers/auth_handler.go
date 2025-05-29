package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"go-newsletter/internal/services"
	"go-newsletter/internal/utils"
	"go-newsletter/pkg/generated"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authService *services.AuthService
	responder   *utils.HTTPResponder
	supabaseURL string
	supabaseKey string
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *services.AuthService, supabaseURL, supabaseKey string, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		responder:   utils.NewHTTPResponder(logger),
		supabaseURL: supabaseURL,
		supabaseKey: supabaseKey,
	}
}

// PostAuthSignup handles POST /auth/signup endpoint
func (h *AuthHandler) PostAuthSignup(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req generated.AuthCredentials
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create Supabase signup request
	supabaseReq := map[string]interface{}{
		"email":    req.Email,
		"password": req.Password,
	}

	// Send request to Supabase
	supabaseResp, err := h.makeSupabaseRequest("/auth/v1/signup", supabaseReq)
	if err != nil {
		h.responder.RespondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Return the Supabase response
	h.responder.RespondJSON(w, http.StatusOK, supabaseResp)
}

// PostAuthSignin handles POST /auth/signin endpoint
func (h *AuthHandler) PostAuthSignin(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req generated.AuthCredentials
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create Supabase signin request
	supabaseReq := map[string]interface{}{
		"email":    req.Email,
		"password": req.Password,
	}

	// Send request to Supabase
	supabaseResp, err := h.makeSupabaseRequest("/auth/v1/token?grant_type=password", supabaseReq)
	if err != nil {
		h.responder.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Return the Supabase response
	h.responder.RespondJSON(w, http.StatusOK, supabaseResp)
}

// PostAuthPasswordResetRequest handles POST /auth/password-reset endpoint
func (h *AuthHandler) PostAuthPasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create Supabase password reset request
	supabaseReq := map[string]interface{}{
		"email": req.Email,
	}

	// Send request to Supabase
	supabaseResp, err := h.makeSupabaseRequest("/auth/v1/recover", supabaseReq)
	if err != nil {
		h.responder.RespondError(w, http.StatusInternalServerError, "Failed to send password reset email")
		return
	}

	// Return the Supabase response
	h.responder.RespondJSON(w, http.StatusOK, supabaseResp)
}

// makeSupabaseRequest is a helper function to make requests to Supabase
func (h *AuthHandler) makeSupabaseRequest(path string, body interface{}) (map[string]interface{}, error) {
	// Create request
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", h.supabaseURL+path, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", h.supabaseKey)
	req.Header.Set("Authorization", "Bearer "+h.supabaseKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result, nil
} 