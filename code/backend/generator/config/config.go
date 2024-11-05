package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds application configuration values
type Config struct {
	GeneratorServerPort   string // Port for the generator server to listen on
	BackendServerEndpoint string // Endpoint URL for the backend server
	RequestPerSecond      int    // Number of requests allowed per second
}

// LoadConfig initializes and returns a Config struct, populated with environment variables or defaults
func LoadConfig() (*Config, error) {
	config := &Config{
		GeneratorServerPort:   getEnv("GENERATOR_SERVER_PORT", "8080"),
		BackendServerEndpoint: ensureNoTrailingSlash(getEnv("BACKEND_SERVER_ENDPOINT", "http://localhost")),
	}

	// Ensure BackendServerEndpoint is set
	if config.BackendServerEndpoint == "" {
		return nil, fmt.Errorf("backend server endpoint is not set in environment variables")
	}

	return config, nil
}

// getEnv retrieves the value of an environment variable or returns a default if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// ensureNoTrailingSlash removes any trailing slash from a URL string
func ensureNoTrailingSlash(url string) string {
	return strings.TrimRight(url, "/")
}
