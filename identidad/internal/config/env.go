package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "identidad_user"),
		DBPassword: getEnv("DB_PASSWORD", "identidad_pass_dev"),
		DBName:     getEnv("DB_NAME", "identidad_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	if config.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s client_encoding=UTF8",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}
