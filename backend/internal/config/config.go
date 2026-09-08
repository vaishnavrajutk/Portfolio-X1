package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port            string
	AllowedOrigins  []string
	SMTPHost        string
	SMTPPort        string
	SMTPUsername    string
	SMTPPassword    string
	ContactFrom     string
	ContactTo       string
	RateLimitPerMin int
}

func Load() Config {
	return Config{
		Port:            getEnv("PORT", "8080"),
		AllowedOrigins:  splitAndTrim(getEnv("ALLOWED_ORIGINS", "http://localhost:5173")),
		SMTPHost:        os.Getenv("SMTP_HOST"),
		SMTPPort:        getEnv("SMTP_PORT", "587"),
		SMTPUsername:    os.Getenv("SMTP_USERNAME"),
		SMTPPassword:    os.Getenv("SMTP_PASSWORD"),
		ContactFrom:     getEnv("CONTACT_FROM", "portfolio@example.com"),
		ContactTo:       os.Getenv("CONTACT_TO"),
		RateLimitPerMin: getEnvInt("CONTACT_RATE_LIMIT_PER_MIN", 5),
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

func splitAndTrim(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
