package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	ServerAddress string
	DatabaseURL   string
	RedisURL      string
}

// LoadConfig loads configuration from environment variables or defaults
func LoadConfig() *Config {
	cfg := &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://user:password@localhost/dbname?sslmode=disable"),
		RedisURL:      getEnv("REDIS_URL", "redis://localhost:6379"),
	}

	if cfg.DatabaseURL == "" {
		panic("DATABASE_URL is required")
	}

	if cfg.RedisURL == "" {
		panic("REDIS_URL is required")
	}

	return cfg
}

// getEnv fetches an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
