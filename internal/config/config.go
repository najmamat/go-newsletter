package config

import (
	"fmt"
	"os"
	"time"

	"go-newsletter/internal/utils"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logging  LoggingConfig
	Supabase SupabaseConfig
	Resend   ResendConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	ApiBaseURL   string
	Port         string
	ApiVersion   string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
	MaxConns int32
	MinConns int32
}

type ResendConfig struct {
	Sender string
	ApiKey string
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level string
}

// SupabaseConfig holds Supabase-related configuration
type SupabaseConfig struct {
	URL       string
	AnonKey   string
	JWTSecret string
}

func (c Config) BuildApiBaseUrl() string {
	return fmt.Sprintf("%s:%s/api/v%s", c.Server.ApiBaseURL, c.Server.Port, c.Server.ApiVersion)
}

// Load loads configuration from environment variables
func Load() *Config {

	return &Config{
		Server: ServerConfig{
			ApiBaseURL:   utils.GetEnvWithDefault("API_BASE_URL", "http://localhost"),
			Port:         utils.GetEnvWithDefault("PORT", "8080"),
			ApiVersion:   utils.GetEnvWithDefault("API_VERSION", "1"),
			ReadTimeout:  utils.GetDurationWithDefault("READ_TIMEOUT", 15*time.Second),
			WriteTimeout: utils.GetDurationWithDefault("WRITE_TIMEOUT", 15*time.Second),
		},
		Database: DatabaseConfig{
			Host:     utils.GetEnvWithDefault("PGHOST", "localhost"),
			Port:     utils.GetEnvWithDefault("PGPORT", "5432"),
			User:     utils.GetEnvWithDefault("PGUSER", "postgres.iiivolgfmqsxvlrggwsh"),
			Password: os.Getenv("PGPASSWORD"),
			Database: utils.GetEnvWithDefault("PGDATABASE", "postgres"),
			SSLMode:  utils.GetEnvWithDefault("PGSSLMODE", "require"),
			MaxConns: utils.GetInt32WithDefault("DB_MAX_CONNS", 10),
			MinConns: utils.GetInt32WithDefault("DB_MIN_CONNS", 2),
		},
		Logging: LoggingConfig{
			Level: utils.GetEnvWithDefault("LOG_LEVEL", "info"),
		},
		Supabase: SupabaseConfig{
			URL:       os.Getenv("SUPABASE_URL"),
			AnonKey:   os.Getenv("SUPABASE_ANON_KEY"),
			JWTSecret: os.Getenv("SUPABASE_JWT_SECRET"),
		},
		Resend: ResendConfig{
			Sender: os.Getenv("RESEND_SENDER"),
			ApiKey: os.Getenv("RESEND_API_KEY"),
		},
	}
}
