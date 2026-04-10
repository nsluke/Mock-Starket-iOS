package polygon

import (
	"testing"
	"time"
)

func TestGetMarketSession(t *testing.T) {
	ny, _ := time.LoadLocation("America/New_York")

	tests := []struct {
		name   string
		t      time.Time
		expect MarketSession
	}{
		{
			"regular hours",
			time.Date(2026, 4, 8, 12, 0, 0, 0, ny), // Wednesday noon
			SessionRegular,
		},
		{
			"pre-market",
			time.Date(2026, 4, 8, 7, 0, 0, 0, ny), // Wednesday 7am
			SessionPreMarket,
		},
		{
			"after hours",
			time.Date(2026, 4, 8, 17, 0, 0, 0, ny), // Wednesday 5pm
			SessionAfterHours,
		},
		{
			"closed overnight",
			time.Date(2026, 4, 8, 2, 0, 0, 0, ny), // Wednesday 2am
			SessionClosed,
		},
		{
			"saturday",
			time.Date(2026, 4, 11, 12, 0, 0, 0, ny), // Saturday noon
			SessionClosed,
		},
		{
			"sunday",
			time.Date(2026, 4, 12, 12, 0, 0, 0, ny), // Sunday noon
			SessionClosed,
		},
		{
			"market open edge",
			time.Date(2026, 4, 8, 9, 30, 0, 0, ny), // Exactly 9:30
			SessionRegular,
		},
		{
			"market close edge",
			time.Date(2026, 4, 8, 16, 0, 0, 0, ny), // Exactly 4:00pm
			SessionAfterHours,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMarketSession(tt.t)
			if got != tt.expect {
				t.Errorf("GetMarketSession(%v) = %q, want %q", tt.t, got, tt.expect)
			}
		})
	}
}

func TestIsMarketOpen(t *testing.T) {
	ny, _ := time.LoadLocation("America/New_York")

	if !IsMarketOpen(time.Date(2026, 4, 8, 12, 0, 0, 0, ny)) {
		t.Error("expected market open at Wednesday noon")
	}
	if IsMarketOpen(time.Date(2026, 4, 11, 12, 0, 0, 0, ny)) {
		t.Error("expected market closed on Saturday")
	}
}

func TestNextMarketOpen(t *testing.T) {
	ny, _ := time.LoadLocation("America/New_York")

	// Friday evening → should return Monday 9:30
	friday := time.Date(2026, 4, 10, 18, 0, 0, 0, ny)
	next := NextMarketOpen(friday)
	if next.Weekday() != time.Monday {
		t.Errorf("expected Monday, got %s", next.Weekday())
	}
	if next.Hour() != 9 || next.Minute() != 30 {
		t.Errorf("expected 9:30, got %d:%d", next.Hour(), next.Minute())
	}
}
