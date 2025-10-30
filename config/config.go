package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Server   string
	Port     int
	Database string
	User     string
	Password string
	Trusted  bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "1433"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Server:   getEnv("DB_SERVER", "MAHMOUD"),
			Port:     dbPort,
			Database: getEnv("DB_DATABASE", "iserver-light"),
			User:     getEnv("DB_USER", "sa"),
			Password: getEnv("DB_PASSWORD", "123"),
			Trusted:  getEnv("DB_TRUSTED", "true") == "true",
		},
	}

	return config, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
