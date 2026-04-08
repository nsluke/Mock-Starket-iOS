# Mock Starket

**Learn the stock market with real data and zero risk.**

Paper trade 46 real assets -- Apple, NVIDIA, Tesla, Bitcoin, SPY, and more -- with $100,000 in virtual cash. Prices come from live market data via [Polygon.io](https://polygon.io), so you're practicing with the same numbers Wall Street sees. Compete on the leaderboard, earn achievements, and build confidence before you invest real money.

Built as a full-stack, multi-platform application spanning Go, Next.js, SwiftUI, and Jetpack Compose.

---

## What Makes This Different

Most paper trading apps use fake data or delayed feeds. Mock Starket pulls **real-time prices from NYSE, NASDAQ, and crypto exchanges** via the Polygon.io API. When you buy AAPL at $253, that's Apple's actual closing price. When Bitcoin moves, your portfolio moves with it.

The app also includes a **simulation mode** with a custom stochastic engine (Geometric Brownian Motion with mean reversion and sector correlation) for development and offline use.

---

## What I Built

### Backend -- Go

A production-grade REST API and WebSocket server with a pluggable market data architecture that supports both live Polygon.io data and simulated prices.

**Highlights:**
- **Dual market data sources** via a `PriceProvider` interface -- switch between live Polygon.io data and simulation with a single env var
- **Polygon.io integration** -- REST client with rate limiting and caching, WebSocket client for real-time streaming (paid tier), per-ticker polling fallback (free tier), SIC-to-sector mapping, and market hours scheduling
- **Order matching engine** evaluating limit, stop, and stop-limit orders against live prices every tick
- **Real-time WebSocket** broadcasting price updates, trade confirmations, and alert notifications to connected clients
- Background workers for price history recording (OHLCV at 4 intervals), leaderboard computation, achievement evaluation, daily challenge generation, options pricing (Black-Scholes), and price alert monitoring
- **42 unit tests** covering order matching logic, OHLCV aggregation, middleware, config loading, and HTTP handlers

**Stack:** Go, Chi router, pgx (PostgreSQL), Gorilla WebSocket, Firebase Auth, Polygon.io API, golang-migrate

### Web -- Next.js

A responsive web app with real-time price updates, interactive charts, and a complete trading workflow.

**Highlights:**
- **Live candlestick charts** (lightweight-charts) with configurable time intervals
- Real-time price updates via WebSocket reflected instantly across all views
- Asset type filtering (Stocks / ETFs / Crypto) with sector groupings
- Full trading flow: market orders, limit/stop order creation, order management
- Portfolio dashboard with P&L breakdown, trade history, and position details
- Onboarding flow explaining real market data and paper trading concept
- Dark theme inspired by GitHub's design system

**Stack:** Next.js 15, React 19, TypeScript, Tailwind CSS, Zustand, lightweight-charts

### iOS -- SwiftUI

A native iOS app with a Robinhood-inspired trading interface and full feature parity with the web.

**Highlights:**
- **Robinhood-style full-screen trade flow** with share/dollar input modes, review confirmation, and success animation
- Native Charts framework for price visualization
- Real-time WebSocket price updates with live UI transitions
- ETF detail view showing weighted constituent holdings
- Keychain-based auth token persistence

**Stack:** SwiftUI, iOS 18+, Swift 6, MVVM with @Observable, async/await, SPM

### Android -- Jetpack Compose

A native Android app with Material 3 design, Hilt dependency injection, and the same trading workflow.

**Highlights:**
- Full MVVM architecture with Hilt DI, Retrofit networking, and DataStore persistence
- Market view with live data, search filtering, and navigation to stock detail
- Trade execution with buy/sell, quantity input, and real-time price display
- Portfolio with P&L cards, holdings breakdown, and achievement tracking

**Stack:** Jetpack Compose, Kotlin, Hilt, Retrofit + Moshi, Material 3, DataStore

---

## Architecture

```
 iOS (SwiftUI)     Web (Next.js)     Android (Compose)
       |                |                  |
       +------------ REST API + WebSocket -+
                        |
              +---------+---------+
              |    Go Backend     |
              |                   |
              |  PriceProvider    |  interface
              |    /        \    |
              | Polygon.io  Simulation |
              |  (live)     (dev)      |
              |                   |
              |  Order Matcher   |  limit / stop / stop-limit
              |  8+ Workers      |  price history, alerts,
              |                   |  leaderboard, achievements,
              |                   |  options pricing, stock sync
              +---------+---------+
                        |
                   PostgreSQL
```

## Assets

All data -- company names, sectors, descriptions, and prices -- sourced from Polygon.io.

| Type | Count | Examples |
|------|-------|---------|
| **Stocks** | 37 | AAPL, MSFT, NVDA, GOOGL, TSLA, META, JPM, LLY, UNH |
| **ETFs** | 5 | SPY, QQQ, DIA, IWM, VTI |
| **Crypto** | 4 | BTC, ETH, SOL, DOGE |

Sectors are derived from SIC codes: Technology, Healthcare, Financial, Energy, Consumer, Industrial.

## Features

| Feature | Description |
|---------|------------|
| Real market data | Live prices from Polygon.io (NYSE, NASDAQ, crypto exchanges) |
| Paper trading | $100K virtual cash, real prices, zero risk |
| Advanced orders | Limit, stop, and stop-limit orders with automatic matching |
| Options trading | Black-Scholes pricing, Greeks, chain generation |
| Portfolio tracking | Net worth, cash, invested value, per-position P&L |
| Market hours | Real US market schedule with pre-market/after-hours awareness |
| Leaderboard | Daily, weekly, and all-time rankings by net worth |
| Achievements | 20 unlockables across trading, portfolio, social, streak, and skill categories |
| Daily challenges | Auto-generated challenges with cash rewards |
| Price alerts | Above/below alerts with real-time WebSocket notifications |
| Candlestick charts | Interactive OHLCV charts at multiple intervals |
| Guest accounts | Start trading immediately with no sign-up |

## Running Locally

### Quick start with Docker

```bash
cd deploy && cp .env.example .env
# Edit .env to add your POLYGON_API_KEY and set MARKET_DATA_SOURCE=polygon
docker compose up -d
```

### Manual setup

```bash
# Start Postgres (if not running)
cd deploy && docker compose up -d postgres

# Backend
cd backend
export DATABASE_URL="postgres://mockstarket:mockstarket_dev@localhost:5432/mockstarket?sslmode=disable"

go run cmd/migrate/main.go up
POLYGON_API_KEY=your_key go run cmd/seed/main.go  # fetches real data from Polygon
DATABASE_URL=$DATABASE_URL MARKET_DATA_SOURCE=polygon POLYGON_API_KEY=your_key DEV_MODE=true go run cmd/server/main.go

# Without Polygon (simulation mode -- no API key needed)
go run cmd/server/main.go

# Web
cd web && npm install && npm run dev

# iOS
open ios/MockStarket.xcodeproj

# Android
# Open android/ in Android Studio
```

| Service | Port |
|---------|------|
| Backend API | 8080 |
| PostgreSQL | 5432 |
| Web (Next.js) | 3000 |
| Nginx (Docker) | 80 |

### Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MARKET_DATA_SOURCE` | `simulation` | `polygon` for live data, `simulation` for fake data |
| `POLYGON_API_KEY` | -- | Required for `polygon` mode. Get one at [polygon.io](https://polygon.io) |
| `POLYGON_WS_ENABLED` | `false` | `true` for real-time WebSocket streaming (paid tier) |
| `POLYGON_POLL_INTERVAL_MS` | `30000` | How often to poll for price updates (ms) |
| `DEV_MODE` | `true` | Bypasses Firebase auth (tokens treated as user IDs) |
| `DATABASE_URL` | -- | PostgreSQL connection string |

### Tests

```bash
cd backend && go test ./...
cd web && npm run type-check
```

## Project Structure

```
mock-starket/
  backend/        Go API server, Polygon.io client, simulation engine
    internal/
      market/     PriceProvider interface and shared types
      polygon/    Polygon.io REST client, WebSocket, feed, scheduler
      simulation/ Stochastic price simulation (dev/fallback mode)
      handler/    HTTP handlers (30+ endpoints)
      worker/     Background workers (price sync, alerts, leaderboard, options)
      service/    Business logic (trading, achievements, challenges)
  web/            Next.js frontend (charts, trading, portfolio)
  ios/            SwiftUI app (Robinhood-style trading)
  android/        Compose app (Material 3)
  deploy/         Docker Compose, nginx, Postgres config
  .github/        CI/CD workflows, PR/issue templates
```

## License

MIT

## Author

**Luke Solomon** -- [GitHub](https://github.com/nsluke)
