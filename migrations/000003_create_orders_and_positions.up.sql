-- Limit orders table (pending trades before execution)
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(10) CHECK (side IN ('buy', 'sell')) NOT NULL,
    limit_price NUMERIC(12,4) NOT NULL,
    contracts NUMERIC(12,4) NOT NULL,
    leverage INT DEFAULT 100,
    stop_loss NUMERIC(12,4),
    take_profit NUMERIC(12,4),
    status VARCHAR(20) DEFAULT 'pending', -- pending, executed, cancelled
    created_at TIMESTAMP DEFAULT NOW()
);

-- Active positions table (executed trades that are open)
CREATE TABLE IF NOT EXISTS positions (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE SET NULL,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(10) CHECK (side IN ('buy', 'sell')) NOT NULL,
    entry_price NUMERIC(12,4) NOT NULL,
    contracts NUMERIC(12,4) NOT NULL,
    leverage INT DEFAULT 100,
    stop_loss NUMERIC(12,4),
    take_profit NUMERIC(12,4),
    status VARCHAR(20) DEFAULT 'open', -- open, closed
    opened_at TIMESTAMP DEFAULT NOW(),
    closed_at TIMESTAMP,
    exit_price NUMERIC(12,4),
    realized_pnl NUMERIC(12,4) DEFAULT 0.0
);

-- Optional: add margin reservation tracking for users
ALTER TABLE users ADD COLUMN IF NOT EXISTS reserved_margin NUMERIC(12,2) DEFAULT 0.0;
