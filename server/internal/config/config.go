package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	Env            string
	Addr           string
	DatabaseURL    string
	AllowedOrigins []string
}

func Load() (Config, error) {
	cfg := Config{
		Env:            getEnv("APP_ENV", "development"),
		Addr:           getEnv("ADDR", ":8080"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		AllowedOrigins: splitList(getEnv("ALLOWED_ORIGINS", "http://localhost:4200")),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitList(s string) []string {
	var out []string
	for part := range strings.SplitSeq(s, ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}
