-- Add asset_type to stocks table
ALTER TABLE stocks ADD COLUMN asset_type TEXT NOT NULL DEFAULT 'stock';

-- ETF composition: which underlying assets an ETF tracks
CREATE TABLE etf_holdings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    etf_ticker      TEXT NOT NULL REFERENCES stocks(ticker) ON DELETE CASCADE,
    holding_ticker  TEXT NOT NULL REFERENCES stocks(ticker) ON DELETE CASCADE,
    weight          DECIMAL(6,4) NOT NULL DEFAULT 0.0,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(etf_ticker, holding_ticker)
);

CREATE INDEX idx_etf_holdings_etf ON etf_holdings(etf_ticker);
