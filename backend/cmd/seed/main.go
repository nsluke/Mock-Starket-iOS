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

	// When using MARKET_DATA_SOURCE=polygon, prices are fetched live from Polygon.io.
	// The prices below are approximate reference values for simulation mode.
	// Volatility/Drift/MeanReversion are used only in simulation mode.

	stocks := []model.Stock{
		// ============ TECHNOLOGY ============
		{Ticker: "AAPL", Name: "Apple Inc.", Sector: "Technology", AssetType: "stock", BasePrice: d("195.00"), CurrentPrice: d("195.00"), DayOpen: d("195.00"), DayHigh: d("195.00"), DayLow: d("195.00"), PrevClose: d("194.00"), Volatility: d("0.0001"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Consumer electronics, software, and services. Maker of iPhone, Mac, and Apple Watch.")},
		{Ticker: "MSFT", Name: "Microsoft Corporation", Sector: "Technology", AssetType: "stock", BasePrice: d("420.00"), CurrentPrice: d("420.00"), DayOpen: d("420.00"), DayHigh: d("420.00"), DayLow: d("420.00"), PrevClose: d("418.00"), Volatility: d("0.00008"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Enterprise software, cloud computing (Azure), and AI. Owns LinkedIn, GitHub, Xbox.")},
		{Ticker: "GOOGL", Name: "Alphabet Inc.", Sector: "Technology", AssetType: "stock", BasePrice: d("175.00"), CurrentPrice: d("175.00"), DayOpen: d("175.00"), DayHigh: d("175.00"), DayLow: d("175.00"), PrevClose: d("174.00"), Volatility: d("0.0001"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Parent company of Google. Search, advertising, cloud, YouTube, and Waymo autonomous vehicles.")},
		{Ticker: "AMZN", Name: "Amazon.com Inc.", Sector: "Technology", AssetType: "stock", BasePrice: d("185.00"), CurrentPrice: d("185.00"), DayOpen: d("185.00"), DayHigh: d("185.00"), DayLow: d("185.00"), PrevClose: d("184.00"), Volatility: d("0.0001"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("E-commerce, cloud computing (AWS), streaming, and AI. The everything store.")},
		{Ticker: "NVDA", Name: "NVIDIA Corporation", Sector: "Technology", AssetType: "stock", BasePrice: d("880.00"), CurrentPrice: d("880.00"), DayOpen: d("880.00"), DayHigh: d("880.00"), DayLow: d("880.00"), PrevClose: d("875.00"), Volatility: d("0.00015"), Drift: d("0.0002"), MeanReversion: d("0.18"), Description: ptr("GPU and AI chip leader. Powers data centers, gaming, autonomous vehicles, and AI training.")},
		{Ticker: "META", Name: "Meta Platforms Inc.", Sector: "Technology", AssetType: "stock", BasePrice: d("500.00"), CurrentPrice: d("500.00"), DayOpen: d("500.00"), DayHigh: d("500.00"), DayLow: d("500.00"), PrevClose: d("497.00"), Volatility: d("0.00012"), Drift: d("0.0001"), MeanReversion: d("0.18"), Description: ptr("Social media (Facebook, Instagram, WhatsApp) and metaverse (Reality Labs). Ad-driven revenue.")},
		{Ticker: "TSLA", Name: "Tesla Inc.", Sector: "Technology", AssetType: "stock", BasePrice: d("245.00"), CurrentPrice: d("245.00"), DayOpen: d("245.00"), DayHigh: d("245.00"), DayLow: d("245.00"), PrevClose: d("243.00"), Volatility: d("0.0002"), Drift: d("0.0001"), MeanReversion: d("0.15"), Description: ptr("Electric vehicles, energy storage, and solar. Also developing autonomous driving and robotics.")},
		{Ticker: "CRM", Name: "Salesforce Inc.", Sector: "Technology", AssetType: "stock", BasePrice: d("270.00"), CurrentPrice: d("270.00"), DayOpen: d("270.00"), DayHigh: d("270.00"), DayLow: d("270.00"), PrevClose: d("268.00"), Volatility: d("0.0001"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Cloud-based CRM platform. Enterprise software for sales, service, marketing, and analytics.")},
		{Ticker: "ORCL", Name: "Oracle Corporation", Sector: "Technology", AssetType: "stock", BasePrice: d("125.00"), CurrentPrice: d("125.00"), DayOpen: d("125.00"), DayHigh: d("125.00"), DayLow: d("125.00"), PrevClose: d("124.00"), Volatility: d("0.0001"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Enterprise database, cloud infrastructure, and applications. Growing cloud business.")},
		{Ticker: "INTC", Name: "Intel Corporation", Sector: "Technology", AssetType: "stock", BasePrice: d("32.00"), CurrentPrice: d("32.00"), DayOpen: d("32.00"), DayHigh: d("32.00"), DayLow: d("32.00"), PrevClose: d("31.50"), Volatility: d("0.00015"), Drift: d("-0.0001"), MeanReversion: d("0.22"), Description: ptr("Semiconductor manufacturer. CPUs, data center chips, and foundry services. Turnaround in progress.")},

		// ============ HEALTHCARE ============
		{Ticker: "JNJ", Name: "Johnson & Johnson", Sector: "Healthcare", AssetType: "stock", BasePrice: d("155.00"), CurrentPrice: d("155.00"), DayOpen: d("155.00"), DayHigh: d("155.00"), DayLow: d("155.00"), PrevClose: d("154.00"), Volatility: d("0.00006"), Drift: d("0.0000"), MeanReversion: d("0.25"), Description: ptr("Pharmaceutical, medical devices, and consumer health products. Defensive dividend stock.")},
		{Ticker: "UNH", Name: "UnitedHealth Group", Sector: "Healthcare", AssetType: "stock", BasePrice: d("520.00"), CurrentPrice: d("520.00"), DayOpen: d("520.00"), DayHigh: d("520.00"), DayLow: d("520.00"), PrevClose: d("517.00"), Volatility: d("0.00008"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Largest US health insurer. Also owns Optum health services and technology division.")},
		{Ticker: "PFE", Name: "Pfizer Inc.", Sector: "Healthcare", AssetType: "stock", BasePrice: d("28.00"), CurrentPrice: d("28.00"), DayOpen: d("28.00"), DayHigh: d("28.00"), DayLow: d("28.00"), PrevClose: d("27.75"), Volatility: d("0.0001"), Drift: d("-0.0001"), MeanReversion: d("0.22"), Description: ptr("Global pharmaceutical company. Vaccines, oncology, and rare disease treatments.")},
		{Ticker: "ABBV", Name: "AbbVie Inc.", Sector: "Healthcare", AssetType: "stock", BasePrice: d("170.00"), CurrentPrice: d("170.00"), DayOpen: d("170.00"), DayHigh: d("170.00"), DayLow: d("170.00"), PrevClose: d("169.00"), Volatility: d("0.00008"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Biopharmaceutical company. Immunology, oncology, and neuroscience therapeutics.")},
		{Ticker: "MRK", Name: "Merck & Co.", Sector: "Healthcare", AssetType: "stock", BasePrice: d("125.00"), CurrentPrice: d("125.00"), DayOpen: d("125.00"), DayHigh: d("125.00"), DayLow: d("125.00"), PrevClose: d("124.00"), Volatility: d("0.00008"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Global pharma leader. Keytruda cancer immunotherapy and animal health division.")},
		{Ticker: "LLY", Name: "Eli Lilly and Company", Sector: "Healthcare", AssetType: "stock", BasePrice: d("780.00"), CurrentPrice: d("780.00"), DayOpen: d("780.00"), DayHigh: d("780.00"), DayLow: d("780.00"), PrevClose: d("775.00"), Volatility: d("0.00012"), Drift: d("0.0002"), MeanReversion: d("0.18"), Description: ptr("Pharmaceutical company. GLP-1 diabetes/obesity drugs driving massive growth.")},

		// ============ FINANCIAL ============
		{Ticker: "JPM", Name: "JPMorgan Chase & Co.", Sector: "Financial", AssetType: "stock", BasePrice: d("195.00"), CurrentPrice: d("195.00"), DayOpen: d("195.00"), DayHigh: d("195.00"), DayLow: d("195.00"), PrevClose: d("194.00"), Volatility: d("0.0001"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Largest US bank. Investment banking, commercial banking, asset management.")},
		{Ticker: "BAC", Name: "Bank of America Corp.", Sector: "Financial", AssetType: "stock", BasePrice: d("35.00"), CurrentPrice: d("35.00"), DayOpen: d("35.00"), DayHigh: d("35.00"), DayLow: d("35.00"), PrevClose: d("34.75"), Volatility: d("0.0001"), Drift: d("0.0000"), MeanReversion: d("0.22"), Description: ptr("Major US bank. Consumer banking, wealth management, and global markets.")},
		{Ticker: "GS", Name: "Goldman Sachs Group", Sector: "Financial", AssetType: "stock", BasePrice: d("385.00"), CurrentPrice: d("385.00"), DayOpen: d("385.00"), DayHigh: d("385.00"), DayLow: d("385.00"), PrevClose: d("382.00"), Volatility: d("0.00012"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Global investment bank. Trading, advisory, asset management, and consumer banking.")},
		{Ticker: "V", Name: "Visa Inc.", Sector: "Financial", AssetType: "stock", BasePrice: d("280.00"), CurrentPrice: d("280.00"), DayOpen: d("280.00"), DayHigh: d("280.00"), DayLow: d("280.00"), PrevClose: d("278.00"), Volatility: d("0.00008"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Global payments technology. Processes billions of transactions annually.")},
		{Ticker: "MA", Name: "Mastercard Inc.", Sector: "Financial", AssetType: "stock", BasePrice: d("460.00"), CurrentPrice: d("460.00"), DayOpen: d("460.00"), DayHigh: d("460.00"), DayLow: d("460.00"), PrevClose: d("457.00"), Volatility: d("0.00008"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Global payment network. Digital payments, cybersecurity, and data analytics.")},

		// ============ ENERGY ============
		{Ticker: "XOM", Name: "Exxon Mobil Corporation", Sector: "Energy", AssetType: "stock", BasePrice: d("105.00"), CurrentPrice: d("105.00"), DayOpen: d("105.00"), DayHigh: d("105.00"), DayLow: d("105.00"), PrevClose: d("104.00"), Volatility: d("0.0001"), Drift: d("0.0000"), MeanReversion: d("0.22"), Description: ptr("Integrated oil and gas. Upstream exploration, refining, and chemical manufacturing.")},
		{Ticker: "CVX", Name: "Chevron Corporation", Sector: "Energy", AssetType: "stock", BasePrice: d("155.00"), CurrentPrice: d("155.00"), DayOpen: d("155.00"), DayHigh: d("155.00"), DayLow: d("155.00"), PrevClose: d("154.00"), Volatility: d("0.0001"), Drift: d("0.0000"), MeanReversion: d("0.22"), Description: ptr("Integrated energy company. Oil, natural gas, and growing renewable energy portfolio.")},
		{Ticker: "COP", Name: "ConocoPhillips", Sector: "Energy", AssetType: "stock", BasePrice: d("115.00"), CurrentPrice: d("115.00"), DayOpen: d("115.00"), DayHigh: d("115.00"), DayLow: d("115.00"), PrevClose: d("114.00"), Volatility: d("0.00012"), Drift: d("0.0000"), MeanReversion: d("0.22"), Description: ptr("Independent E&P company. Focused on low-cost oil and gas production.")},
		{Ticker: "SLB", Name: "Schlumberger Limited", Sector: "Energy", AssetType: "stock", BasePrice: d("48.00"), CurrentPrice: d("48.00"), DayOpen: d("48.00"), DayHigh: d("48.00"), DayLow: d("48.00"), PrevClose: d("47.50"), Volatility: d("0.00012"), Drift: d("0.0000"), MeanReversion: d("0.22"), Description: ptr("Oilfield services. Technology, drilling, and production solutions for oil and gas.")},

		// ============ CONSUMER ============
		{Ticker: "WMT", Name: "Walmart Inc.", Sector: "Consumer", AssetType: "stock", BasePrice: d("165.00"), CurrentPrice: d("165.00"), DayOpen: d("165.00"), DayHigh: d("165.00"), DayLow: d("165.00"), PrevClose: d("164.00"), Volatility: d("0.00006"), Drift: d("0.0001"), MeanReversion: d("0.25"), Description: ptr("World's largest retailer. Grocery, general merchandise, e-commerce, and Sam's Club.")},
		{Ticker: "KO", Name: "The Coca-Cola Company", Sector: "Consumer", AssetType: "stock", BasePrice: d("60.00"), CurrentPrice: d("60.00"), DayOpen: d("60.00"), DayHigh: d("60.00"), DayLow: d("60.00"), PrevClose: d("59.75"), Volatility: d("0.00005"), Drift: d("0.0000"), MeanReversion: d("0.25"), Description: ptr("Global beverage company. Coke, Sprite, Fanta, and 200+ brands in 200+ countries.")},
		{Ticker: "PEP", Name: "PepsiCo Inc.", Sector: "Consumer", AssetType: "stock", BasePrice: d("170.00"), CurrentPrice: d("170.00"), DayOpen: d("170.00"), DayHigh: d("170.00"), DayLow: d("170.00"), PrevClose: d("169.00"), Volatility: d("0.00005"), Drift: d("0.0000"), MeanReversion: d("0.25"), Description: ptr("Beverages and snacks. Pepsi, Lay's, Gatorade, Quaker, and Frito-Lay brands.")},
		{Ticker: "MCD", Name: "McDonald's Corporation", Sector: "Consumer", AssetType: "stock", BasePrice: d("290.00"), CurrentPrice: d("290.00"), DayOpen: d("290.00"), DayHigh: d("290.00"), DayLow: d("290.00"), PrevClose: d("288.00"), Volatility: d("0.00006"), Drift: d("0.0001"), MeanReversion: d("0.22"), Description: ptr("Global fast-food chain. 40,000+ restaurants in 100+ countries. Franchise model.")},
		{Ticker: "NKE", Name: "NIKE Inc.", Sector: "Consumer", AssetType: "stock", BasePrice: d("95.00"), CurrentPrice: d("95.00"), DayOpen: d("95.00"), DayHigh: d("95.00"), DayLow: d("95.00"), PrevClose: d("94.50"), Volatility: d("0.0001"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Athletic footwear, apparel, and equipment. World's largest sportswear company.")},
		{Ticker: "SBUX", Name: "Starbucks Corporation", Sector: "Consumer", AssetType: "stock", BasePrice: d("92.00"), CurrentPrice: d("92.00"), DayOpen: d("92.00"), DayHigh: d("92.00"), DayLow: d("92.00"), PrevClose: d("91.50"), Volatility: d("0.0001"), Drift: d("0.0000"), MeanReversion: d("0.22"), Description: ptr("Global coffeehouse chain. 35,000+ stores. Loyalty program and mobile ordering leader.")},
		{Ticker: "DIS", Name: "The Walt Disney Company", Sector: "Consumer", AssetType: "stock", BasePrice: d("110.00"), CurrentPrice: d("110.00"), DayOpen: d("110.00"), DayHigh: d("110.00"), DayLow: d("110.00"), PrevClose: d("109.00"), Volatility: d("0.0001"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Entertainment conglomerate. Theme parks, streaming (Disney+), Marvel, Star Wars, Pixar.")},

		// ============ INDUSTRIAL ============
		{Ticker: "CAT", Name: "Caterpillar Inc.", Sector: "Industrial", AssetType: "stock", BasePrice: d("340.00"), CurrentPrice: d("340.00"), DayOpen: d("340.00"), DayHigh: d("340.00"), DayLow: d("340.00"), PrevClose: d("338.00"), Volatility: d("0.0001"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Construction and mining equipment. Diesel engines, turbines, locomotives.")},
		{Ticker: "BA", Name: "The Boeing Company", Sector: "Industrial", AssetType: "stock", BasePrice: d("190.00"), CurrentPrice: d("190.00"), DayOpen: d("190.00"), DayHigh: d("190.00"), DayLow: d("190.00"), PrevClose: d("188.00"), Volatility: d("0.00015"), Drift: d("0.0000"), MeanReversion: d("0.18"), Description: ptr("Aerospace and defense. Commercial aircraft, military systems, and space technology.")},
		{Ticker: "HON", Name: "Honeywell International", Sector: "Industrial", AssetType: "stock", BasePrice: d("200.00"), CurrentPrice: d("200.00"), DayOpen: d("200.00"), DayHigh: d("200.00"), DayLow: d("200.00"), PrevClose: d("199.00"), Volatility: d("0.00008"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Diversified industrial conglomerate. Aerospace, building tech, materials, and safety.")},
		{Ticker: "UPS", Name: "United Parcel Service", Sector: "Industrial", AssetType: "stock", BasePrice: d("145.00"), CurrentPrice: d("145.00"), DayOpen: d("145.00"), DayHigh: d("145.00"), DayLow: d("145.00"), PrevClose: d("144.00"), Volatility: d("0.00008"), Drift: d("0.0000"), MeanReversion: d("0.22"), Description: ptr("Global package delivery and supply chain management. E-commerce logistics backbone.")},
		{Ticker: "GE", Name: "GE Aerospace", Sector: "Industrial", AssetType: "stock", BasePrice: d("160.00"), CurrentPrice: d("160.00"), DayOpen: d("160.00"), DayHigh: d("160.00"), DayLow: d("160.00"), PrevClose: d("159.00"), Volatility: d("0.0001"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Aviation engines and services. Spun off from the original General Electric conglomerate.")},

		// ============ ETFs ============
		{Ticker: "SPY", Name: "SPDR S&P 500 ETF Trust", Sector: "ETF", AssetType: "etf", BasePrice: d("510.00"), CurrentPrice: d("510.00"), DayOpen: d("510.00"), DayHigh: d("510.00"), DayLow: d("510.00"), PrevClose: d("508.00"), Volatility: d("0.00006"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Tracks the S&P 500 index. The most traded ETF in the world.")},
		{Ticker: "QQQ", Name: "Invesco QQQ Trust", Sector: "ETF", AssetType: "etf", BasePrice: d("440.00"), CurrentPrice: d("440.00"), DayOpen: d("440.00"), DayHigh: d("440.00"), DayLow: d("440.00"), PrevClose: d("438.00"), Volatility: d("0.00008"), Drift: d("0.0001"), MeanReversion: d("0.18"), Description: ptr("Tracks the Nasdaq-100 index. Heavy tech weighting.")},
		{Ticker: "DIA", Name: "SPDR Dow Jones Industrial", Sector: "ETF", AssetType: "etf", BasePrice: d("390.00"), CurrentPrice: d("390.00"), DayOpen: d("390.00"), DayHigh: d("390.00"), DayLow: d("390.00"), PrevClose: d("388.00"), Volatility: d("0.00006"), Drift: d("0.0001"), MeanReversion: d("0.22"), Description: ptr("Tracks the Dow Jones Industrial Average. 30 blue-chip US stocks.")},
		{Ticker: "IWM", Name: "iShares Russell 2000 ETF", Sector: "ETF", AssetType: "etf", BasePrice: d("200.00"), CurrentPrice: d("200.00"), DayOpen: d("200.00"), DayHigh: d("200.00"), DayLow: d("200.00"), PrevClose: d("199.00"), Volatility: d("0.0001"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Tracks Russell 2000 small-cap index. Broad small-cap US equity exposure.")},
		{Ticker: "VTI", Name: "Vanguard Total Stock Market", Sector: "ETF", AssetType: "etf", BasePrice: d("260.00"), CurrentPrice: d("260.00"), DayOpen: d("260.00"), DayHigh: d("260.00"), DayLow: d("260.00"), PrevClose: d("259.00"), Volatility: d("0.00006"), Drift: d("0.0001"), MeanReversion: d("0.20"), Description: ptr("Tracks entire US stock market. Large, mid, and small-cap exposure.")},

		// ============ CRYPTO ============
		{Ticker: "X:BTCUSD", Name: "Bitcoin", Sector: "Crypto", AssetType: "crypto", BasePrice: d("68000.00"), CurrentPrice: d("68000.00"), DayOpen: d("68000.00"), DayHigh: d("68000.00"), DayLow: d("68000.00"), PrevClose: d("67500.00"), Volatility: d("0.00015"), Drift: d("0.0001"), MeanReversion: d("0.12"), Description: ptr("The original cryptocurrency. Digital gold and decentralized store of value.")},
		{Ticker: "X:ETHUSD", Name: "Ethereum", Sector: "Crypto", AssetType: "crypto", BasePrice: d("3500.00"), CurrentPrice: d("3500.00"), DayOpen: d("3500.00"), DayHigh: d("3500.00"), DayLow: d("3500.00"), PrevClose: d("3450.00"), Volatility: d("0.0002"), Drift: d("0.0001"), MeanReversion: d("0.12"), Description: ptr("Programmable blockchain platform. Powers DeFi, NFTs, and smart contracts.")},
		{Ticker: "X:SOLUSD", Name: "Solana", Sector: "Crypto", AssetType: "crypto", BasePrice: d("180.00"), CurrentPrice: d("180.00"), DayOpen: d("180.00"), DayHigh: d("180.00"), DayLow: d("180.00"), PrevClose: d("178.00"), Volatility: d("0.00025"), Drift: d("0.0001"), MeanReversion: d("0.15"), Description: ptr("High-performance blockchain with sub-second finality.")},
		{Ticker: "X:DOGEUSD", Name: "Dogecoin", Sector: "Crypto", AssetType: "crypto", BasePrice: d("0.15"), CurrentPrice: d("0.15"), DayOpen: d("0.15"), DayHigh: d("0.15"), DayLow: d("0.15"), PrevClose: d("0.14"), Volatility: d("0.0003"), Drift: d("0.0000"), MeanReversion: d("0.18"), Description: ptr("The people's crypto. Much wow. Very currency.")},
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
		// SPY - S&P 500 (top holdings by weight)
		{"SPY", "AAPL", "0.07"}, {"SPY", "MSFT", "0.07"}, {"SPY", "NVDA", "0.06"},
		{"SPY", "AMZN", "0.04"}, {"SPY", "META", "0.03"}, {"SPY", "GOOGL", "0.04"},
		{"SPY", "TSLA", "0.02"}, {"SPY", "JPM", "0.02"}, {"SPY", "V", "0.02"},
		{"SPY", "UNH", "0.02"}, {"SPY", "JNJ", "0.02"}, {"SPY", "XOM", "0.02"},

		// QQQ - Nasdaq 100 (tech-heavy)
		{"QQQ", "AAPL", "0.09"}, {"QQQ", "MSFT", "0.08"}, {"QQQ", "NVDA", "0.07"},
		{"QQQ", "AMZN", "0.06"}, {"QQQ", "META", "0.05"}, {"QQQ", "GOOGL", "0.05"},
		{"QQQ", "TSLA", "0.04"}, {"QQQ", "CRM", "0.03"}, {"QQQ", "INTC", "0.02"},

		// DIA - Dow Jones
		{"DIA", "AAPL", "0.06"}, {"DIA", "MSFT", "0.06"}, {"DIA", "UNH", "0.08"},
		{"DIA", "GS", "0.06"}, {"DIA", "MCD", "0.05"}, {"DIA", "CAT", "0.05"},
		{"DIA", "HON", "0.04"}, {"DIA", "BA", "0.04"}, {"DIA", "V", "0.04"},
		{"DIA", "JPM", "0.04"}, {"DIA", "NKE", "0.03"}, {"DIA", "DIS", "0.03"},

		// IWM - Russell 2000 (equal-ish weight, small caps aren't in our list but use what we have)
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
