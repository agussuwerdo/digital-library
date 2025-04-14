package config

import (
	"os"

	"log"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	DatabaseURL string
	JWTSecret   string
}

// LoadConfig loads configuration from environment variables or a .env file
func LoadConfig() *Config {
	// Attempt to load .env file, ignore error if it doesn't exist
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading config from environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return &Config{
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
	}
}
