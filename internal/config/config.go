package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName      string
	Environment  string
	Port         string
	LogLevel     string
	GeminiAPIKey string
	Model        string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppName:      getEnv("APP_NAME", "Tooling Intelligence"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		Port:         getEnv("PORT", "8080"),
		LogLevel:     getEnv("LOG_LEVEL", "INFO"),
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
		Model:        getEnv("GEMINI_MODEL", "gemini-2.5-flash"),
	}

	if cfg.Port == "" {
		return nil, fmt.Errorf("PORT cannot be empty")
	}

	if cfg.GeminiAPIKey == "" {
		return nil, fmt.Errorf("missing GEMINI_API_KEY")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
