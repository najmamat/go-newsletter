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
	"github.com/google/uuid"
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
	profileService := services.NewProfileService(profileRepo, logger)
	authService := services.NewAuthService(cfg.Supabase.JWTSecret, logger)
	apiServer := server.NewServer(profileService, authService, logger)

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
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure connection pool settings
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30

	// Disable automatic prepared statement caching to avoid conflicts
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	// Initialize connection pool
	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)
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

	// Protected routes (require authentication, any editor)
	apiRouter.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)
		
		// Profile management
		r.Get("/me", apiServer.GetMe)
		r.Put("/me", apiServer.PutMe)
		
		// Newsletter management (editor-owned)
		r.Get("/newsletters", apiServer.GetNewsletters)
		r.Post("/newsletters", apiServer.PostNewsletters)
		r.Get("/newsletters/{newsletterId}", func(w http.ResponseWriter, r *http.Request) {
			rawID := chi.URLParam(r, "newsletterId")
			if err := validateUUID(rawID); err != nil {
				http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
				return
			}
			apiServer.GetNewslettersNewsletterId(w, r, *utils.StringToUUIDPtr(rawID))
		})
		r.Put("/newsletters/{newsletterId}", func(w http.ResponseWriter, r *http.Request) {
			rawID := chi.URLParam(r, "newsletterId")
			if err := validateUUID(rawID); err != nil {
				http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
				return
			}
			apiServer.PutNewslettersNewsletterId(w, r, *utils.StringToUUIDPtr(rawID))
		})
		r.Delete("/newsletters/{newsletterId}", func(w http.ResponseWriter, r *http.Request) {
			rawID := chi.URLParam(r, "newsletterId")
			if err := validateUUID(rawID); err != nil {
				http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
				return
			}
			apiServer.DeleteNewslettersNewsletterId(w, r, *utils.StringToUUIDPtr(rawID))
		})
		
		// Post management (editor-owned)
		r.Get("/newsletters/{newsletterId}/posts", func(w http.ResponseWriter, r *http.Request) {
			rawID := chi.URLParam(r, "newsletterId")
			if err := validateUUID(rawID); err != nil {
				http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
				return
			}
			apiServer.GetNewslettersNewsletterIdPosts(w, r, *utils.StringToUUIDPtr(rawID))
		})
		r.Post("/newsletters/{newsletterId}/posts", func(w http.ResponseWriter, r *http.Request) {
			rawID := chi.URLParam(r, "newsletterId")
			if err := validateUUID(rawID); err != nil {
				http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
				return
			}
			apiServer.PostNewslettersNewsletterIdPosts(w, r, *utils.StringToUUIDPtr(rawID))
		})
		r.Get("/newsletters/{newsletterId}/scheduled-posts", func(w http.ResponseWriter, r *http.Request) {
			rawID := chi.URLParam(r, "newsletterId")
			if err := validateUUID(rawID); err != nil {
				http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
				return
			}
			apiServer.GetNewslettersNewsletterIdScheduledPosts(w, r, *utils.StringToUUIDPtr(rawID))
		})
		r.Get("/newsletters/{newsletterId}/scheduled-posts/{postId}", func(w http.ResponseWriter, r *http.Request) {
			newsletterID := chi.URLParam(r, "newsletterId")
			postID := chi.URLParam(r, "postId")
			if err := validateUUID(newsletterID); err != nil {
				http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
				return
			}
			if err := validateUUID(postID); err != nil {
				http.Error(w, "Invalid post ID", http.StatusBadRequest)
				return
			}
			apiServer.GetNewslettersNewsletterIdScheduledPostsPostId(w, r, *utils.StringToUUIDPtr(newsletterID), *utils.StringToUUIDPtr(postID))
		})
		r.Put("/newsletters/{newsletterId}/scheduled-posts/{postId}", func(w http.ResponseWriter, r *http.Request) {
			newsletterID := chi.URLParam(r, "newsletterId")
			postID := chi.URLParam(r, "postId")
			if err := validateUUID(newsletterID); err != nil {
				http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
				return
			}
			if err := validateUUID(postID); err != nil {
				http.Error(w, "Invalid post ID", http.StatusBadRequest)
				return
			}
			apiServer.PutNewslettersNewsletterIdScheduledPostsPostId(w, r, *utils.StringToUUIDPtr(newsletterID), *utils.StringToUUIDPtr(postID))
		})
		r.Delete("/newsletters/{newsletterId}/scheduled-posts/{postId}", func(w http.ResponseWriter, r *http.Request) {
			newsletterID := chi.URLParam(r, "newsletterId")
			postID := chi.URLParam(r, "postId")
			if err := validateUUID(newsletterID); err != nil {
				http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
				return
			}
			if err := validateUUID(postID); err != nil {
				http.Error(w, "Invalid post ID", http.StatusBadRequest)
				return
			}
			apiServer.DeleteNewslettersNewsletterIdScheduledPostsPostId(w, r, *utils.StringToUUIDPtr(newsletterID), *utils.StringToUUIDPtr(postID))
		})
	})

	// Admin routes
	apiRouter.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequireAdmin)
		r.Get("/admin/users", apiServer.GetAdminUsers)
		r.Get("/admin/newsletters", apiServer.GetAdminNewsletters)
		// Wrap handlers to match chi router's expected function signature
		r.Delete("/admin/newsletters/{newsletterId}", func(w http.ResponseWriter, r *http.Request) {
			rawID := chi.URLParam(r, "newsletterId")
			if err := validateUUID(rawID); err != nil {
				http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
				return
			}
			apiServer.DeleteAdminNewslettersNewsletterId(w, r, *utils.StringToUUIDPtr(rawID))
		})
		r.Put("/admin/users/{userId}/grant-admin", func(w http.ResponseWriter, r *http.Request) {
			rawID := chi.URLParam(r, "userId")
			if err := validateUUID(rawID); err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}
			apiServer.PutAdminUsersUserIdGrantAdmin(w, r, *utils.StringToUUIDPtr(rawID))
		})
		r.Put("/admin/users/{userId}/revoke-admin", func(w http.ResponseWriter, r *http.Request) {
			rawID := chi.URLParam(r, "userId")
			if err := validateUUID(rawID); err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}
			apiServer.PutAdminUsersUserIdRevokeAdmin(w, r, *utils.StringToUUIDPtr(rawID))
		})
	})

	// Public routes (no authentication required)
	apiRouter.Group(func(r chi.Router) {
		// Authentication documentation endpoints
		r.Post("/auth/signup", apiServer.PostAuthSignup)
		r.Post("/auth/signin", apiServer.PostAuthSignin)
		r.Post("/auth/password-reset-request", apiServer.PostAuthPasswordResetRequest)
		
		// Public subscription endpoints
		r.Post("/newsletters/{newsletterId}/subscribe", func(w http.ResponseWriter, r *http.Request) {
			rawID := chi.URLParam(r, "newsletterId")
			if err := validateUUID(rawID); err != nil {
				http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
				return
			}
			apiServer.PostNewslettersNewsletterIdSubscribe(w, r, *utils.StringToUUIDPtr(rawID))
		})
		r.Get("/subscribe/confirm/{confirmationToken}", func(w http.ResponseWriter, r *http.Request) {
			token := chi.URLParam(r, "confirmationToken")
			apiServer.GetSubscribeConfirmConfirmationToken(w, r, token)
		})
		r.Get("/unsubscribe/{unsubscribeToken}", func(w http.ResponseWriter, r *http.Request) {
			token := chi.URLParam(r, "unsubscribeToken")
			apiServer.GetUnsubscribeUnsubscribeToken(w, r, token)
		})
	})

	// Mount the API router
	r.Mount("/api/v1", apiRouter)

	return r
}

// validateUUID checks if a string is a valid UUID
func validateUUID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("invalid UUID format: %w", err)
	}
	return nil
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