-- Option contracts: catalog of available options, priced by the engine each tick.
CREATE TABLE option_contracts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker          TEXT NOT NULL REFERENCES stocks(ticker) ON DELETE CASCADE,
    option_type     TEXT NOT NULL CHECK (option_type IN ('call', 'put')),
    strike_price    DECIMAL(12,4) NOT NULL,
    expiration      TIMESTAMPTZ NOT NULL,
    contract_symbol TEXT UNIQUE NOT NULL,

    -- Pricing (updated every tick)
    bid_price       DECIMAL(12,4) NOT NULL DEFAULT 0,
    ask_price       DECIMAL(12,4) NOT NULL DEFAULT 0,
    last_price      DECIMAL(12,4) NOT NULL DEFAULT 0,
    mark_price      DECIMAL(12,4) NOT NULL DEFAULT 0,
    open_interest   INT NOT NULL DEFAULT 0,
    volume          INT NOT NULL DEFAULT 0,
    implied_vol     DECIMAL(8,6) NOT NULL DEFAULT 0.3,

    -- Greeks
    delta           DECIMAL(10,6) NOT NULL DEFAULT 0,
    gamma           DECIMAL(10,6) NOT NULL DEFAULT 0,
    theta           DECIMAL(10,6) NOT NULL DEFAULT 0,
    vega            DECIMAL(10,6) NOT NULL DEFAULT 0,
    rho             DECIMAL(10,6) NOT NULL DEFAULT 0,

    -- Status
    status          TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'expired', 'exercised')),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_option_contracts_chain ON option_contracts(ticker, expiration, option_type, strike_price);
CREATE INDEX idx_option_contracts_expiry ON option_contracts(expiration, status);
CREATE INDEX idx_option_contracts_symbol ON option_contracts(contract_symbol);

-- Option positions: user holdings. Positive quantity = long, negative = short/written.
CREATE TABLE option_positions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    portfolio_id    UUID NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
    contract_id     UUID NOT NULL REFERENCES option_contracts(id),
    quantity        INT NOT NULL,
    avg_cost        DECIMAL(12,4) NOT NULL,
    collateral      DECIMAL(14,4) NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(portfolio_id, contract_id)
);

-- Option trades: immutable trade log.
CREATE TABLE option_trades (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contract_id     UUID NOT NULL REFERENCES option_contracts(id),
    side            TEXT NOT NULL CHECK (side IN ('buy_to_open', 'buy_to_close', 'sell_to_open', 'sell_to_close')),
    quantity        INT NOT NULL,
    price           DECIMAL(12,4) NOT NULL,
    total           DECIMAL(14,4) NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Option orders: pending limit orders for options.
CREATE TABLE option_orders (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contract_id     UUID NOT NULL REFERENCES option_contracts(id),
    side            TEXT NOT NULL CHECK (side IN ('buy_to_open', 'buy_to_close', 'sell_to_open', 'sell_to_close')),
    order_type      TEXT NOT NULL CHECK (order_type IN ('market', 'limit')),
    quantity        INT NOT NULL,
    limit_price     DECIMAL(12,4),
    status          TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'filled', 'cancelled', 'expired')),
    filled_price    DECIMAL(12,4),
    filled_at       TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
