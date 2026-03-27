package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	data := map[string]string{"status": "ok"}
	writeJSON(rr, http.StatusOK, data)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %q", ct)
	}

	var result map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", result["status"])
	}
}

func TestWriteJSON_CustomStatusCode(t *testing.T) {
	rr := httptest.NewRecorder()
	writeJSON(rr, http.StatusCreated, map[string]string{"id": "123"})

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}
}

func TestWriteError(t *testing.T) {
	rr := httptest.NewRecorder()
	writeError(rr, http.StatusBadRequest, "invalid input")

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}

	var result map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if result["error"] != "invalid input" {
		t.Errorf("expected error 'invalid input', got %q", result["error"])
	}
}

func TestGetQueryInt_Default(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	val := getQueryInt(req, "limit", 50)
	if val != 50 {
		t.Errorf("expected default 50, got %d", val)
	}
}

func TestGetQueryInt_Provided(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?limit=25", nil)
	val := getQueryInt(req, "limit", 50)
	if val != 25 {
		t.Errorf("expected 25, got %d", val)
	}
}

func TestGetQueryInt_Invalid(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?limit=abc", nil)
	val := getQueryInt(req, "limit", 50)
	if val != 50 {
		t.Errorf("expected fallback 50 for invalid int, got %d", val)
	}
}

func TestHealthCheck(t *testing.T) {
	// HealthCheck only depends on hub.ClientCount() — we can test with a nil repo
	// by creating a handler with a real hub
	h := &Handler{
		hub: newTestHub(),
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/system/health", nil)

	h.HealthCheck(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode health response: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got %v", result["status"])
	}
	if _, ok := result["time"]; !ok {
		t.Error("expected 'time' field in health response")
	}
	if _, ok := result["ws_clients"]; !ok {
		t.Error("expected 'ws_clients' field in health response")
	}
}
