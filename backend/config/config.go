package config

import "os"

// Config holds runtime configuration sourced from environment variables.
type Config struct {
	Port        string
	JWTSecret   string
	DBPath      string
	CORSOrigins string

	// Email (welcome message). If SMTPHost is empty, emails are skipped.
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPass     string
	SMTPFrom     string // e.g. "no-reply@lumora.app"
	SMTPFromName string // e.g. "Lumora"
	AppURL       string // link target in the email (the web app)
	LogoURL      string // optional hosted fox logo image; falls back to emoji
}

// Load reads configuration from the environment, applying sensible defaults so
// the server runs out of the box for local development.
func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "lumora-dev-secret-change-me"),
		DBPath:      getEnv("DB_PATH", "lumora.db"),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:3000"),

		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPass:     getEnv("SMTP_PASS", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "no-reply@lumora.app"),
		SMTPFromName: getEnv("SMTP_FROM_NAME", "Lumora"),
		AppURL:       getEnv("APP_URL", "http://localhost:3000"),
		LogoURL:      getEnv("LOGO_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
