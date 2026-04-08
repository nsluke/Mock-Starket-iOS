-- Track which data source populated each stock's prices
ALTER TABLE stocks ADD COLUMN IF NOT EXISTS data_source TEXT NOT NULL DEFAULT 'simulation';
ALTER TABLE stocks ADD COLUMN IF NOT EXISTS last_polygon_update TIMESTAMPTZ;

-- Cache market open/close status
CREATE TABLE IF NOT EXISTS market_status (
    id         SERIAL PRIMARY KEY,
    market     TEXT NOT NULL UNIQUE,     -- "stocks", "crypto", "forex"
    is_open    BOOLEAN NOT NULL DEFAULT false,
    next_open  TIMESTAMPTZ,
    next_close TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
