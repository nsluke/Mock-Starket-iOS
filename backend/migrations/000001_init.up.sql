-- Mock Starket Initial Schema

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users (Firebase UID as external identity)
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firebase_uid    TEXT UNIQUE NOT NULL,
    display_name    TEXT NOT NULL,
    avatar_url      TEXT,
    is_guest        BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    last_login_at   TIMESTAMPTZ,
    login_streak    INT DEFAULT 0,
    longest_streak  INT DEFAULT 0
);

-- Stocks (seeded by the application)
CREATE TABLE stocks (
    ticker          TEXT PRIMARY KEY,
    name            TEXT NOT NULL,
    sector          TEXT NOT NULL,
    base_price      DECIMAL(12,4) NOT NULL,
    current_price   DECIMAL(12,4) NOT NULL,
    day_open        DECIMAL(12,4),
    day_high        DECIMAL(12,4),
    day_low         DECIMAL(12,4),
    prev_close      DECIMAL(12,4),
    volume          BIGINT DEFAULT 0,
    volatility      DECIMAL(6,4) DEFAULT 0.02,
    drift           DECIMAL(6,4) DEFAULT 0.0,
    mean_reversion  DECIMAL(6,4) DEFAULT 0.1,
    description     TEXT,
    logo_url        TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Portfolios (one per user)
CREATE TABLE portfolios (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    cash            DECIMAL(14,4) NOT NULL DEFAULT 100000.0000,
    net_worth       DECIMAL(14,4) NOT NULL DEFAULT 100000.0000,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Holdings (positions in stocks)
CREATE TABLE holdings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    portfolio_id    UUID REFERENCES portfolios(id) ON DELETE CASCADE,
    ticker          TEXT REFERENCES stocks(ticker),
    shares          INT NOT NULL DEFAULT 0,
    avg_cost        DECIMAL(12,4) NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(portfolio_id, ticker)
);

-- Trades (immutable log)
CREATE TABLE trades (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    ticker          TEXT REFERENCES stocks(ticker),
    side            TEXT NOT NULL CHECK (side IN ('buy', 'sell')),
    shares          INT NOT NULL,
    price           DECIMAL(12,4) NOT NULL,
    total           DECIMAL(14,4) NOT NULL,
    order_id        UUID,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Orders (limit, stop, stop-limit)
CREATE TABLE orders (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    ticker          TEXT REFERENCES stocks(ticker),
    side            TEXT NOT NULL CHECK (side IN ('buy', 'sell')),
    order_type      TEXT NOT NULL CHECK (order_type IN ('limit', 'stop', 'stop_limit')),
    shares          INT NOT NULL,
    limit_price     DECIMAL(12,4),
    stop_price      DECIMAL(12,4),
    status          TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'filled', 'cancelled', 'expired')),
    filled_price    DECIMAL(12,4),
    filled_at       TIMESTAMPTZ,
    expires_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Price history (for charts)
CREATE TABLE price_history (
    id              BIGSERIAL PRIMARY KEY,
    ticker          TEXT REFERENCES stocks(ticker),
    price           DECIMAL(12,4) NOT NULL,
    open_price      DECIMAL(12,4),
    high            DECIMAL(12,4),
    low             DECIMAL(12,4),
    close_price     DECIMAL(12,4),
    volume          BIGINT DEFAULT 0,
    interval        TEXT NOT NULL CHECK (interval IN ('1s', '1m', '5m', '1h', '1d')),
    recorded_at     TIMESTAMPTZ NOT NULL,
    UNIQUE(ticker, interval, recorded_at)
);
CREATE INDEX idx_price_history_lookup ON price_history(ticker, interval, recorded_at DESC);

-- Portfolio value history
CREATE TABLE portfolio_history (
    id              BIGSERIAL PRIMARY KEY,
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    net_worth       DECIMAL(14,4) NOT NULL,
    cash            DECIMAL(14,4) NOT NULL,
    recorded_at     TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_portfolio_history_lookup ON portfolio_history(user_id, recorded_at DESC);

-- Leaderboard (materialized, recomputed periodically)
CREATE TABLE leaderboard (
    id              BIGSERIAL PRIMARY KEY,
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    display_name    TEXT NOT NULL,
    net_worth       DECIMAL(14,4) NOT NULL,
    total_return    DECIMAL(8,4),
    rank            INT NOT NULL,
    period          TEXT NOT NULL CHECK (period IN ('daily', 'weekly', 'alltime')),
    computed_at     TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, period, computed_at)
);
CREATE INDEX idx_leaderboard_rank ON leaderboard(period, rank);

-- Achievements (definitions)
CREATE TABLE achievements (
    id              TEXT PRIMARY KEY,
    name            TEXT NOT NULL,
    description     TEXT NOT NULL,
    icon            TEXT NOT NULL,
    category        TEXT NOT NULL,
    criteria_json   JSONB
);

-- User achievements (earned)
CREATE TABLE user_achievements (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    achievement_id  TEXT REFERENCES achievements(id),
    earned_at       TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, achievement_id)
);

-- Price alerts
CREATE TABLE price_alerts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    ticker          TEXT REFERENCES stocks(ticker),
    condition       TEXT NOT NULL CHECK (condition IN ('above', 'below')),
    target_price    DECIMAL(12,4) NOT NULL,
    triggered       BOOLEAN DEFAULT FALSE,
    triggered_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Daily challenges
CREATE TABLE daily_challenges (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date            DATE UNIQUE NOT NULL,
    challenge_type  TEXT NOT NULL,
    description     TEXT NOT NULL,
    target_json     JSONB NOT NULL,
    reward_cash     DECIMAL(10,4) DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- User challenge completions
CREATE TABLE user_challenges (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    challenge_id    UUID REFERENCES daily_challenges(id),
    completed       BOOLEAN DEFAULT FALSE,
    completed_at    TIMESTAMPTZ,
    claimed         BOOLEAN DEFAULT FALSE,
    UNIQUE(user_id, challenge_id)
);

-- Watchlist
CREATE TABLE watchlist (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    ticker          TEXT REFERENCES stocks(ticker),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, ticker)
);
