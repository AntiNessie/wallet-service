CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE wallets (
                         id UUID PRIMARY KEY,
                         balance DECIMAL(15,2) NOT NULL DEFAULT 0.00,
                         created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                         updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);