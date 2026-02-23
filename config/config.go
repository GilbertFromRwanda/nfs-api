package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	Port           string
	AppEnv         string
	AllowedOrigins []string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		DBHost:         getEnv("DB_HOST", "127.0.0.1"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBUser:         getEnv("DB_USER", "root"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "nfs_share_db_v2"),
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production!"),
		Port:           getEnv("PORT", "5002"),
		AppEnv:         getEnv("APP_ENV", "production"),
		AllowedOrigins: parseOrigins(getEnv("ALLOWED_ORIGINS", "*")),
	}
}

func parseOrigins(s string) []string {
	parts := strings.Split(s, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	return origins
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
