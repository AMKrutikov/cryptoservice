BEGIN;
    CREATE SCHEMA IF NOT EXISTS crypto;

    CREATE TABLE IF NOT EXISTS crypto.coins (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        price NUMERIC(18, 8) NOT NULL,
        actual_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

    CREATE INDEX IF NOT EXISTS idx_crypto_coins_title_actual_at
    ON crypto.coins (title, actual_at DESC);
COMMIT;