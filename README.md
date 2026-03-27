# Mock Starket

**A real-time stock market simulator built from the ground up as a full-stack, multi-platform application.**

Trade 31 assets across stocks, ETFs, crypto, and commodities with $100,000 in virtual cash. Prices move in real-time using a custom simulation engine. Compete on the leaderboard, earn achievements, and complete daily challenges.

Built solo as a portfolio project to demonstrate full-stack architecture, real-time systems, and native mobile development across iOS, Android, and web.

<!-- Screenshots go here once captured -->
<!-- ![Market View](docs/screenshots/market.png) -->

---

## What I Built

### Backend -- Go

A production-grade REST API and WebSocket server handling 30+ endpoints, real-time price simulation, and 8 concurrent background workers.

**Highlights:**
- Custom stock price simulation using **Geometric Brownian Motion** with mean reversion, sector correlation, and random market events (earnings surprises, sector shifts, Fed announcements)
- **Order matching engine** that evaluates limit, stop, and stop-limit orders against live prices every simulation tick
- **Real-time WebSocket** broadcasting price updates, trade confirmations, and alert notifications to connected clients
- Background workers for price history recording (OHLCV at 4 intervals), leaderboard computation, achievement evaluation, daily challenge generation, and price alert monitoring
- **42 unit tests** covering order matching logic, OHLCV aggregation, middleware (rate limiting, auth, panic recovery), config loading, and HTTP handlers

**Stack:** Go 1.23, Chi router, pgx (PostgreSQL), Gorilla WebSocket, Firebase Auth, golang-migrate

### Web -- Next.js

A 12-page responsive web app with real-time price updates, interactive charts, and a complete trading workflow.

**Highlights:**
- **Live candlestick charts** (lightweight-charts) with configurable time intervals
- Real-time price updates via WebSocket reflected instantly across all views
- Asset type filtering (Stocks / ETFs / Crypto / Commodities)
- Full trading flow: market orders, limit/stop order creation, order management
- Portfolio dashboard with P&L breakdown, trade history, and position details
- Price alerts, daily challenges, leaderboard, and achievement tracking
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
- Asset type badges and user position display on stock detail

**Stack:** SwiftUI, iOS 18+, Swift 6, MVVM with @Observable, async/await, SPM

### Android -- Jetpack Compose

A native Android app with Material 3 design, Hilt dependency injection, and the same trading workflow.

**Highlights:**
- Full MVVM architecture with Hilt DI, Retrofit networking, and DataStore persistence
- Market view with live data, search filtering, and navigation to stock detail
- Trade execution with buy/sell, quantity input, and real-time price display
- Portfolio with P&L cards, holdings breakdown, and achievement tracking
- Material 3 dark theme matching the web and iOS color palette

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
              |  Simulation Engine|  GBM + mean reversion
              |  Order Matcher   |  limit / stop / stop-limit
              |  8 Workers       |  price history, alerts,
              |                   |  leaderboard, achievements,
              |                   |  challenges, stock sync
              +---------+---------+
                        |
                   PostgreSQL
```

## Assets

| Type | Count | Examples |
|------|-------|---------|
| **Stocks** | 20 | Pied Piper, Stark Industries, Dunder Mifflin, Krusty Krab |
| **Crypto** | 4 | Bitcoin, Ethereum, Solana, Dogecoin |
| **Commodities** | 3 | Gold, Silver, Crude Oil |
| **ETFs** | 4 | Total Market, Tech, Defense, Crypto |

ETFs track their underlying holdings with configurable weights and derive their price from constituent assets.

## Features

| Feature | Description |
|---------|------------|
| Real-time trading | Market orders execute instantly at simulated prices |
| Advanced orders | Limit, stop, and stop-limit orders with automatic matching |
| Portfolio tracking | Net worth, cash, invested value, per-position P&L |
| Price simulation | Geometric Brownian Motion with drift, volatility, and mean reversion |
| Market events | Random earnings surprises, sector events, and macro announcements |
| Leaderboard | Daily, weekly, and all-time rankings by net worth |
| Achievements | 20 unlockables across trading, portfolio, social, streak, and skill categories |
| Daily challenges | Auto-generated challenges with cash rewards |
| Price alerts | Above/below alerts with real-time WebSocket notifications |
| Candlestick charts | Interactive OHLCV charts at 1-second to 1-hour intervals |
| ETF holdings | View constituent assets and their weights |
| Guest accounts | Start trading immediately with no sign-up |

## Running Locally

**Quick start with Docker:**
```bash
cd deploy && cp .env.example .env && docker compose up -d
```

**Manual setup:**
```bash
# Backend
cd backend
export DATABASE_URL="postgres://user:pass@localhost:5432/mockstarket?sslmode=disable"
go run cmd/migrate/main.go up
go run cmd/seed/main.go
go run cmd/server/main.go

# Web
cd web && npm install && npm run dev

# iOS
open ios/MockStarket.xcodeproj

# Android
# Open android/ in Android Studio
```

**Tests:**
```bash
cd backend && go test ./...    # 42 tests
cd web && npm run type-check   # TypeScript validation
```

## Project Structure

```
mock-starket/
  backend/        Go API server (30+ endpoints, 8 workers, simulation engine)
  web/            Next.js frontend (12 pages, charts, WebSocket)
  ios/            SwiftUI app (28 Swift files, MVVM, Robinhood-style trading)
  android/        Compose app (23 Kotlin files, Hilt, Material 3)
  deploy/         Docker Compose, nginx, Postgres
  .github/        7 CI/CD workflows, PR/issue templates
```

## License

MIT

## Author

**Luke Solomon** -- [GitHub](https://github.com/ares42)
