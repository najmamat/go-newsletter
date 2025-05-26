package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthService handles JWT validation and user authentication
type AuthService struct {
	jwtSecret string
	logger    *slog.Logger
}

// UserClaims represents the claims in our JWT token
type UserClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"sub"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	AAL    string `json:"aal,omitempty"` // Authentication Assurance Level
}

// UserContext represents authenticated user information
type UserContext struct {
	UserID uuid.UUID
	Email  string
	Role   string
	AAL    string
}

// NewAuthService creates a new auth service
func NewAuthService(jwtSecret string, logger *slog.Logger) *AuthService {
	return &AuthService{
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

// ValidateJWT validates a JWT token and returns user claims
func (s *AuthService) ValidateJWT(tokenString string) (*UserClaims, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// GetUserFromToken extracts user context from JWT token
func (s *AuthService) GetUserFromToken(tokenString string) (*UserContext, error) {
	claims, err := s.ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in token: %w", err)
	}

	return &UserContext{
		UserID: userID,
		Email:  claims.Email,
		Role:   claims.Role,
		AAL:    claims.AAL,
	}, nil
}

// IsAdmin checks if the user has admin role
func (uc *UserContext) IsAdmin() bool {
	// Note: Admin status is stored in profiles table and checked at the handler level
	// This method is kept for future use if we decide to include admin status in JWT
	return false
}

// Key for storing user context in request context
type contextKey string

const UserContextKey contextKey = "user"

// AddUserToContext adds user context to the request context
func AddUserToContext(ctx context.Context, user *UserContext) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

// GetUserFromContext extracts user context from request context
func GetUserFromContext(ctx context.Context) (*UserContext, bool) {
	user, ok := ctx.Value(UserContextKey).(*UserContext)
	return user, ok
} 