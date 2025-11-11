CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username varchar(50) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    balance NUMERIC(12,2) DEFAULT 10000.00,
    created_at TIMESTAMP DEFAULT NOW()
);