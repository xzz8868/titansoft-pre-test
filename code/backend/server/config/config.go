package config

import (
    "fmt"
    "os"
)

type Config struct {
    DBUser     string
    DBPassword string
    DBHost     string
    DBPort     string
    DBName     string
    ServerPort string
}

func LoadConfig() (*Config, error) {
    config := &Config{
        DBUser:     getEnv("DB_USER", "test"),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBHost:     getEnv("DB_HOST", "mariadb"),
        DBPort:     getEnv("DB_PORT", "3306"),
        DBName:     getEnv("DB_NAME", "pretest"),
        ServerPort: getEnv("PORT", "8080"),
    }

    return config, nil
}

func (c *Config) GetDSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
