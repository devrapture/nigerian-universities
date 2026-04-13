package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv               string
	Port                 string
	DatabaseURL          string
	FrontendURL          string
	JwtSecret            string
	JwtExpires           int
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	GITHUB_CLIENT_ID     string
	GITHUB_CLIENT_SECRET string
	GOOGLE_REDIRECT_URL  string
	GITHUB_REDIRECT_URL  string
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

	jwtHours, err := strconv.Atoi(getEnv("JWT_EXPIRES_IN_HOURS", "24"))
	if err != nil {
		log.Println("invalid JWT_EXPIRES_IN_HOURS, defaulting to 24")
		jwtHours = 24
	}

	jwtSecret := getEnv("JWT_SECRET", "")

	if jwtSecret == "" {
		if appEnv == "production" {
			return nil, fmt.Errorf("JWT_SECRET must be set in production")
		}
		log.Println("WARNING: using insecure default JWT_SECRET for development")
		jwtSecret = "dev-secret-do-not-use-in-production"
	}

	googleClientID := getEnv("GOOGLE_CLIENT_ID", "")
	googleClientSecret := getEnv("GOOGLE_CLIENT_SECRET", "")

	if googleClientID == "" || googleClientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be set")
	}

	githubClientID := getEnv("GITHUB_CLIENT_ID", "")
	githubClientSecret := getEnv("GITHUB_CLIENT_SECRET", "")

	if githubClientID == "" || githubClientSecret == "" {
		return nil, fmt.Errorf("GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET must be set")
	}

	return &Config{
		AppEnv:               appEnv,
		Port:                 getEnv("PORT", "8080"),
		DatabaseURL:          dbURL,
		FrontendURL:          getEnv("FRONTEND_URL", "http://localhost:3000"),
		JwtSecret:            jwtSecret,
		JwtExpires:           jwtHours,
		GOOGLE_CLIENT_ID:     googleClientID,
		GOOGLE_CLIENT_SECRET: googleClientSecret,
		GITHUB_CLIENT_ID:     githubClientID,
		GITHUB_CLIENT_SECRET: githubClientSecret,
		GOOGLE_REDIRECT_URL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		GITHUB_REDIRECT_URL:  getEnv("GITHUB_REDIRECT_URL", "http://localhost:8080/auth/github/callback"),
	}, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
