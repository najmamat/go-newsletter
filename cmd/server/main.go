package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"go-newsletter/internal/repository"
	"go-newsletter/internal/server"
	"go-newsletter/internal/services"
	"go-newsletter/internal/utils"
	"go-newsletter/pkg/generated"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	// Initialize dependencies using dependency injection
	profileRepo := repository.NewProfileRepository(dbpool, logger)
	profileService := services.NewProfileService(profileRepo, logger)
	apiServer := server.NewServer(profileService, logger)

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
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(SlogMiddleware(logger))
	r.Use(middleware.Recoverer)

	// Health check route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Mount the generated API routes
	r.Mount("/api/v1", generated.HandlerFromMux(apiServer, chi.NewRouter()))

	return r
}

// SlogMiddleware is a chi middleware for logging requests using slog.
func SlogMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tstart := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				logger.Info("Served request",
					"method", r.Method,
					"path", r.URL.Path,
					"status", ww.Status(),
					"latency_ms", time.Since(tstart).Milliseconds(),
					"bytes_out", ww.BytesWritten(),
					"request_id", middleware.GetReqID(r.Context()),
					"remote_ip", r.RemoteAddr,
					"user_agent", r.UserAgent(),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
} 