// internal/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds app configuration
type Config struct {
	ServerPort string
	AppName    string
	AppEnv     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// LoadConfig reads .env and returns Config
func LoadConfig() *Config {
	// Load .env file (if exists)
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found — using system env or defaults")
	}

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		AppName:    getEnv("APP_NAME", "Brutal"),
		AppEnv:     getEnv("APP_ENV", "development"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5433"),     // ← confirm this
		DBUser:     getEnv("DB_USER", "dev"),
		DBPassword: getEnv("DB_PASSWORD", "devpass"),
		DBName:     getEnv("DB_NAME", "brutal"),
	}
}

// Helper to get env var or fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}