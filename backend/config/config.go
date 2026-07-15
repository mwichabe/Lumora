package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// loadDotEnv reads a .env file (if present) and sets any variables that aren't
// already defined in the real environment. Keeps secrets out of the codebase
// without pulling in a dependency. Real env vars always win.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return // no .env — that's fine
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.Trim(strings.TrimSpace(val), `"'`)
		if key != "" && os.Getenv(key) == "" {
			_ = os.Setenv(key, val)
		}
	}
}

// Config holds runtime configuration sourced from environment variables.
type Config struct {
	Port      string
	JWTSecret string

	// DatabaseURL is a Postgres connection string. When set it wins and the app
	// runs on Postgres (production). When empty the app falls back to a local
	// SQLite file at DBPath, so local dev needs no database server.
	DatabaseURL string
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

	// Paystack (monetization). If PaystackSecret is empty, payments are
	// disabled and the exam stays free — handy for local development.
	PaystackSecret string
	PaystackPublic string
	ExamPriceKES   int // price of the exam + certificate unlock, in whole KES
	KESPerUSD      int // approx KES per 1 USD, used to show a USD equivalent
}

// Load reads configuration from the environment, applying sensible defaults so
// the server runs out of the box for local development.
func Load() Config {
	loadDotEnv(".env")
	return Config{
		Port:        getEnv("PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "lumora-dev-secret-change-me"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
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

		PaystackSecret: getEnv("PAYSTACK_SECRET_KEY", ""),
		PaystackPublic: getEnv("PAYSTACK_PUBLIC_KEY", ""),
		ExamPriceKES:   getEnvInt("EXAM_PRICE_KES", 500),
		KESPerUSD:      getEnvInt("KES_PER_USD", 130),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
