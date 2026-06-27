package config

import "os"

// Config holds runtime configuration sourced from environment variables.
type Config struct {
	Port        string
	JWTSecret   string
	DBPath      string
	CORSOrigins string
}

// Load reads configuration from the environment, applying sensible defaults so
// the server runs out of the box for local development.
func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "lumora-dev-secret-change-me"),
		DBPath:      getEnv("DB_PATH", "lumora.db"),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
