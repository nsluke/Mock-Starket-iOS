package polygon

import "testing"

func TestSectorFromSIC(t *testing.T) {
	tests := []struct {
		sic    string
		expect string
	}{
		{"3571", "Technology"},  // Electronic Computers (Apple)
		{"7372", "Technology"},  // Prepackaged Software (Microsoft)
		{"3674", "Technology"},  // Semiconductors (NVIDIA)
		{"5961", "Technology"},  // Catalog/Mail-Order (Amazon)
		{"2834", "Healthcare"},  // Pharmaceutical Preparations (Pfizer)
		{"6324", "Healthcare"},  // Hospital & Medical Service Plans (UNH)
		{"6022", "Financial"},   // State commercial banks (JPM)
		{"6199", "Financial"},   // Credit Institutions
		{"1311", "Energy"},      // Crude Petroleum & Gas (Exxon)
		{"2911", "Energy"},      // Petroleum Refining
		{"5812", "Consumer"},    // Eating Places (MCD, SBUX)
		{"5331", "Consumer"},    // Variety Stores (WMT)
		{"3711", "Industrial"},  // Motor Vehicles (Tesla's SIC)
		{"3721", "Industrial"},  // Aircraft (Boeing)
		{"", "Other"},           // Empty
		{"0000", "Other"},       // Zero
		{"notanumber", "Other"}, // Invalid
	}

	for _, tt := range tests {
		got := SectorFromSIC(tt.sic)
		if got != tt.expect {
			t.Errorf("SectorFromSIC(%q) = %q, want %q", tt.sic, got, tt.expect)
		}
	}
}

func TestSectorFromTickerDetail(t *testing.T) {
	tests := []struct {
		name   string
		detail *TickerDetail
		expect string
	}{
		{
			"nil detail",
			nil,
			"Other",
		},
		{
			"crypto",
			&TickerDetail{Market: "crypto", Ticker: "X:BTCUSD"},
			"Crypto",
		},
		{
			"ETF",
			&TickerDetail{Type: "ETF", Ticker: "SPY"},
			"ETF",
		},
		{
			"stock with SIC",
			&TickerDetail{Type: "CS", Ticker: "AAPL", SICCode: "3571"},
			"Technology",
		},
		{
			"stock without SIC",
			&TickerDetail{Type: "CS", Ticker: "UNKNOWN"},
			"Other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SectorFromTickerDetail(tt.detail)
			if got != tt.expect {
				t.Errorf("SectorFromTickerDetail() = %q, want %q", got, tt.expect)
			}
		})
	}
}
