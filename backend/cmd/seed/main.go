package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luke/mockstarket/internal/model"
	"github.com/luke/mockstarket/internal/repository"
	"github.com/shopspring/decimal"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://mockstarket:mockstarket_dev@localhost:5432/mockstarket?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	repo := repository.New(pool)
	ctx := context.Background()

	stocks := []model.Stock{
		// ============ STOCKS ============
		// Tech Sector
		{Ticker: "PLNX", Name: "Planetronix", Sector: "Tech", AssetType: "stock", BasePrice: d("142.00"), CurrentPrice: d("142.00"), DayOpen: d("142.00"), DayHigh: d("142.00"), DayLow: d("142.00"), PrevClose: d("140.50"), Volatility: d("0.025"), Drift: d("0.001"), MeanReversion: d("0.10"), Description: ptr("Leading manufacturer of interstellar navigation systems and quantum computing platforms.")},
		{Ticker: "GWNT", Name: "Gwent Industries", Sector: "Tech", AssetType: "stock", BasePrice: d("156.00"), CurrentPrice: d("156.00"), DayOpen: d("156.00"), DayHigh: d("156.00"), DayLow: d("156.00"), PrevClose: d("154.25"), Volatility: d("0.032"), Drift: d("0.002"), MeanReversion: d("0.08"), Description: ptr("Enterprise card-based analytics and strategic decision-making software.")},
		{Ticker: "PIPE", Name: "Pied Piper", Sector: "Tech", AssetType: "stock", BasePrice: d("267.00"), CurrentPrice: d("267.00"), DayOpen: d("267.00"), DayHigh: d("267.00"), DayLow: d("267.00"), PrevClose: d("264.80"), Volatility: d("0.040"), Drift: d("0.003"), MeanReversion: d("0.06"), Description: ptr("Revolutionary middle-out compression technology for decentralized internet infrastructure.")},
		{Ticker: "INIT", Name: "Initech Systems", Sector: "Tech", AssetType: "stock", BasePrice: d("31.00"), CurrentPrice: d("31.00"), DayOpen: d("31.00"), DayHigh: d("31.00"), DayLow: d("31.00"), PrevClose: d("30.75"), Volatility: d("0.035"), Drift: d("-0.001"), MeanReversion: d("0.12"), Description: ptr("Legacy enterprise software solutions. Known for their TPS report management platform.")},
		{Ticker: "LUMN", Name: "Lumon Industries", Sector: "Tech", AssetType: "stock", BasePrice: d("145.00"), CurrentPrice: d("145.00"), DayOpen: d("145.00"), DayHigh: d("145.00"), DayLow: d("145.00"), PrevClose: d("143.50"), Volatility: d("0.025"), Drift: d("0.001"), MeanReversion: d("0.10"), Description: ptr("Biotech and workplace optimization company. Pioneering consciousness-based productivity solutions.")},

		// Consumer Sector
		{Ticker: "DM", Name: "Dunder Mifflin", Sector: "Consumer", AssetType: "stock", BasePrice: d("45.00"), CurrentPrice: d("45.00"), DayOpen: d("45.00"), DayHigh: d("45.00"), DayLow: d("45.00"), PrevClose: d("44.50"), Volatility: d("0.020"), Drift: d("0.000"), MeanReversion: d("0.15"), Description: ptr("Mid-tier paper and office supply company. Strong regional presence in the Northeast.")},
		{Ticker: "MSPC", Name: "Michael Scott Paper", Sector: "Consumer", AssetType: "stock", BasePrice: d("8.50"), CurrentPrice: d("8.50"), DayOpen: d("8.50"), DayHigh: d("8.50"), DayLow: d("8.50"), PrevClose: d("8.40"), Volatility: d("0.055"), Drift: d("-0.002"), MeanReversion: d("0.08"), Description: ptr("Scrappy paper startup disrupting the industry with competitive pricing and personal service.")},
		{Ticker: "SWTE", Name: "Sweet Tea Co", Sector: "Consumer", AssetType: "stock", BasePrice: d("23.00"), CurrentPrice: d("23.00"), DayOpen: d("23.00"), DayHigh: d("23.00"), DayLow: d("23.00"), PrevClose: d("22.80"), Volatility: d("0.022"), Drift: d("0.001"), MeanReversion: d("0.12"), Description: ptr("Premium artisanal beverage company specializing in Southern-style sweet tea blends.")},
		{Ticker: "CHUX", Name: "Chux Headwear", Sector: "Consumer", AssetType: "stock", BasePrice: d("12.00"), CurrentPrice: d("12.00"), DayOpen: d("12.00"), DayHigh: d("12.00"), DayLow: d("12.00"), PrevClose: d("11.85"), Volatility: d("0.050"), Drift: d("0.000"), MeanReversion: d("0.10"), Description: ptr("Trendy headwear and accessories brand popular with Gen Z consumers.")},

		// Defense Sector
		{Ticker: "ZONE", Name: "Danger Zone Defense", Sector: "Defense", AssetType: "stock", BasePrice: d("210.00"), CurrentPrice: d("210.00"), DayOpen: d("210.00"), DayHigh: d("210.00"), DayLow: d("210.00"), PrevClose: d("208.00"), Volatility: d("0.035"), Drift: d("0.002"), MeanReversion: d("0.08"), Description: ptr("Advanced aerospace and defense contractor. Specializes in next-gen fighter jet systems.")},
		{Ticker: "STRK", Name: "Stark Industries", Sector: "Defense", AssetType: "stock", BasePrice: d("412.00"), CurrentPrice: d("412.00"), DayOpen: d("412.00"), DayHigh: d("412.00"), DayLow: d("412.00"), PrevClose: d("408.50"), Volatility: d("0.020"), Drift: d("0.002"), MeanReversion: d("0.10"), Description: ptr("Multinational defense and clean energy conglomerate. Industry leader in arc reactor technology.")},
		{Ticker: "ACME", Name: "Acme Corporation", Sector: "Defense", AssetType: "stock", BasePrice: d("88.00"), CurrentPrice: d("88.00"), DayOpen: d("88.00"), DayHigh: d("88.00"), DayLow: d("88.00"), PrevClose: d("87.25"), Volatility: d("0.028"), Drift: d("0.000"), MeanReversion: d("0.12"), Description: ptr("Diversified industrial manufacturer. Known for creative solutions and rapid prototyping.")},

		// Food Sector
		{Ticker: "KRAB", Name: "Krusty Krab Holdings", Sector: "Food", AssetType: "stock", BasePrice: d("19.00"), CurrentPrice: d("19.00"), DayOpen: d("19.00"), DayHigh: d("19.00"), DayLow: d("19.00"), PrevClose: d("18.80"), Volatility: d("0.042"), Drift: d("0.001"), MeanReversion: d("0.10"), Description: ptr("Fast-casual seafood restaurant chain. Famous for their proprietary secret formula burger.")},
		{Ticker: "BUBS", Name: "Bubs Concessions", Sector: "Food", AssetType: "stock", BasePrice: d("8.50"), CurrentPrice: d("8.50"), DayOpen: d("8.50"), DayHigh: d("8.50"), DayLow: d("8.50"), PrevClose: d("8.40"), Volatility: d("0.055"), Drift: d("-0.001"), MeanReversion: d("0.08"), Description: ptr("Micro-cap concession stand operator. Questionable accounting but loyal customer base.")},
		{Ticker: "BOWL", Name: "Big Kahuna Burger", Sector: "Food", AssetType: "stock", BasePrice: d("34.00"), CurrentPrice: d("34.00"), DayOpen: d("34.00"), DayHigh: d("34.00"), DayLow: d("34.00"), PrevClose: d("33.60"), Volatility: d("0.030"), Drift: d("0.001"), MeanReversion: d("0.10"), Description: ptr("Hawaiian-themed fast food chain. That IS a tasty burger.")},

		// Industrial Sector
		{Ticker: "MOMR", Name: "Mom's Friendly Robotics", Sector: "Industrial", AssetType: "stock", BasePrice: d("340.00"), CurrentPrice: d("340.00"), DayOpen: d("340.00"), DayHigh: d("340.00"), DayLow: d("340.00"), PrevClose: d("337.50"), Volatility: d("0.030"), Drift: d("0.002"), MeanReversion: d("0.08"), Description: ptr("Leading robotics and automation manufacturer. Supplies 70% of all industrial robots worldwide.")},
		{Ticker: "FIGG", Name: "Figgis Financial", Sector: "Industrial", AssetType: "stock", BasePrice: d("67.00"), CurrentPrice: d("67.00"), DayOpen: d("67.00"), DayHigh: d("67.00"), DayLow: d("67.00"), PrevClose: d("66.25"), Volatility: d("0.028"), Drift: d("0.000"), MeanReversion: d("0.12"), Description: ptr("Private investigation firm turned financial services company. Unconventional but effective.")},
		{Ticker: "PDLK", Name: "Paddle King Sports", Sector: "Industrial", AssetType: "stock", BasePrice: d("34.00"), CurrentPrice: d("34.00"), DayOpen: d("34.00"), DayHigh: d("34.00"), DayLow: d("34.00"), PrevClose: d("33.75"), Volatility: d("0.040"), Drift: d("0.001"), MeanReversion: d("0.10"), Description: ptr("Premium sporting goods manufacturer specializing in paddle sports equipment and apparel.")},
		{Ticker: "CBIO", Name: "Sebio Streaming", Sector: "Industrial", AssetType: "stock", BasePrice: d("95.00"), CurrentPrice: d("95.00"), DayOpen: d("95.00"), DayHigh: d("95.00"), DayLow: d("95.00"), PrevClose: d("94.00"), Volatility: d("0.038"), Drift: d("0.001"), MeanReversion: d("0.08"), Description: ptr("Next-generation streaming platform with AI-curated content and neural-direct viewing technology.")},
		{Ticker: "CHSM", Name: "Chu Supply Materials", Sector: "Industrial", AssetType: "stock", BasePrice: d("52.00"), CurrentPrice: d("52.00"), DayOpen: d("52.00"), DayHigh: d("52.00"), DayLow: d("52.00"), PrevClose: d("51.50"), Volatility: d("0.025"), Drift: d("0.000"), MeanReversion: d("0.12"), Description: ptr("Wholesale building materials and supply chain logistics. Reliable dividend payer.")},

		// ============ CRYPTO ============
		{Ticker: "BTC", Name: "Bitcoin", Sector: "Crypto", AssetType: "crypto", BasePrice: d("68420.00"), CurrentPrice: d("68420.00"), DayOpen: d("68420.00"), DayHigh: d("68420.00"), DayLow: d("68420.00"), PrevClose: d("67800.00"), Volatility: d("0.045"), Drift: d("0.002"), MeanReversion: d("0.04"), Description: ptr("The original cryptocurrency. Digital gold and decentralized store of value.")},
		{Ticker: "ETH", Name: "Ethereum", Sector: "Crypto", AssetType: "crypto", BasePrice: d("3450.00"), CurrentPrice: d("3450.00"), DayOpen: d("3450.00"), DayHigh: d("3450.00"), DayLow: d("3450.00"), PrevClose: d("3400.00"), Volatility: d("0.055"), Drift: d("0.003"), MeanReversion: d("0.05"), Description: ptr("Programmable blockchain platform. Powers DeFi, NFTs, and smart contracts.")},
		{Ticker: "SOL", Name: "Solana", Sector: "Crypto", AssetType: "crypto", BasePrice: d("185.00"), CurrentPrice: d("185.00"), DayOpen: d("185.00"), DayHigh: d("185.00"), DayLow: d("185.00"), PrevClose: d("182.00"), Volatility: d("0.070"), Drift: d("0.003"), MeanReversion: d("0.06"), Description: ptr("High-performance blockchain with sub-second finality. Popular for memecoins and DePIN.")},
		{Ticker: "DOGE", Name: "Dogecoin", Sector: "Crypto", AssetType: "crypto", BasePrice: d("0.42"), CurrentPrice: d("0.42"), DayOpen: d("0.42"), DayHigh: d("0.42"), DayLow: d("0.42"), PrevClose: d("0.41"), Volatility: d("0.090"), Drift: d("0.000"), MeanReversion: d("0.08"), Description: ptr("The people's crypto. Much wow. Very currency. Started as a joke, now a movement.")},

		// ============ COMMODITIES ============
		{Ticker: "GOLD", Name: "Gold", Sector: "Commodities", AssetType: "commodity", BasePrice: d("2340.00"), CurrentPrice: d("2340.00"), DayOpen: d("2340.00"), DayHigh: d("2340.00"), DayLow: d("2340.00"), PrevClose: d("2330.00"), Volatility: d("0.012"), Drift: d("0.001"), MeanReversion: d("0.15"), Description: ptr("The timeless safe-haven asset. Historically maintains value during market uncertainty.")},
		{Ticker: "SLVR", Name: "Silver", Sector: "Commodities", AssetType: "commodity", BasePrice: d("29.50"), CurrentPrice: d("29.50"), DayOpen: d("29.50"), DayHigh: d("29.50"), DayLow: d("29.50"), PrevClose: d("29.20"), Volatility: d("0.022"), Drift: d("0.001"), MeanReversion: d("0.12"), Description: ptr("Industrial and precious metal. Used in electronics, solar panels, and jewelry.")},
		{Ticker: "OIL", Name: "Crude Oil", Sector: "Commodities", AssetType: "commodity", BasePrice: d("78.50"), CurrentPrice: d("78.50"), DayOpen: d("78.50"), DayHigh: d("78.50"), DayLow: d("78.50"), PrevClose: d("77.80"), Volatility: d("0.030"), Drift: d("0.000"), MeanReversion: d("0.10"), Description: ptr("Black gold. Global benchmark for energy prices and economic health.")},

		// ============ ETFs ============
		{Ticker: "MKTX", Name: "Mock Total Market ETF", Sector: "ETF", AssetType: "etf", BasePrice: d("100.00"), CurrentPrice: d("100.00"), DayOpen: d("100.00"), DayHigh: d("100.00"), DayLow: d("100.00"), PrevClose: d("99.50"), Volatility: d("0.015"), Drift: d("0.001"), MeanReversion: d("0.10"), Description: ptr("Tracks all stocks in the Mock Starket. Broad market exposure in a single trade.")},
		{Ticker: "TEKX", Name: "Mock Tech ETF", Sector: "ETF", AssetType: "etf", BasePrice: d("150.00"), CurrentPrice: d("150.00"), DayOpen: d("150.00"), DayHigh: d("150.00"), DayLow: d("150.00"), PrevClose: d("148.50"), Volatility: d("0.028"), Drift: d("0.002"), MeanReversion: d("0.08"), Description: ptr("Tracks the top tech stocks. Concentrated exposure to the technology sector.")},
		{Ticker: "DEFX", Name: "Mock Defense ETF", Sector: "ETF", AssetType: "etf", BasePrice: d("200.00"), CurrentPrice: d("200.00"), DayOpen: d("200.00"), DayHigh: d("200.00"), DayLow: d("200.00"), PrevClose: d("198.00"), Volatility: d("0.020"), Drift: d("0.002"), MeanReversion: d("0.10"), Description: ptr("Tracks defense and aerospace stocks. Stable growth with government contract tailwinds.")},
		{Ticker: "CPTX", Name: "Mock Crypto ETF", Sector: "ETF", AssetType: "etf", BasePrice: d("50.00"), CurrentPrice: d("50.00"), DayOpen: d("50.00"), DayHigh: d("50.00"), DayLow: d("50.00"), PrevClose: d("49.00"), Volatility: d("0.050"), Drift: d("0.002"), MeanReversion: d("0.06"), Description: ptr("Tracks major cryptocurrencies. Crypto exposure without managing wallets.")},
	}

	fmt.Println("Seeding stocks...")
	for _, stock := range stocks {
		if err := repo.UpsertStock(ctx, &stock); err != nil {
			log.Printf("failed to seed stock %s: %v", stock.Ticker, err)
		} else {
			fmt.Printf("  ✓ %s - %s ($%.2f)\n", stock.Ticker, stock.Name, stock.BasePrice.InexactFloat64())
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
			log.Printf("failed to seed achievement %s: %v", a.id, err)
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
		// MKTX - Total Market (equal weight across all stocks)
		{"MKTX", "PLNX", "0.05"}, {"MKTX", "GWNT", "0.05"}, {"MKTX", "PIPE", "0.05"},
		{"MKTX", "INIT", "0.05"}, {"MKTX", "LUMN", "0.05"}, {"MKTX", "DM", "0.05"},
		{"MKTX", "MSPC", "0.05"}, {"MKTX", "SWTE", "0.05"}, {"MKTX", "CHUX", "0.05"},
		{"MKTX", "ZONE", "0.05"}, {"MKTX", "STRK", "0.05"}, {"MKTX", "ACME", "0.05"},
		{"MKTX", "KRAB", "0.05"}, {"MKTX", "BUBS", "0.05"}, {"MKTX", "BOWL", "0.05"},
		{"MKTX", "MOMR", "0.05"}, {"MKTX", "FIGG", "0.05"}, {"MKTX", "PDLK", "0.05"},
		{"MKTX", "CBIO", "0.05"}, {"MKTX", "CHSM", "0.05"},

		// TEKX - Tech ETF
		{"TEKX", "PLNX", "0.20"}, {"TEKX", "GWNT", "0.20"}, {"TEKX", "PIPE", "0.25"},
		{"TEKX", "INIT", "0.10"}, {"TEKX", "LUMN", "0.25"},

		// DEFX - Defense ETF
		{"DEFX", "ZONE", "0.30"}, {"DEFX", "STRK", "0.45"}, {"DEFX", "ACME", "0.25"},

		// CPTX - Crypto ETF
		{"CPTX", "BTC", "0.50"}, {"CPTX", "ETH", "0.30"}, {"CPTX", "SOL", "0.15"}, {"CPTX", "DOGE", "0.05"},
	}

	for _, h := range etfHoldings {
		if err := repo.UpsertETFHolding(ctx, h.etf, h.holding, d(h.weight)); err != nil {
			log.Printf("failed to seed ETF holding %s->%s: %v", h.etf, h.holding, err)
		} else {
			fmt.Printf("  ✓ %s holds %s (%.0f%%)\n", h.etf, h.holding, d(h.weight).InexactFloat64()*100)
		}
	}

	fmt.Println("\nSeed complete!")
}

func d(s string) decimal.Decimal {
	v, _ := decimal.NewFromString(s)
	return v
}

func ptr(s string) *string {
	return &s
}
