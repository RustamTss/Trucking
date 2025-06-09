package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	MongoURI  string
	JWTSecret string
	DBName    string
}

func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		Port:      getEnv("PORT", "8080"),
		MongoURI:  getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		JWTSecret: getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		DBName:    getEnv("DB_NAME", "business_schedule"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
