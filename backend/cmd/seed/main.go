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
		// Tech Sector
		{Ticker: "PLNX", Name: "Planetronix", Sector: "Tech", BasePrice: d("142.00"), CurrentPrice: d("142.00"), DayOpen: d("142.00"), DayHigh: d("142.00"), DayLow: d("142.00"), PrevClose: d("140.50"), Volatility: d("0.025"), Drift: d("0.001"), MeanReversion: d("0.10"), Description: ptr("Leading manufacturer of interstellar navigation systems and quantum computing platforms.")},
		{Ticker: "GWNT", Name: "Gwent Industries", Sector: "Tech", BasePrice: d("156.00"), CurrentPrice: d("156.00"), DayOpen: d("156.00"), DayHigh: d("156.00"), DayLow: d("156.00"), PrevClose: d("154.25"), Volatility: d("0.032"), Drift: d("0.002"), MeanReversion: d("0.08"), Description: ptr("Enterprise card-based analytics and strategic decision-making software.")},
		{Ticker: "PIPE", Name: "Pied Piper", Sector: "Tech", BasePrice: d("267.00"), CurrentPrice: d("267.00"), DayOpen: d("267.00"), DayHigh: d("267.00"), DayLow: d("267.00"), PrevClose: d("264.80"), Volatility: d("0.040"), Drift: d("0.003"), MeanReversion: d("0.06"), Description: ptr("Revolutionary middle-out compression technology for decentralized internet infrastructure.")},
		{Ticker: "INIT", Name: "Initech Systems", Sector: "Tech", BasePrice: d("31.00"), CurrentPrice: d("31.00"), DayOpen: d("31.00"), DayHigh: d("31.00"), DayLow: d("31.00"), PrevClose: d("30.75"), Volatility: d("0.035"), Drift: d("-0.001"), MeanReversion: d("0.12"), Description: ptr("Legacy enterprise software solutions. Known for their TPS report management platform.")},
		{Ticker: "LUMN", Name: "Lumon Industries", Sector: "Tech", BasePrice: d("145.00"), CurrentPrice: d("145.00"), DayOpen: d("145.00"), DayHigh: d("145.00"), DayLow: d("145.00"), PrevClose: d("143.50"), Volatility: d("0.025"), Drift: d("0.001"), MeanReversion: d("0.10"), Description: ptr("Biotech and workplace optimization company. Pioneering consciousness-based productivity solutions.")},

		// Consumer Sector
		{Ticker: "DM", Name: "Dunder Mifflin", Sector: "Consumer", BasePrice: d("45.00"), CurrentPrice: d("45.00"), DayOpen: d("45.00"), DayHigh: d("45.00"), DayLow: d("45.00"), PrevClose: d("44.50"), Volatility: d("0.020"), Drift: d("0.000"), MeanReversion: d("0.15"), Description: ptr("Mid-tier paper and office supply company. Strong regional presence in the Northeast.")},
		{Ticker: "MSPC", Name: "Michael Scott Paper", Sector: "Consumer", BasePrice: d("8.50"), CurrentPrice: d("8.50"), DayOpen: d("8.50"), DayHigh: d("8.50"), DayLow: d("8.50"), PrevClose: d("8.40"), Volatility: d("0.055"), Drift: d("-0.002"), MeanReversion: d("0.08"), Description: ptr("Scrappy paper startup disrupting the industry with competitive pricing and personal service.")},
		{Ticker: "SWTE", Name: "Sweet Tea Co", Sector: "Consumer", BasePrice: d("23.00"), CurrentPrice: d("23.00"), DayOpen: d("23.00"), DayHigh: d("23.00"), DayLow: d("23.00"), PrevClose: d("22.80"), Volatility: d("0.022"), Drift: d("0.001"), MeanReversion: d("0.12"), Description: ptr("Premium artisanal beverage company specializing in Southern-style sweet tea blends.")},
		{Ticker: "CHUX", Name: "Chux Headwear", Sector: "Consumer", BasePrice: d("12.00"), CurrentPrice: d("12.00"), DayOpen: d("12.00"), DayHigh: d("12.00"), DayLow: d("12.00"), PrevClose: d("11.85"), Volatility: d("0.050"), Drift: d("0.000"), MeanReversion: d("0.10"), Description: ptr("Trendy headwear and accessories brand popular with Gen Z consumers.")},

		// Defense Sector
		{Ticker: "ZONE", Name: "Danger Zone Defense", Sector: "Defense", BasePrice: d("210.00"), CurrentPrice: d("210.00"), DayOpen: d("210.00"), DayHigh: d("210.00"), DayLow: d("210.00"), PrevClose: d("208.00"), Volatility: d("0.035"), Drift: d("0.002"), MeanReversion: d("0.08"), Description: ptr("Advanced aerospace and defense contractor. Specializes in next-gen fighter jet systems.")},
		{Ticker: "STRK", Name: "Stark Industries", Sector: "Defense", BasePrice: d("412.00"), CurrentPrice: d("412.00"), DayOpen: d("412.00"), DayHigh: d("412.00"), DayLow: d("412.00"), PrevClose: d("408.50"), Volatility: d("0.020"), Drift: d("0.002"), MeanReversion: d("0.10"), Description: ptr("Multinational defense and clean energy conglomerate. Industry leader in arc reactor technology.")},
		{Ticker: "ACME", Name: "Acme Corporation", Sector: "Defense", BasePrice: d("88.00"), CurrentPrice: d("88.00"), DayOpen: d("88.00"), DayHigh: d("88.00"), DayLow: d("88.00"), PrevClose: d("87.25"), Volatility: d("0.028"), Drift: d("0.000"), MeanReversion: d("0.12"), Description: ptr("Diversified industrial manufacturer. Known for creative solutions and rapid prototyping.")},

		// Food Sector
		{Ticker: "KRAB", Name: "Krusty Krab Holdings", Sector: "Food", BasePrice: d("19.00"), CurrentPrice: d("19.00"), DayOpen: d("19.00"), DayHigh: d("19.00"), DayLow: d("19.00"), PrevClose: d("18.80"), Volatility: d("0.042"), Drift: d("0.001"), MeanReversion: d("0.10"), Description: ptr("Fast-casual seafood restaurant chain. Famous for their proprietary secret formula burger.")},
		{Ticker: "BUBS", Name: "Bubs Concessions", Sector: "Food", BasePrice: d("8.50"), CurrentPrice: d("8.50"), DayOpen: d("8.50"), DayHigh: d("8.50"), DayLow: d("8.50"), PrevClose: d("8.40"), Volatility: d("0.055"), Drift: d("-0.001"), MeanReversion: d("0.08"), Description: ptr("Micro-cap concession stand operator. Questionable accounting but loyal customer base.")},
		{Ticker: "BOWL", Name: "Big Kahuna Burger", Sector: "Food", BasePrice: d("34.00"), CurrentPrice: d("34.00"), DayOpen: d("34.00"), DayHigh: d("34.00"), DayLow: d("34.00"), PrevClose: d("33.60"), Volatility: d("0.030"), Drift: d("0.001"), MeanReversion: d("0.10"), Description: ptr("Hawaiian-themed fast food chain. That IS a tasty burger.")},

		// Industrial Sector
		{Ticker: "MOMR", Name: "Mom's Friendly Robotics", Sector: "Industrial", BasePrice: d("340.00"), CurrentPrice: d("340.00"), DayOpen: d("340.00"), DayHigh: d("340.00"), DayLow: d("340.00"), PrevClose: d("337.50"), Volatility: d("0.030"), Drift: d("0.002"), MeanReversion: d("0.08"), Description: ptr("Leading robotics and automation manufacturer. Supplies 70% of all industrial robots worldwide.")},
		{Ticker: "FIGG", Name: "Figgis Financial", Sector: "Industrial", BasePrice: d("67.00"), CurrentPrice: d("67.00"), DayOpen: d("67.00"), DayHigh: d("67.00"), DayLow: d("67.00"), PrevClose: d("66.25"), Volatility: d("0.028"), Drift: d("0.000"), MeanReversion: d("0.12"), Description: ptr("Private investigation firm turned financial services company. Unconventional but effective.")},
		{Ticker: "PDLK", Name: "Paddle King Sports", Sector: "Industrial", BasePrice: d("34.00"), CurrentPrice: d("34.00"), DayOpen: d("34.00"), DayHigh: d("34.00"), DayLow: d("34.00"), PrevClose: d("33.75"), Volatility: d("0.040"), Drift: d("0.001"), MeanReversion: d("0.10"), Description: ptr("Premium sporting goods manufacturer specializing in paddle sports equipment and apparel.")},
		{Ticker: "CBIO", Name: "Sebio Streaming", Sector: "Industrial", BasePrice: d("95.00"), CurrentPrice: d("95.00"), DayOpen: d("95.00"), DayHigh: d("95.00"), DayLow: d("95.00"), PrevClose: d("94.00"), Volatility: d("0.038"), Drift: d("0.001"), MeanReversion: d("0.08"), Description: ptr("Next-generation streaming platform with AI-curated content and neural-direct viewing technology.")},
		{Ticker: "CHSM", Name: "Chu Supply Materials", Sector: "Industrial", BasePrice: d("52.00"), CurrentPrice: d("52.00"), DayOpen: d("52.00"), DayHigh: d("52.00"), DayLow: d("52.00"), PrevClose: d("51.50"), Volatility: d("0.025"), Drift: d("0.000"), MeanReversion: d("0.12"), Description: ptr("Wholesale building materials and supply chain logistics. Reliable dividend payer.")},
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

	fmt.Println("\nSeed complete!")
}

func d(s string) decimal.Decimal {
	v, _ := decimal.NewFromString(s)
	return v
}

func ptr(s string) *string {
	return &s
}
