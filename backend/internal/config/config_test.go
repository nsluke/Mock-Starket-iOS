package config

import (
	"os"
	"testing"
)

func TestLoad_RequiresDatabaseURL(t *testing.T) {
	// Ensure DATABASE_URL is unset
	os.Unsetenv("DATABASE_URL")

	_, err := Load()
	if err == nil {
		t.Error("expected error when DATABASE_URL is empty")
	}
}

func TestLoad_Defaults(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost/test")
	defer os.Unsetenv("DATABASE_URL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("default port: got %d, want 8080", cfg.Port)
	}
	if cfg.SimTickMS != 30000 {
		t.Errorf("default sim tick: got %d, want 30000", cfg.SimTickMS)
	}
	if cfg.StartingCash != 100000 {
		t.Errorf("default starting cash: got %f, want 100000", cfg.StartingCash)
	}
	if cfg.MarketEventFreq != 150 {
		t.Errorf("default event freq: got %d, want 150", cfg.MarketEventFreq)
	}
	if cfg.MaxWSClients != 1000 {
		t.Errorf("default max ws clients: got %d, want 1000", cfg.MaxWSClients)
	}
	if cfg.FirebaseProjectID != "mock-starket" {
		t.Errorf("default firebase project: got %q, want 'mock-starket'", cfg.FirebaseProjectID)
	}
	if !cfg.DevMode {
		t.Error("default dev mode should be true")
	}
}

func TestLoad_OverridesFromEnv(t *testing.T) {
	envs := map[string]string{
		"DATABASE_URL":     "postgres://test:test@localhost/test",
		"PORT":             "9090",
		"SIM_TICK_MS":      "500",
		"STARTING_CASH":    "50000",
		"MARKET_EVENT_FREQ": "30",
		"MAX_WS_CLIENTS":   "500",
		"LOG_LEVEL":        "debug",
		"DEV_MODE":         "false",
		"ADMIN_API_KEY":    "secret-key",
	}

	for k, v := range envs {
		os.Setenv(k, v)
		defer os.Unsetenv(k)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 9090 {
		t.Errorf("port: got %d, want 9090", cfg.Port)
	}
	if cfg.SimTickMS != 500 {
		t.Errorf("sim tick: got %d, want 500", cfg.SimTickMS)
	}
	if cfg.StartingCash != 50000 {
		t.Errorf("starting cash: got %f, want 50000", cfg.StartingCash)
	}
	if cfg.MarketEventFreq != 30 {
		t.Errorf("event freq: got %d, want 30", cfg.MarketEventFreq)
	}
	if cfg.MaxWSClients != 500 {
		t.Errorf("max ws clients: got %d, want 500", cfg.MaxWSClients)
	}
	if cfg.DevMode {
		t.Error("dev mode should be false when DEV_MODE=false")
	}
	if cfg.AdminAPIKey != "secret-key" {
		t.Errorf("admin key: got %q, want 'secret-key'", cfg.AdminAPIKey)
	}
}

func TestLoad_CORSOrigins(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost/test")
	os.Setenv("CORS_ORIGINS", "http://localhost:3000,https://app.example.com")
	defer os.Unsetenv("DATABASE_URL")
	defer os.Unsetenv("CORS_ORIGINS")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.CORSOrigins) != 2 {
		t.Fatalf("expected 2 CORS origins, got %d", len(cfg.CORSOrigins))
	}
	if cfg.CORSOrigins[0] != "http://localhost:3000" {
		t.Errorf("CORS origin 0: got %q", cfg.CORSOrigins[0])
	}
	if cfg.CORSOrigins[1] != "https://app.example.com" {
		t.Errorf("CORS origin 1: got %q", cfg.CORSOrigins[1])
	}
}

func TestGetEnvStr_Fallback(t *testing.T) {
	os.Unsetenv("NONEXISTENT_KEY")
	val := getEnvStr("NONEXISTENT_KEY", "default")
	if val != "default" {
		t.Errorf("expected 'default', got %q", val)
	}
}

func TestGetEnvInt_InvalidValue(t *testing.T) {
	os.Setenv("BAD_INT", "not-a-number")
	defer os.Unsetenv("BAD_INT")

	val := getEnvInt("BAD_INT", 42)
	if val != 42 {
		t.Errorf("expected fallback 42 for invalid int, got %d", val)
	}
}

func TestGetEnvFloat_InvalidValue(t *testing.T) {
	os.Setenv("BAD_FLOAT", "not-a-float")
	defer os.Unsetenv("BAD_FLOAT")

	val := getEnvFloat("BAD_FLOAT", 3.14)
	if val != 3.14 {
		t.Errorf("expected fallback 3.14 for invalid float, got %f", val)
	}
}
