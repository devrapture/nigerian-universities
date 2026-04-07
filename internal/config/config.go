package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	Port        string
	DatabaseURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	appEnv := getEnv("App_Env", "development")

	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		if appEnv != "development" {
			return nil, fmt.Errorf("DATABASE_URL is required in %s environment", appEnv)
		}
		dbURL = "postgres://localhost:5432/todo?sslmode=disable"
	}
	return &Config{
		AppEnv:      appEnv,
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: dbURL,
	}, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
