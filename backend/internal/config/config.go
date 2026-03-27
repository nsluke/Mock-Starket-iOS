package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                    int
	DatabaseURL             string
	FirebaseProjectID       string
	FirebaseCredentialsFile string
	DevMode                 bool
	CORSOrigins             []string
	LogLevel                string
	SimTickMS               int
	StartingCash            float64
	MarketEventFreq         int
	AdminAPIKey             string
	MaxWSClients            int
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:                    getEnvInt("PORT", 8080),
		DatabaseURL:             getEnvStr("DATABASE_URL", ""),
		FirebaseProjectID:       getEnvStr("FIREBASE_PROJECT_ID", "mock-starket"),
		FirebaseCredentialsFile: getEnvStr("FIREBASE_CREDENTIALS_FILE", ""),
		DevMode:                 getEnvStr("DEV_MODE", "true") == "true",
		CORSOrigins:             strings.Split(getEnvStr("CORS_ORIGINS", "http://localhost:3000"), ","),
		LogLevel:                getEnvStr("LOG_LEVEL", "info"),
		SimTickMS:               getEnvInt("SIM_TICK_MS", 2000),
		StartingCash:            getEnvFloat("STARTING_CASH", 100000),
		MarketEventFreq:         getEnvInt("MARKET_EVENT_FREQ", 60),
		AdminAPIKey:             getEnvStr("ADMIN_API_KEY", ""),
		MaxWSClients:            getEnvInt("MAX_WS_CLIENTS", 1000),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnvStr(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if val := os.Getenv(key); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return fallback
}
