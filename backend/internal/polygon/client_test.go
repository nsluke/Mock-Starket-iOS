package polygon

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetPreviousClose(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := PreviousCloseResponse{
			Status:       "OK",
			Ticker:       "AAPL",
			ResultsCount: 1,
			Results: []PrevCloseBar{
				{Open: 150.0, High: 155.0, Low: 148.0, Close: 153.0, Volume: 50000000},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewClient("test-key", srv.URL, 100, nil)
	bar, err := client.GetPreviousClose(context.Background(), "AAPL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bar.Close != 153.0 {
		t.Errorf("expected close 153.0, got %f", bar.Close)
	}
	if bar.Volume != 50000000 {
		t.Errorf("expected volume 50000000, got %f", bar.Volume)
	}
}

func TestClient_GetTickerDetails(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := TickerDetailResponse{
			Status: "OK",
			Results: TickerDetail{
				Ticker:         "AAPL",
				Name:           "Apple Inc.",
				Market:         "stocks",
				Type:           "CS",
				SICCode:        "3571",
				SICDescription: "ELECTRONIC COMPUTERS",
				Description:    "Apple makes iPhones.",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewClient("test-key", srv.URL, 100, nil)
	detail, err := client.GetTickerDetails(context.Background(), "AAPL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if detail.Name != "Apple Inc." {
		t.Errorf("expected name 'Apple Inc.', got %q", detail.Name)
	}
	if detail.SICCode != "3571" {
		t.Errorf("expected SIC 3571, got %q", detail.SICCode)
	}
}

func TestClient_GetPreviousClose_NoResults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := PreviousCloseResponse{Status: "OK", ResultsCount: 0, Results: nil}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewClient("test-key", srv.URL, 100, nil)
	_, err := client.GetPreviousClose(context.Background(), "FAKE")
	if err == nil {
		t.Error("expected error for empty results")
	}
}

func TestClient_RateLimitError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"status":"ERROR","message":"rate limit exceeded"}`))
	}))
	defer srv.Close()

	client := NewClient("test-key", srv.URL, 100, nil)
	_, err := client.GetPreviousClose(context.Background(), "AAPL")
	if err == nil {
		t.Error("expected error for 429 response")
	}
}
