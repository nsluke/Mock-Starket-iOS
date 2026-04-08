package polygon

// SnapshotResponse is the top-level response from GET /v2/snapshot/locale/us/markets/stocks/tickers.
type SnapshotResponse struct {
	Status  string     `json:"status"`
	Count   int        `json:"count"`
	Tickers []Snapshot `json:"tickers"`
}

// SingleSnapshotResponse is the response from GET /v2/snapshot/locale/us/markets/stocks/tickers/{ticker}.
type SingleSnapshotResponse struct {
	Status string   `json:"status"`
	Ticker Snapshot `json:"ticker"`
}

// Snapshot represents a single ticker's snapshot from Polygon.
type Snapshot struct {
	Ticker  string `json:"ticker"`
	Day     Agg    `json:"day"`
	PrevDay Agg    `json:"prevDay"`
	Min     Agg    `json:"min"`
	Updated int64  `json:"updated"`
}

// Agg holds aggregate (OHLCV) data.
type Agg struct {
	Open   float64 `json:"o"`
	High   float64 `json:"h"`
	Low    float64 `json:"l"`
	Close  float64 `json:"c"`
	Volume float64 `json:"v"`
	VWAP   float64 `json:"vw"`
}

// AggregateResponse is the response from GET /v2/aggs/ticker/{ticker}/range/...
type AggregateResponse struct {
	Status       string   `json:"status"`
	Ticker       string   `json:"ticker"`
	ResultsCount int      `json:"resultsCount"`
	Results      []AggBar `json:"results"`
}

// AggBar is a single OHLCV bar from Polygon aggregates.
type AggBar struct {
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	Volume    float64 `json:"v"`
	VWAP      float64 `json:"vw"`
	Timestamp int64   `json:"t"` // Unix milliseconds
	NumTrades int     `json:"n"`
}

// TickerDetailResponse is the response from GET /v3/reference/tickers/{ticker}.
type TickerDetailResponse struct {
	Status  string       `json:"status"`
	Results TickerDetail `json:"results"`
}

// TickerDetail holds reference data for a ticker.
type TickerDetail struct {
	Ticker         string          `json:"ticker"`
	Name           string          `json:"name"`
	Market         string          `json:"market"`          // "stocks", "crypto", "fx"
	Locale         string          `json:"locale"`          // "us", "global"
	Type           string          `json:"type"`            // "CS" (common stock), "ETF"
	Active         bool            `json:"active"`
	SICCode        string          `json:"sic_code"`
	SICDescription string          `json:"sic_description"` // e.g. "ELECTRONIC COMPUTERS"
	Description    string          `json:"description"`
	HomepageURL    string          `json:"homepage_url"`
	MarketCap      float64         `json:"market_cap"`
	Branding       *TickerBranding `json:"branding"`
}

// TickerBranding holds logo/icon URLs from Polygon.
type TickerBranding struct {
	LogoURL string `json:"logo_url"`
	IconURL string `json:"icon_url"`
}

// MarketStatusResponse is the response from GET /v1/marketstatus/now.
type MarketStatusResponse struct {
	Market     string            `json:"market"`     // "open", "closed", "extended-hours"
	EarlyHours bool              `json:"earlyHours"`
	AfterHours bool              `json:"afterHours"`
	Exchanges  map[string]string `json:"exchanges"`
	Currencies map[string]string `json:"currencies"`
}

// PreviousCloseResponse is the response from GET /v2/aggs/ticker/{ticker}/prev.
type PreviousCloseResponse struct {
	Status       string         `json:"status"`
	Ticker       string         `json:"ticker"`
	ResultsCount int            `json:"resultsCount"`
	Results      []PrevCloseBar `json:"results"`
}

// PrevCloseBar is like AggBar but with a flexible timestamp field
// (Polygon free tier may return it as float or int).
type PrevCloseBar struct {
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	Volume    float64 `json:"v"`
	VWAP      float64 `json:"vw"`
	Timestamp any `json:"t"` // Unix milliseconds — type varies by endpoint
	NumTrades int     `json:"n"`
}

// WSMessage represents a message sent/received on the Polygon WebSocket.
type WSMessage struct {
	Action string `json:"action,omitempty"`
	Params string `json:"params,omitempty"`
}

// WSAggregateMessage is a per-minute aggregate received on the WebSocket (AM.* channel).
type WSAggregateMessage struct {
	Event     string  `json:"ev"`  // "AM" (per-minute) or "A" (per-second)
	Ticker    string  `json:"sym"`
	Volume    float64 `json:"v"`
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	VWAP      float64 `json:"vw"`
	StartTime int64   `json:"s"` // Unix milliseconds
	EndTime   int64   `json:"e"` // Unix milliseconds
}

// WSStatusMessage is a control message from the Polygon WebSocket.
type WSStatusMessage struct {
	Event   string `json:"ev"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
