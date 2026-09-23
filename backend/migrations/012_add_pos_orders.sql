-- --------------------------------------------------------
-- Migration 012: POS order integration (stock deduction + daily revenue)
-- --------------------------------------------------------

-- Extend orders with POS fields
ALTER TABLE orders ADD COLUMN IF NOT EXISTS order_number VARCHAR(20) DEFAULT '';
ALTER TABLE orders ADD COLUMN IF NOT EXISTS order_type VARCHAR(20) NOT NULL DEFAULT 'dine_in';
ALTER TABLE orders ADD COLUMN IF NOT EXISTS table_number VARCHAR(20) DEFAULT '';
ALTER TABLE orders ADD COLUMN IF NOT EXISTS customer_name VARCHAR(255) DEFAULT '';
ALTER TABLE orders ADD COLUMN IF NOT EXISTS customer_phone VARCHAR(30) DEFAULT '';
ALTER TABLE orders ADD COLUMN IF NOT EXISTS discount_amount NUMERIC(15,2) NOT NULL DEFAULT 0;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS payment_method VARCHAR(20) NOT NULL DEFAULT 'cash';
ALTER TABLE orders ADD COLUMN IF NOT EXISTS outlet_id VARCHAR(64) DEFAULT '';

-- Allow orders without an active shift (POS must keep working)
ALTER TABLE orders ALTER COLUMN shift_id DROP NOT NULL;

-- Extend order_items with POS fields
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS variant_label VARCHAR(255) DEFAULT '';
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS note TEXT DEFAULT '';

-- Daily revenue booking (real-time, upserted on every POS sale)
CREATE TABLE IF NOT EXISTS daily_revenue (
    date DATE PRIMARY KEY,
    total_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    transaction_count INT NOT NULL DEFAULT 0,
    product_count INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
CREATE INDEX IF NOT EXISTS idx_orders_outlet_id ON orders(outlet_id);
