package utils

import (
	"os"
	"strconv"
	"time"
)

// GetEnvWithDefault returns the environment variable value or a default value if not set
func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetDurationWithDefault returns the environment variable as duration or a default value if not set/invalid
func GetDurationWithDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// GetInt32WithDefault returns the environment variable as int32 or a default value if not set/invalid
func GetInt32WithDefault(key string, defaultValue int32) int32 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 32); err == nil {
			return int32(intValue)
		}
	}
	return defaultValue
} 