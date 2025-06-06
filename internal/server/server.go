package server

import (
	"log/slog"
	"net/http"

	"go-newsletter/internal/handlers"
	"go-newsletter/internal/services"
	"go-newsletter/internal/utils"
	"go-newsletter/pkg/generated"
)

// Server implements the generated ServerInterface
type Server struct {
	profileHandler    *handlers.ProfileHandler
	authHandler       *handlers.AuthHandler
	authService       *services.AuthService
	mailingService    *services.MailingService
	postService       *services.PostService
	newsletterHandler *handlers.NewsletterHandler
	subscriberHandler *handlers.SubscriberHandler
	postHandler       *handlers.PostHandler
	responder         *utils.HTTPResponder
	logger            *slog.Logger // Keep logger for non-HTTP operations
}

// NewServer creates a new server instance
func NewServer(profileService *services.ProfileService, authService *services.AuthService, logger *slog.Logger, mailingService *services.MailingService, newsletterService *services.NewsletterService, subscriberService *services.SubscriberService, postService *services.PostService, responder *utils.HTTPResponder) *Server {
	return &Server{
		logger:            logger,
		profileHandler:    handlers.NewProfileHandler(profileService, authService, logger),
		authHandler:       handlers.NewAuthHandler(authService, logger),
		authService:       authService,
		mailingService:    mailingService,
		postService:       postService,
		newsletterHandler: handlers.NewNewsletterHandler(newsletterService, profileService, responder),
		subscriberHandler: handlers.NewSubscriberHandler(subscriberService, responder),
		postHandler:       handlers.NewPostHandler(postService, responder),
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
	s.newsletterHandler.GetAllNewsletters(w, r)
}

func (s *Server) DeleteAdminNewslettersNewsletterId(w http.ResponseWriter, r *http.Request) {
	s.newsletterHandler.DeleteNewsletterByID(w, r)
}

func (s *Server) GetAdminUsers(w http.ResponseWriter, r *http.Request) {
	s.profileHandler.GetAllProfiles(w, r)
}

func (s *Server) PutAdminUsersUserIdGrantAdmin(w http.ResponseWriter, r *http.Request) {

	s.profileHandler.GrantAdmin(w, r)
}

func (s *Server) PutAdminUsersUserIdRevokeAdmin(w http.ResponseWriter, r *http.Request) {

	s.profileHandler.RevokeAdmin(w, r)
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

// GetNewsletters handles GET /newsletters - get newsletters owned by current editor
func (s *Server) GetNewsletters(w http.ResponseWriter, r *http.Request) {
	s.newsletterHandler.GetNewslettersOwnedByEditor(w, r)
}

// PostNewsletters handles POST /newsletters - create newsletter
func (s *Server) PostNewsletters(w http.ResponseWriter, r *http.Request) {
	s.newsletterHandler.PostNewsletters(w, r)
}

func (s *Server) DeleteNewslettersNewsletterId(w http.ResponseWriter, r *http.Request) {
	s.newsletterHandler.DeleteNewsletter(w, r)
}

// GetNewslettersNewsletterId handles GET /newsletters/{newsletterId}
func (s *Server) GetNewslettersNewsletterId(w http.ResponseWriter, r *http.Request) {
	s.newsletterHandler.GetNewsletterByID(w, r)
}

// PutNewslettersNewsletterId handles PUT /newsletters/{newsletterId}
func (s *Server) PutNewslettersNewsletterId(w http.ResponseWriter, r *http.Request) {
	s.newsletterHandler.PutNewsletters(w, r)
}

// GetNewslettersNewsletterIdPosts handles GET /newsletters/{newsletterId}/posts and returns only published posts
func (s *Server) GetNewslettersNewsletterIdPosts(w http.ResponseWriter, r *http.Request) {
	s.postHandler.GetPostsByNewsletterId(w, r, true)
}

func (s *Server) PostNewslettersNewsletterIdPosts(w http.ResponseWriter, r *http.Request) {
	s.postHandler.PostPost(w, r)
}

// GetNewslettersNewsletterIdScheduledPosts handles GET /newsletters/{newsletterId}/posts and returns only unpublished posts
func (s *Server) GetNewslettersNewsletterIdScheduledPosts(w http.ResponseWriter, r *http.Request) {
	s.postHandler.GetPostsByNewsletterId(w, r, false)
}

func (s *Server) DeleteNewslettersNewsletterIdScheduledPostsPostId(w http.ResponseWriter, r *http.Request) {
	s.postHandler.DeletePostById(w, r)
}

func (s *Server) GetNewslettersNewsletterIdScheduledPostsPostId(w http.ResponseWriter, r *http.Request) {
	s.postHandler.GetPostById(w, r)
}

func (s *Server) PutNewslettersNewsletterIdScheduledPostsPostId(w http.ResponseWriter, r *http.Request) {
	s.postHandler.PutPost(w, r)
}

func (s *Server) PostNewslettersNewsletterIdSubscribe(w http.ResponseWriter, r *http.Request) {
	s.subscriberHandler.Subscribe(w, r)
}

func (s *Server) GetNewslettersNewsletterIdSubscribers(w http.ResponseWriter, r *http.Request) {
	s.subscriberHandler.ListSubscribers(w, r)
}

func (s *Server) GetSubscribeConfirmConfirmationToken(w http.ResponseWriter, r *http.Request, confirmationToken string) {
	s.subscriberHandler.ConfirmSubscription(w, r, confirmationToken)
}

func (s *Server) GetUnsubscribeUnsubscribeToken(w http.ResponseWriter, r *http.Request, unsubscribeToken string) {
	s.subscriberHandler.Unsubscribe(w, r, unsubscribeToken)
}

func (s *Server) notImplemented(w http.ResponseWriter, r *http.Request) {
	errorResponse := generated.Error{
		Code:    501,
		Message: "Endpoint not yet implemented",
	}
	s.responder.RespondJSON(w, http.StatusNotImplemented, errorResponse)
}
