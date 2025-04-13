package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

func Load() *Config {
	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/pvs_db?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret"),
		Port:        getEnv("PORT", "8080"),
	}

	return cfg
}

func getEnv(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("⚠️  %s not set, using default: %s", key, defaultVal)
		return defaultVal
	}
	return val
}
