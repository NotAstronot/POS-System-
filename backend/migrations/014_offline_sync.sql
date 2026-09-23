-- --------------------------------------------------------
-- Migration 014: Offline sync (idempotency untuk order offline)
-- --------------------------------------------------------

ALTER TABLE orders ADD COLUMN IF NOT EXISTS client_order_id VARCHAR(64) DEFAULT '';
CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_client_order_id ON orders(client_order_id) WHERE client_order_id <> '';