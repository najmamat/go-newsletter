package server

import (
	"log/slog"
	"net/http"

	"go-newsletter/internal/config"
	"go-newsletter/internal/handlers"
	"go-newsletter/internal/services"
	"go-newsletter/internal/utils"
	"go-newsletter/pkg/generated"
)

// Server implements the generated ServerInterface
type Server struct {
	profileHandler *handlers.ProfileHandler
	authHandler    *handlers.AuthHandler
	authService    *services.AuthService
	mailingService *services.MailingService
	responder      *utils.HTTPResponder
	logger         *slog.Logger // Keep logger for non-HTTP operations
}

// NewServer creates a new server instance
func NewServer(profileService *services.ProfileService, authService *services.AuthService, cfg *config.Config, logger *slog.Logger, mailingService *services.MailingService) *Server {
	return &Server{
		profileHandler: handlers.NewProfileHandler(profileService, authService, logger),
		authHandler:    handlers.NewAuthHandler(authService, cfg.Supabase.URL, cfg.Supabase.ServiceKey, logger),
		authService:    authService,
		mailingService: mailingService,
		responder:      utils.NewHTTPResponder(logger),
		logger:         logger,
	}
}

func (s *Server) GetAuthService() *services.AuthService {
	return s.authService
}

func (s *Server) GetMe(w http.ResponseWriter, r *http.Request) {
	s.profileHandler.GetMe(w, r)
}

func (s *Server) PutMe(w http.ResponseWriter, r *http.Request) {
	s.profileHandler.PutMe(w, r)
}

func (s *Server) GetAdminNewsletters(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) DeleteAdminNewslettersNewsletterId(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) GetAdminUsers(w http.ResponseWriter, r *http.Request) {
	s.profileHandler.GetAllProfiles(w, r)
}

func (s *Server) PutAdminUsersUserIdGrantAdmin(w http.ResponseWriter, r *http.Request) {

	s.notImplemented(w, r)
}

func (s *Server) PutAdminUsersUserIdRevokeAdmin(w http.ResponseWriter, r *http.Request) {

	s.notImplemented(w, r)
}

func (s *Server) DeleteAdminUsersUserId(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

// PostAuthSignup handles POST /auth/signup endpoint
func (s *Server) PostAuthSignup(w http.ResponseWriter, r *http.Request) {
	s.authHandler.PostAuthSignup(w, r)
}

// PostAuthSignin handles POST /auth/signin endpoint
func (s *Server) PostAuthSignin(w http.ResponseWriter, r *http.Request) {
	s.authHandler.PostAuthSignin(w, r)
}

// PostAuthPasswordResetRequest handles POST /auth/password-reset endpoint
func (s *Server) PostAuthPasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	s.authHandler.PostAuthPasswordResetRequest(w, r)
}

func (s *Server) GetNewsletters(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) PostNewsletters(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) DeleteNewslettersNewsletterId(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterId(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) PutNewslettersNewsletterId(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterIdPosts(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) PostNewslettersNewsletterIdPosts(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterIdScheduledPosts(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) DeleteNewslettersNewsletterIdScheduledPostsPostId(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterIdScheduledPostsPostId(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) PutNewslettersNewsletterIdScheduledPostsPostId(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) PostNewslettersNewsletterIdSubscribe(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) PostNewslettersNewsletterIdUnsubscribe(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterIdConfirmSubscription(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) GetNewslettersNewsletterIdSubscribers(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r)
}

func (s *Server) GetSubscribeConfirmConfirmationToken(w http.ResponseWriter, r *http.Request, confirmationToken string) {
	s.notImplemented(w, r)
}

func (s *Server) GetUnsubscribeUnsubscribeToken(w http.ResponseWriter, r *http.Request, unsubscribeToken string) {
	s.notImplemented(w, r)
}

func (s *Server) notImplemented(w http.ResponseWriter, r *http.Request) {
	errorResponse := generated.Error{
		Code:    501,
		Message: "Endpoint not yet implemented",
	}
	s.responder.RespondJSON(w, http.StatusNotImplemented, errorResponse)
}
