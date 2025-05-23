package config

import (
	"os"
	"time"

	"go-newsletter/internal/utils"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logging  LoggingConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
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

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         utils.GetEnvWithDefault("PORT", "8080"),
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
	}
} 