package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	BcryptCost int
}

func LoadConfig() (*Config, error) {
	// Parse BcryptCost
	bcryptCostStr := getEnv("BCRYPT_COST", "12")
	bcryptCost, err := strconv.Atoi(bcryptCostStr)
	if err != nil {
		return nil, fmt.Errorf("BCRYPT_COST debe ser un número válido: %w", err)
	}

	// Validar que BcryptCost esté en rango válido (10-14)
	if bcryptCost < 10 || bcryptCost > 14 {
		return nil, fmt.Errorf("BCRYPT_COST debe estar entre 10 y 14, se obtuvo: %d", bcryptCost)
	}

	config := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "identidad_user"),
		DBPassword: getEnv("DB_PASSWORD", "identidad_pass_dev"),
		DBName:     getEnv("DB_NAME", "identidad_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		BcryptCost: bcryptCost,
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
