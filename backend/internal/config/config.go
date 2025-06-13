package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	ClerkSecretKey string // Chiave segreta di Clerk
	Environment    string
}

func Load() *Config {
	// Carica file .env se esiste
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	config := &Config{
		Port:           getEnv("PORT", "30000"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "@Admin123"),
		DBName:         getEnv("DB_NAME", "thothix-db"),
		ClerkSecretKey: getEnv("CLERK_SECRET_KEY", "development_key"),
		Environment:    getEnv("ENVIRONMENT", "development"),
	}

	if config.ClerkSecretKey == "" || config.ClerkSecretKey == "development_key" {
		log.Println("WARNING: CLERK_SECRET_KEY not set or using development key")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
