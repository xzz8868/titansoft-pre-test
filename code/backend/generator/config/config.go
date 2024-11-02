package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	GeneratorServerPort   string
	BackendServerEndpoint string
	RequestPerSecond      int
}

func LoadConfig() (*Config, error) {
	config := &Config{
		GeneratorServerPort:   getEnv("GENERATOR_SERVER_PORT", "8080"),
		BackendServerEndpoint: ensureNoTrailingSlash(getEnv("BACKEND_SERVER_ENDPOINT", "http://localhost")),
		RequestPerSecond:      getRequestPerSecond("REQUESTS_PER_SECOND", 100),
	}

	if config.BackendServerEndpoint == "" {
		return nil, fmt.Errorf("backend server endpoint is not set in environment variables")
	}

	return config, nil
}

func getRequestPerSecond(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Invalid %s value: %v", key, err)
	}
	return value
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func ensureNoTrailingSlash(url string) string {
	return strings.TrimRight(url, "/")
}
