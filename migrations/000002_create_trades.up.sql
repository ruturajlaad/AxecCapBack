CREATE TABLE IF NOT EXISTS trades (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(10) CHECK (side IN ('buy', 'sell')),
    price NUMERIC(12,4) NOT NULL,
    qty NUMERIC(12,4) NOT NULL,
    pnl NUMERIC(12,4) DEFAULT 0.0,
    created_at TIMESTAMP DEFAULT NOW()
);
