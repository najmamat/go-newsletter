package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Profile struct to match the database table
type Profile struct {
	ID        string    `json:"id"`
	FullName  *string   `json:"full_name,omitempty"` // Use pointers for nullable fields
	AvatarURL *string   `json:"avatar_url,omitempty"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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
	port := getEnvWithDefault("PORT", "8080")

	// Setup database connection
	dbpool, err := initializeDatabase(logger)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	// Initialize router and middleware
	r := setupRouter(logger, dbpool)

	// Start server
	logger.Info("Starting server", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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

	// Initialize connection pool
	dbpool, err := pgxpool.New(context.Background(), databaseURL)
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

func setupRouter(logger *slog.Logger, dbpool *pgxpool.Pool) chi.Router {
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

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/profiles", handleGetProfiles(logger, dbpool))
	})

	return r
}

//TODO: only temporary for testing purposes, later delete this
func handleGetProfiles(logger *slog.Logger, dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := dbpool.Query(context.Background(), "SELECT id, full_name, avatar_url, is_admin, created_at, updated_at FROM public.profiles")
		if err != nil {
			logger.ErrorContext(r.Context(), "Failed to query profiles", "error", err)
			http.Error(w, "Failed to retrieve profiles", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		profiles := []Profile{}
		for rows.Next() {
			var p Profile
			if err := rows.Scan(&p.ID, &p.FullName, &p.AvatarURL, &p.IsAdmin, &p.CreatedAt, &p.UpdatedAt); err != nil {
				logger.ErrorContext(r.Context(), "Failed to scan profile row", "error", err)
				http.Error(w, "Failed to process profiles", http.StatusInternalServerError)
				return
			}
			profiles = append(profiles, p)
		}

		if rows.Err() != nil {
			logger.ErrorContext(r.Context(), "Error iterating profile rows", "error", rows.Err())
			http.Error(w, "Failed to retrieve profiles data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(profiles); err != nil {
			logger.ErrorContext(r.Context(), "Failed to encode profiles to JSON", "error", err)
			http.Error(w, "Failed to prepare response", http.StatusInternalServerError)
		}
	}
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