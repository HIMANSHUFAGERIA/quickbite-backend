package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
	Port        string
	Environment string
	FrontendURL string
}

func Load() *Config {
	env := os.Getenv("ENVIRONMENT")

	if env != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: .env file not found, using environment variables")
		}
	}

	databaseURL := os.Getenv("DATABASE_URL")

	return &Config{
		DatabaseURL: databaseURL,
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "quickbite"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
	}
}
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
