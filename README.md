# Mock Starket

A full-stack stock market simulation where users trade fictional stocks with fake money in real-time. Built as a mono-repo with a Go backend, Next.js web app, SwiftUI iOS app, and Jetpack Compose Android app.

## Architecture

```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   iOS App   │  │   Web App   │  │ Android App │
│   SwiftUI   │  │   Next.js   │  │   Compose   │
└──────┬──────┘  └──────┬──────┘  └──────┬──────┘
       │                │                │
       └────────────────┼────────────────┘
                        │
              REST API + WebSocket
                        │
               ┌────────┴────────┐
               │   Go Backend    │
               │  Chi + pgx +   │
               │  Gorilla WS    │
               ├─────────────────┤
               │ Simulation Eng. │ ← Geometric Brownian Motion
               │ Order Matching  │ ← Limit/Stop/Stop-Limit
               │ Price History   │ ← OHLCV at 1s/1m/5m/1h
               │ Leaderboard    │ ← Ranked by net worth
               │ Achievements   │ ← 20 achievements
               │ Daily Challenge │ ← Auto-generated daily
               │ Price Alerts   │ ← WebSocket notifications
               └────────┬────────┘
                        │
               ┌────────┴────────┐
               │   PostgreSQL    │
               └─────────────────┘
```

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go 1.23, Chi router, pgx, Gorilla WebSocket, Firebase Auth |
| **Web** | Next.js 15, React 19, TypeScript, Tailwind CSS, Zustand, lightweight-charts |
| **iOS** | SwiftUI, iOS 18+, Swift 6, MVVM, async/await, Charts framework |
| **Android** | Jetpack Compose, Material 3, Hilt, Retrofit, Kotlin Coroutines |
| **Database** | PostgreSQL with golang-migrate |
| **Deploy** | Docker Compose, nginx reverse proxy |
| **CI/CD** | GitHub Actions (7 workflows) |

## Features

- **Real-time trading** -- Buy and sell 23 fictional stocks with simulated prices
- **Live price simulation** -- Geometric Brownian Motion with mean reversion, sector correlation, and market events
- **Order types** -- Market, limit, stop, and stop-limit orders with automatic matching
- **Portfolio tracking** -- Net worth, P&L, position breakdown, portfolio history
- **Leaderboard** -- Compete for the top rank (daily, weekly, all-time)
- **Achievements** -- 20 unlockable achievements across trading, portfolio, social, and skill categories
- **Daily challenges** -- Auto-generated challenges with cash rewards
- **Price alerts** -- Set above/below alerts with real-time WebSocket notifications
- **Price charts** -- Interactive candlestick charts with multiple time intervals
- **Guest accounts** -- Start trading immediately, no sign-up required

## Project Structure

```
.
├── backend/          # Go REST API + WebSocket server
│   ├── cmd/          # Server, migrate, seed entrypoints
│   ├── internal/     # Handler, service, repository, simulation, workers
│   └── migrations/   # PostgreSQL schema
├── web/              # Next.js web frontend
│   └── src/          # Pages, components, stores, types
├── ios/              # SwiftUI iOS app
│   └── MockStarket/  # App, Core, Features, Models
├── android/          # Jetpack Compose Android app
│   └── app/src/      # UI, data, domain, DI
├── deploy/           # Docker Compose, nginx, Postgres init
├── .github/          # CI/CD workflows, PR/issue templates
└── scripts/          # Setup and test runner scripts
```

## Getting Started

### Prerequisites

- Go 1.23+
- Node.js 22+
- PostgreSQL 15+
- Docker & Docker Compose (optional, for containerized setup)

### Quick Start (Docker)

```bash
cd deploy
cp .env.example .env
docker compose up -d
```

This starts PostgreSQL, the Go backend, the Next.js web app, and nginx.

### Manual Setup

**Backend:**
```bash
cd backend
export DATABASE_URL="postgres://user:pass@localhost:5432/mockstarket?sslmode=disable"
go run cmd/migrate/main.go up
go run cmd/seed/main.go
go run cmd/server/main.go
```

**Web:**
```bash
cd web
npm install
npm run dev
```

**iOS:**
Open `ios/MockStarket.xcodeproj` in Xcode. The Firebase SPM package will resolve on first open.

**Android:**
Open the `android/` directory in Android Studio. Gradle will sync dependencies automatically.

### Running Tests

```bash
# Backend (42 tests)
cd backend && go test ./...

# Web
cd web && npm run type-check
```

## API Overview

The backend exposes 30+ REST endpoints under `/api/v1/`:

| Group | Endpoints |
|-------|----------|
| Auth | Register, guest login, get/update/delete profile |
| Stocks | List, detail, history, market summary |
| Trading | Execute trades, trade history |
| Orders | Create/list/cancel limit & stop orders |
| Portfolio | Holdings, net worth, portfolio history |
| Leaderboard | Rankings by period |
| Achievements | List all, user progress |
| Challenges | Today's challenge, check progress, claim reward |
| Alerts | Create/list/delete price alerts |
| Watchlist | Add/remove/list |

Real-time updates via WebSocket at `/ws` (price batches, trade confirmations, alert triggers).

## License

MIT

## Author

Luke Solomon
