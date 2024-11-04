package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration values.
type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	ServerPort string
	Salt       string
}

// LoadConfig initializes the configuration with environment variables,
// or defaults if the variables are not set.
func LoadConfig() (*Config, error) {
	config := &Config{
		DBUser:     getEnv("DB_USER", "test"),
		DBPassword: getEnv("DB_PASSWORD", "test"),
		DBHost:     getEnv("DB_HOST", "mariadb"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "pretest"),
		ServerPort: getEnv("PORT", "8080"),
		Salt:       getEnv("SALT", "default_salt_value"),
	}

	return config, nil
}

// GetDSN constructs the Data Source Name (DSN) for database connection.
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// getEnv retrieves the value of an environment variable or returns
// a default value if the variable is not set.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
