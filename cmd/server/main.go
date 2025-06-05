package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"go-newsletter/internal/config"
	"go-newsletter/internal/middleware"
	"go-newsletter/internal/repository"
	"go-newsletter/internal/server"
	"go-newsletter/internal/services"
	"go-newsletter/internal/utils"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Load environment variables
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		logger.Warn("Error loading .env file", "error", err)
	}

	// Setup server configuration
	port := utils.GetEnvWithDefault("PORT", "8080")

	// Setup database connection
	dbpool, err := initializeDatabase(logger)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	// Load configuration
	cfg := config.Load()

	// Initialize dependencies using dependency injection
	profileRepo := repository.NewProfileRepository(dbpool, logger)
	newsletterRepo := repository.NewNewsletterRepository(dbpool, logger)
	subscriberRepo := repository.NewSubscriberRepository(dbpool, logger)
	newsletterService := services.NewNewsletterService(newsletterRepo, logger)
	profileService := services.NewProfileService(profileRepo, logger)
	authService := services.NewAuthService(cfg.Supabase.JWTSecret, logger)
	mailingService := services.NewMailingService(&cfg.Resend, logger)
	subscriberService := services.NewSubscriberService(subscriberRepo, newsletterRepo, mailingService, cfg, logger)
	postRepo := repository.NewPostRepository(dbpool, logger)
	postService := services.NewPostService(postRepo, newsletterService, subscriberService, logger)
	apiServer := server.NewServer(profileService, authService, logger, mailingService, newsletterService, subscriberService, postService)

	// Initialize router and middleware
	r := setupRouter(logger, apiServer)

	// Start server
	logger.Info("Starting server", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

func initializeDatabase(logger *slog.Logger) (*pgxpool.Pool, error) {
	// Build connection string from individual parameters
	connConfig := map[string]string{
		"user":     "postgres.iiivolgfmqsxvlrggwsh",
		"password": os.Getenv("PGPASSWORD"),
		"host":     os.Getenv("PGHOST"),
		"port":     os.Getenv("PGPORT"),
		"dbname":   "postgres",
		"sslmode":  "require",
	}

	// Convert map to connection string
	var connStr []string
	for k, v := range connConfig {
		if v != "" {
			connStr = append(connStr, fmt.Sprintf("%s=%s", k, v))
		}
	}

	databaseURL := strings.Join(connStr, " ")
	if databaseURL == "" {
		return nil, fmt.Errorf("database connection parameters are not set properly")
	}

	// Parse config and configure connection pool
	parsedConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure connection pool settings
	parsedConfig.MaxConns = 10
	parsedConfig.MinConns = 2
	parsedConfig.MaxConnLifetime = time.Hour
	parsedConfig.MaxConnIdleTime = time.Minute * 30

	// Disable automatic prepared statement caching to avoid conflicts
	parsedConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	// Initialize connection pool
	dbpool, err := pgxpool.NewWithConfig(context.Background(), parsedConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Verify connection
	if err := dbpool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to the database")
	return dbpool, nil
}

func setupRouter(logger *slog.Logger, apiServer *server.Server) chi.Router {
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(SlogMiddleware(logger))
	r.Use(chimiddleware.Recoverer)

	// Health check route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Create API router with auth middleware
	apiRouter := chi.NewRouter()
	authMiddleware := middleware.NewAuthMiddleware(apiServer.GetAuthService(), logger)

	// Public routes (no auth required)
	apiRouter.Group(func(r chi.Router) {
		// Auth
		r.Post("/auth/signup", apiServer.PostAuthSignup)
		r.Post("/auth/signin", apiServer.PostAuthSignin)
		r.Post("/auth/password-reset", apiServer.PostAuthPasswordResetRequest)

		// Newsletter Subscription
		r.Route("/newsletters/{newsletterId}/subscribe", func(r chi.Router) {
			r.Use(middleware.UUIDParamValidationMiddleware("newsletterId"))
			r.Post("/", apiServer.PostNewslettersNewsletterIdSubscribe)
		})
		r.Route("/newsletters/{newsletterId}/unsubscribe", func(r chi.Router) {
			r.Use(middleware.UUIDParamValidationMiddleware("newsletterId"))
			r.Post("/", apiServer.PostNewslettersNewsletterIdUnsubscribe)
		})
		r.Route("/newsletters/{newsletterId}/confirm-subscription", func(r chi.Router) {
			r.Use(middleware.UUIDParamValidationMiddleware("newsletterId"))
			// Assuming token is a query param, not a UUID path param here
			r.Get("/", apiServer.GetNewslettersNewsletterIdConfirmSubscription)
		})
		r.Route("/subscribe/confirm/{confirmationToken}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				token := chi.URLParam(r, "confirmationToken")
				apiServer.GetSubscribeConfirmConfirmationToken(w, r, token)
			})
		})
		r.Route("/unsubscribe/{unsubscribeToken}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				token := chi.URLParam(r, "unsubscribeToken")
				apiServer.GetUnsubscribeUnsubscribeToken(w, r, token)
			})
		})
	})

	// Protected routes (require authentication, any editor)
	apiRouter.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)

		// Profile management
		r.Get("/me", apiServer.GetMe)
		r.Put("/me", apiServer.PutMe)

		// Newsletter management (editor-owned)
		r.Get("/newsletters", apiServer.GetNewsletters)
		r.Post("/newsletters", apiServer.PostNewsletters)

		r.Route("/newsletters/{newsletterId}", func(r chi.Router) {
			r.Use(middleware.UUIDParamValidationMiddleware("newsletterId"))
			r.Get("/", apiServer.GetNewslettersNewsletterId)
			r.Put("/", apiServer.PutNewslettersNewsletterId)
			r.Delete("/", apiServer.DeleteNewslettersNewsletterId)

			// Subscriber management
			r.Get("/subscribers", apiServer.GetNewslettersNewsletterIdSubscribers)

			// Post management (editor-owned)
			r.Route("/posts", func(r chi.Router) {
				r.Get("/", apiServer.GetNewslettersNewsletterIdPosts)
				r.Post("/", apiServer.PostNewslettersNewsletterIdPosts)
			})

			// Scheduled Post management (editor-owned)
			r.Route("/scheduled-posts", func(r chi.Router) {
				r.Get("/", apiServer.GetNewslettersNewsletterIdScheduledPosts)
				r.Route("/{postId}", func(r chi.Router) {
					r.Use(middleware.UUIDParamValidationMiddleware("postId"))
					r.Get("/", apiServer.GetNewslettersNewsletterIdScheduledPostsPostId)
					r.Put("/", apiServer.PutNewslettersNewsletterIdScheduledPostsPostId)
					r.Delete("/", apiServer.DeleteNewslettersNewsletterIdScheduledPostsPostId)
				})
			})
		})
	})

	// Admin routes
	apiRouter.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequireAdmin)
		r.Get("/admin/users", apiServer.GetAdminUsers)
		r.Get("/admin/newsletters", apiServer.GetAdminNewsletters)
		r.With(middleware.UUIDParamValidationMiddleware("newsletterId")).Delete("/admin/newsletters/{newsletterId}", apiServer.DeleteAdminNewslettersNewsletterId)
		r.With(middleware.UUIDParamValidationMiddleware("userId")).Delete("/admin/users/{userId}", apiServer.DeleteAdminUsersUserId)
		r.With(middleware.UUIDParamValidationMiddleware("userId")).Put("/admin/users/{userId}/grant-admin", apiServer.PutAdminUsersUserIdGrantAdmin)
		r.With(middleware.UUIDParamValidationMiddleware("userId")).Put("/admin/users/{userId}/revoke-admin", apiServer.PutAdminUsersUserIdRevokeAdmin)
	})

	// Mount the API router
	r.Mount("/api/v1", apiRouter)

	return r
}

// SlogMiddleware is a chi middleware for logging requests using slog.
func SlogMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tstart := time.Now()
			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				logger.Info("Served request",
					"method", r.Method,
					"path", r.URL.Path,
					"status", ww.Status(),
					"latency_ms", time.Since(tstart).Milliseconds(),
					"bytes_out", ww.BytesWritten(),
					"request_id", chimiddleware.GetReqID(r.Context()),
					"remote_ip", r.RemoteAddr,
					"user_agent", r.UserAgent(),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
