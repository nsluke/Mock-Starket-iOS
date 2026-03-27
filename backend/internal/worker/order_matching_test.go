package worker

import (
	"testing"

	"github.com/shopspring/decimal"
)

func d(s string) decimal.Decimal {
	v, _ := decimal.NewFromString(s)
	return v
}

func dptr(s string) *decimal.Decimal {
	v := d(s)
	return &v
}

func TestShouldFill_LimitBuy(t *testing.T) {
	w := &OrderMatchingWorker{}

	tests := []struct {
		name     string
		price    string
		limit    string
		expected bool
	}{
		{"price below limit — fill", "95.00", "100.00", true},
		{"price equals limit — fill", "100.00", "100.00", true},
		{"price above limit — no fill", "105.00", "100.00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := w.shouldFill("limit", "buy", d(tt.price), dptr(tt.limit), nil)
			if got != tt.expected {
				t.Errorf("shouldFill(limit, buy, %s, %s) = %v, want %v", tt.price, tt.limit, got, tt.expected)
			}
		})
	}
}

func TestShouldFill_LimitSell(t *testing.T) {
	w := &OrderMatchingWorker{}

	tests := []struct {
		name     string
		price    string
		limit    string
		expected bool
	}{
		{"price above limit — fill", "105.00", "100.00", true},
		{"price equals limit — fill", "100.00", "100.00", true},
		{"price below limit — no fill", "95.00", "100.00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := w.shouldFill("limit", "sell", d(tt.price), dptr(tt.limit), nil)
			if got != tt.expected {
				t.Errorf("shouldFill(limit, sell, %s, %s) = %v, want %v", tt.price, tt.limit, got, tt.expected)
			}
		})
	}
}

func TestShouldFill_StopBuy(t *testing.T) {
	w := &OrderMatchingWorker{}

	tests := []struct {
		name     string
		price    string
		stop     string
		expected bool
	}{
		{"price above stop — fill (breakout)", "105.00", "100.00", true},
		{"price equals stop — fill", "100.00", "100.00", true},
		{"price below stop — no fill", "95.00", "100.00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := w.shouldFill("stop", "buy", d(tt.price), nil, dptr(tt.stop))
			if got != tt.expected {
				t.Errorf("shouldFill(stop, buy, %s, stop=%s) = %v, want %v", tt.price, tt.stop, got, tt.expected)
			}
		})
	}
}

func TestShouldFill_StopSell(t *testing.T) {
	w := &OrderMatchingWorker{}

	tests := []struct {
		name     string
		price    string
		stop     string
		expected bool
	}{
		{"price below stop — fill (stop-loss)", "95.00", "100.00", true},
		{"price equals stop — fill", "100.00", "100.00", true},
		{"price above stop — no fill", "105.00", "100.00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := w.shouldFill("stop", "sell", d(tt.price), nil, dptr(tt.stop))
			if got != tt.expected {
				t.Errorf("shouldFill(stop, sell, %s, stop=%s) = %v, want %v", tt.price, tt.stop, got, tt.expected)
			}
		})
	}
}

func TestShouldFill_StopLimitBuy(t *testing.T) {
	w := &OrderMatchingWorker{}

	tests := []struct {
		name     string
		price    string
		limit    string
		stop     string
		expected bool
	}{
		{"price in range (stop <= price <= limit) — fill", "102.00", "105.00", "100.00", true},
		{"price equals stop — fill", "100.00", "105.00", "100.00", true},
		{"price equals limit — fill", "105.00", "105.00", "100.00", true},
		{"price below stop — no fill", "98.00", "105.00", "100.00", false},
		{"price above limit — no fill", "110.00", "105.00", "100.00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := w.shouldFill("stop_limit", "buy", d(tt.price), dptr(tt.limit), dptr(tt.stop))
			if got != tt.expected {
				t.Errorf("shouldFill(stop_limit, buy, %s, limit=%s, stop=%s) = %v, want %v",
					tt.price, tt.limit, tt.stop, got, tt.expected)
			}
		})
	}
}

func TestShouldFill_StopLimitSell(t *testing.T) {
	w := &OrderMatchingWorker{}

	tests := []struct {
		name     string
		price    string
		limit    string
		stop     string
		expected bool
	}{
		{"price in range (limit <= price <= stop) — fill", "98.00", "95.00", "100.00", true},
		{"price equals stop — fill", "100.00", "95.00", "100.00", true},
		{"price equals limit — fill", "95.00", "95.00", "100.00", true},
		{"price above stop — no fill", "105.00", "95.00", "100.00", false},
		{"price below limit — no fill", "90.00", "95.00", "100.00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := w.shouldFill("stop_limit", "sell", d(tt.price), dptr(tt.limit), dptr(tt.stop))
			if got != tt.expected {
				t.Errorf("shouldFill(stop_limit, sell, %s, limit=%s, stop=%s) = %v, want %v",
					tt.price, tt.limit, tt.stop, got, tt.expected)
			}
		})
	}
}

func TestShouldFill_NilPrices(t *testing.T) {
	w := &OrderMatchingWorker{}

	tests := []struct {
		name      string
		orderType string
		limit     *decimal.Decimal
		stop      *decimal.Decimal
		expected  bool
	}{
		{"limit with nil limit_price", "limit", nil, nil, false},
		{"stop with nil stop_price", "stop", nil, nil, false},
		{"stop_limit with nil limit", "stop_limit", nil, dptr("100"), false},
		{"stop_limit with nil stop", "stop_limit", dptr("100"), nil, false},
		{"stop_limit with both nil", "stop_limit", nil, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := w.shouldFill(tt.orderType, "buy", d("100"), tt.limit, tt.stop)
			if got != tt.expected {
				t.Errorf("shouldFill(%s, buy, 100, ...) = %v, want %v", tt.orderType, got, tt.expected)
			}
		})
	}
}

func TestShouldFill_UnknownOrderType(t *testing.T) {
	w := &OrderMatchingWorker{}
	got := w.shouldFill("market", "buy", d("100"), dptr("100"), dptr("100"))
	if got != false {
		t.Error("unknown order type should return false")
	}
}
