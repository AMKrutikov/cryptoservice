BEGIN;
    DROP INDEX IF EXISTS crypto.idx_crypto_coins_title_actual_at;
    DROP TABLE IF EXISTS crypto.coins;
    DROP SCHEMA IF EXISTS crypto CASCADE;
COMMIT;