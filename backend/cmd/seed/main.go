package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luke/mockstarket/internal/model"
	"github.com/luke/mockstarket/internal/polygon"
	"github.com/luke/mockstarket/internal/repository"
	"github.com/shopspring/decimal"
)

// Tickers to seed — Polygon.io will provide names, sectors, descriptions, and prices.
var tickers = []string{
	// US Stocks
	"AAPL", "MSFT", "GOOGL", "AMZN", "NVDA", "META", "TSLA", "CRM", "ORCL", "INTC",
	"JNJ", "UNH", "PFE", "ABBV", "MRK", "LLY",
	"JPM", "BAC", "GS", "V", "MA",
	"XOM", "CVX", "COP", "SLB",
	"WMT", "KO", "PEP", "MCD", "NKE", "SBUX", "DIS",
	"CAT", "BA", "HON", "UPS", "GE",
	// ETFs
	"SPY", "QQQ", "DIA", "IWM", "VTI",
	// Crypto
	"X:BTCUSD", "X:ETHUSD", "X:SOLUSD", "X:DOGEUSD",
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://mockstarket:mockstarket_dev@localhost:5432/mockstarket?sslmode=disable"
	}

	apiKey := os.Getenv("POLYGON_API_KEY")
	if apiKey == "" {
		fmt.Println("POLYGON_API_KEY not set — seeding with placeholder data only.")
		fmt.Println("Set POLYGON_API_KEY to fetch real names, sectors, and descriptions from Polygon.io.")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	repo := repository.New(pool)
	ctx := context.Background()

	// Clean up old fictional tickers that are not in our real ticker list
	fmt.Println("Cleaning up old tickers...")
	realSet := make(map[string]bool, len(tickers))
	for _, t := range tickers {
		realSet[t] = true
	}

	existingStocks, err := repo.GetAllStocks(ctx)
	if err != nil {
		log.Printf("warning: could not fetch existing stocks: %v", err)
	} else {
		removed := 0
		for _, s := range existingStocks {
			if !realSet[s.Ticker] {
				// Delete dependent records first, then the stock
				_, _ = pool.Exec(ctx, `DELETE FROM price_history WHERE ticker = $1`, s.Ticker)
				_, _ = pool.Exec(ctx, `DELETE FROM etf_holdings WHERE etf_ticker = $1 OR holding_ticker = $1`, s.Ticker)
				_, _ = pool.Exec(ctx, `DELETE FROM option_contracts WHERE ticker = $1`, s.Ticker)
				_, _ = pool.Exec(ctx, `DELETE FROM holdings WHERE ticker = $1`, s.Ticker)
				_, _ = pool.Exec(ctx, `DELETE FROM trades WHERE ticker = $1`, s.Ticker)
				_, _ = pool.Exec(ctx, `DELETE FROM orders WHERE ticker = $1`, s.Ticker)
				_, _ = pool.Exec(ctx, `DELETE FROM price_alerts WHERE ticker = $1`, s.Ticker)
				_, _ = pool.Exec(ctx, `DELETE FROM watchlist WHERE ticker = $1`, s.Ticker)
				_, err := pool.Exec(ctx, `DELETE FROM stocks WHERE ticker = $1`, s.Ticker)
				if err != nil {
					log.Printf("  failed to remove %s: %v", s.Ticker, err)
				} else {
					removed++
				}
			}
		}
		if removed > 0 {
			fmt.Printf("  Removed %d old tickers\n", removed)
		}
	}

	// Build Polygon client if API key is available
	var polygonClient *polygon.Client
	if apiKey != "" {
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))
		polygonClient = polygon.NewClient(apiKey, "https://api.polygon.io", 5, logger)
	}

	fmt.Println("\nSeeding stocks...")
	for _, ticker := range tickers {
		stock := buildStock(ctx, ticker, polygonClient)
		if err := repo.UpsertStock(ctx, &stock); err != nil {
			log.Printf("  failed to seed %s: %v", ticker, err)
		} else {
			fmt.Printf("  ✓ %-12s %-8s %s ($%.2f)\n", stock.Ticker, stock.Sector, stock.Name, stock.BasePrice.InexactFloat64())
		}
	}

	// Seed achievements
	fmt.Println("\nSeeding achievements...")
	achievements := []struct {
		id, name, description, icon, category string
	}{
		{"first_trade", "First Steps", "Execute your first trade", "trophy", "trading"},
		{"ten_trades", "Getting Serious", "Execute 10 trades", "chart.bar", "trading"},
		{"hundred_trades", "Day Trader", "Execute 100 trades", "chart.line.uptrend.xyaxis", "trading"},
		{"first_profit", "In the Green", "Close a profitable trade", "dollarsign.circle", "portfolio"},
		{"double_up", "Double Up", "Reach $200,000 net worth", "arrow.up.forward", "portfolio"},
		{"millionaire", "Millionaire", "Reach $1,000,000 net worth", "star.fill", "portfolio"},
		{"diversified", "Diversified", "Own shares in 10 different stocks", "square.grid.3x3", "portfolio"},
		{"all_in", "All In", "Put 90%+ of your portfolio in one stock", "exclamationmark.triangle", "portfolio"},
		{"collector", "Collector", "Own at least 1 share of every stock", "checklist", "portfolio"},
		{"top_ten", "Leaderboard Climber", "Reach top 10 on the leaderboard", "medal", "social"},
		{"top_three", "Podium Finish", "Reach top 3 on the leaderboard", "trophy.fill", "social"},
		{"number_one", "Champion", "Reach #1 on the leaderboard", "crown", "social"},
		{"streak_3", "Getting Started", "Log in 3 days in a row", "flame", "streak"},
		{"streak_7", "On Fire", "Log in 7 days in a row", "flame.fill", "streak"},
		{"streak_30", "Unstoppable", "Log in 30 days in a row", "bolt.fill", "streak"},
		{"crash_survivor", "Crash Survivor", "Maintain portfolio value during a market crash", "shield.fill", "special"},
		{"buy_the_dip", "Buy the Dip", "Buy a stock at its daily low", "arrow.down.to.line", "skill"},
		{"sell_the_top", "Sell the Top", "Sell a stock at its daily high", "arrow.up.to.line", "skill"},
		{"diamond_hands", "Diamond Hands", "Hold a position for 7+ days", "diamond", "skill"},
		{"paper_hands", "Paper Hands", "Sell within 1 minute of buying", "hand.wave", "skill"},
	}

	for _, a := range achievements {
		_, err := pool.Exec(ctx,
			`INSERT INTO achievements (id, name, description, icon, category) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
			a.id, a.name, a.description, a.icon, a.category)
		if err != nil {
			log.Printf("  failed to seed achievement %s: %v", a.id, err)
		} else {
			fmt.Printf("  ✓ %s - %s\n", a.id, a.name)
		}
	}

	// Seed ETF compositions
	fmt.Println("\nSeeding ETF compositions...")
	etfHoldings := []struct {
		etf, holding string
		weight       string
	}{
		// SPY - S&P 500 (top holdings)
		{"SPY", "AAPL", "0.07"}, {"SPY", "MSFT", "0.07"}, {"SPY", "NVDA", "0.06"},
		{"SPY", "AMZN", "0.04"}, {"SPY", "META", "0.03"}, {"SPY", "GOOGL", "0.04"},
		{"SPY", "TSLA", "0.02"}, {"SPY", "JPM", "0.02"}, {"SPY", "V", "0.02"},
		{"SPY", "UNH", "0.02"}, {"SPY", "JNJ", "0.02"}, {"SPY", "XOM", "0.02"},
		// QQQ - Nasdaq 100
		{"QQQ", "AAPL", "0.09"}, {"QQQ", "MSFT", "0.08"}, {"QQQ", "NVDA", "0.07"},
		{"QQQ", "AMZN", "0.06"}, {"QQQ", "META", "0.05"}, {"QQQ", "GOOGL", "0.05"},
		{"QQQ", "TSLA", "0.04"}, {"QQQ", "CRM", "0.03"}, {"QQQ", "INTC", "0.02"},
		// DIA - Dow Jones
		{"DIA", "AAPL", "0.06"}, {"DIA", "MSFT", "0.06"}, {"DIA", "UNH", "0.08"},
		{"DIA", "GS", "0.06"}, {"DIA", "MCD", "0.05"}, {"DIA", "CAT", "0.05"},
		{"DIA", "HON", "0.04"}, {"DIA", "BA", "0.04"}, {"DIA", "V", "0.04"},
		{"DIA", "JPM", "0.04"}, {"DIA", "NKE", "0.03"}, {"DIA", "DIS", "0.03"},
		// IWM - Russell 2000
		{"IWM", "INTC", "0.10"}, {"IWM", "PFE", "0.10"}, {"IWM", "BAC", "0.10"},
		{"IWM", "SLB", "0.10"}, {"IWM", "NKE", "0.10"}, {"IWM", "SBUX", "0.10"},
		// VTI - Total Market
		{"VTI", "AAPL", "0.06"}, {"VTI", "MSFT", "0.06"}, {"VTI", "NVDA", "0.05"},
		{"VTI", "AMZN", "0.04"}, {"VTI", "META", "0.03"}, {"VTI", "GOOGL", "0.03"},
		{"VTI", "JPM", "0.02"}, {"VTI", "JNJ", "0.02"}, {"VTI", "V", "0.02"},
		{"VTI", "UNH", "0.02"}, {"VTI", "XOM", "0.02"}, {"VTI", "WMT", "0.02"},
	}

	for _, h := range etfHoldings {
		if err := repo.UpsertETFHolding(ctx, h.etf, h.holding, d(h.weight)); err != nil {
			log.Printf("  failed to seed ETF holding %s->%s: %v", h.etf, h.holding, err)
		} else {
			fmt.Printf("  ✓ %s holds %s (%.0f%%)\n", h.etf, h.holding, d(h.weight).InexactFloat64()*100)
		}
	}

	fmt.Println("\nSeed complete!")
}

// buildStock creates a Stock model, fetching real data from Polygon when available.
func buildStock(ctx context.Context, ticker string, client *polygon.Client) model.Stock {
	stock := model.Stock{
		Ticker:        ticker,
		Name:          ticker,
		Sector:        "Other",
		AssetType:     "stock",
		BasePrice:     decimal.Zero,
		CurrentPrice:  decimal.Zero,
		DayOpen:       decimal.Zero,
		DayHigh:       decimal.Zero,
		DayLow:        decimal.Zero,
		PrevClose:     decimal.Zero,
		Volatility:    d("0.0010"),
		Drift:         d("0.0000"),
		MeanReversion: d("0.20"),
	}

	if client == nil {
		// No API key — use fallback data
		fb := fallbackData[ticker]
		if fb.Name != "" {
			stock.Name = fb.Name
		}
		if fb.Sector != "" {
			stock.Sector = fb.Sector
		}
		if fb.AssetType != "" {
			stock.AssetType = fb.AssetType
		}
		return stock
	}

	// Fetch ticker details from Polygon
	detail, err := client.GetTickerDetails(ctx, ticker)
	if err != nil {
		fmt.Printf("    (Polygon details unavailable for %s, using fallback)\n", ticker)
		fb := fallbackData[ticker]
		if fb.Name != "" {
			stock.Name = fb.Name
		}
		if fb.Sector != "" {
			stock.Sector = fb.Sector
		}
		if fb.AssetType != "" {
			stock.AssetType = fb.AssetType
		}
		return stock
	}

	// Populate from Polygon data
	stock.Name = detail.Name
	stock.Sector = polygon.SectorFromTickerDetail(detail)

	if detail.Description != "" {
		stock.Description = ptr(detail.Description)
	}

	if detail.Branding != nil && detail.Branding.IconURL != "" {
		logoURL := detail.Branding.IconURL + "?apiKey=" + os.Getenv("POLYGON_API_KEY")
		stock.LogoURL = &logoURL
	}

	// Determine asset type
	switch {
	case detail.Market == "crypto":
		stock.AssetType = "crypto"
	case detail.Type == "ETF":
		stock.AssetType = "etf"
	default:
		stock.AssetType = "stock"
	}

	// Fetch previous close for initial price
	// (rate limited — sleep briefly between API calls)
	time.Sleep(200 * time.Millisecond)
	bar, err := client.GetPreviousClose(ctx, ticker)
	if err == nil && bar.Close > 0 {
		stock.BasePrice = decimal.NewFromFloat(bar.Close).Round(4)
		stock.CurrentPrice = stock.BasePrice
		stock.DayOpen = decimal.NewFromFloat(bar.Open).Round(4)
		stock.DayHigh = decimal.NewFromFloat(bar.High).Round(4)
		stock.DayLow = decimal.NewFromFloat(bar.Low).Round(4)
		stock.PrevClose = decimal.NewFromFloat(bar.Open).Round(4)
		stock.Volume = int64(bar.Volume)
	}

	return stock
}

// fallbackData provides names/sectors when Polygon API is unavailable.
var fallbackData = map[string]struct {
	Name, Sector, AssetType string
}{
	"AAPL": {"Apple Inc.", "Technology", "stock"}, "MSFT": {"Microsoft Corporation", "Technology", "stock"},
	"GOOGL": {"Alphabet Inc.", "Technology", "stock"}, "AMZN": {"Amazon.com Inc.", "Technology", "stock"},
	"NVDA": {"NVIDIA Corporation", "Technology", "stock"}, "META": {"Meta Platforms Inc.", "Technology", "stock"},
	"TSLA": {"Tesla Inc.", "Technology", "stock"}, "CRM": {"Salesforce Inc.", "Technology", "stock"},
	"ORCL": {"Oracle Corporation", "Technology", "stock"}, "INTC": {"Intel Corporation", "Technology", "stock"},
	"JNJ": {"Johnson & Johnson", "Healthcare", "stock"}, "UNH": {"UnitedHealth Group", "Healthcare", "stock"},
	"PFE": {"Pfizer Inc.", "Healthcare", "stock"}, "ABBV": {"AbbVie Inc.", "Healthcare", "stock"},
	"MRK": {"Merck & Co.", "Healthcare", "stock"}, "LLY": {"Eli Lilly and Company", "Healthcare", "stock"},
	"JPM": {"JPMorgan Chase & Co.", "Financial", "stock"}, "BAC": {"Bank of America Corp.", "Financial", "stock"},
	"GS": {"Goldman Sachs Group", "Financial", "stock"}, "V": {"Visa Inc.", "Financial", "stock"},
	"MA": {"Mastercard Inc.", "Financial", "stock"},
	"XOM": {"Exxon Mobil Corporation", "Energy", "stock"}, "CVX": {"Chevron Corporation", "Energy", "stock"},
	"COP": {"ConocoPhillips", "Energy", "stock"}, "SLB": {"Schlumberger Limited", "Energy", "stock"},
	"WMT": {"Walmart Inc.", "Consumer", "stock"}, "KO": {"The Coca-Cola Company", "Consumer", "stock"},
	"PEP": {"PepsiCo Inc.", "Consumer", "stock"}, "MCD": {"McDonald's Corporation", "Consumer", "stock"},
	"NKE": {"NIKE Inc.", "Consumer", "stock"}, "SBUX": {"Starbucks Corporation", "Consumer", "stock"},
	"DIS": {"The Walt Disney Company", "Consumer", "stock"},
	"CAT": {"Caterpillar Inc.", "Industrial", "stock"}, "BA": {"The Boeing Company", "Industrial", "stock"},
	"HON": {"Honeywell International", "Industrial", "stock"}, "UPS": {"United Parcel Service", "Industrial", "stock"},
	"GE": {"GE Aerospace", "Industrial", "stock"},
	"SPY": {"SPDR S&P 500 ETF Trust", "ETF", "etf"}, "QQQ": {"Invesco QQQ Trust", "ETF", "etf"},
	"DIA": {"SPDR Dow Jones Industrial", "ETF", "etf"}, "IWM": {"iShares Russell 2000 ETF", "ETF", "etf"},
	"VTI": {"Vanguard Total Stock Market", "ETF", "etf"},
	"X:BTCUSD": {"Bitcoin", "Crypto", "crypto"}, "X:ETHUSD": {"Ethereum", "Crypto", "crypto"},
	"X:SOLUSD": {"Solana", "Crypto", "crypto"}, "X:DOGEUSD": {"Dogecoin", "Crypto", "crypto"},
}

func d(s string) decimal.Decimal {
	v, _ := decimal.NewFromString(s)
	return v
}

func ptr(s string) *string {
	return &s
}
