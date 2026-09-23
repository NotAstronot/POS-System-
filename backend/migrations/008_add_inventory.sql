-- Inventory / Gudang Module

ALTER TABLE products ADD COLUMN IF NOT EXISTS barcode VARCHAR(100) NOT NULL DEFAULT '';
ALTER TABLE products ADD COLUMN IF NOT EXISTS purchase_price NUMERIC(15,2) NOT NULL DEFAULT 0;
ALTER TABLE products ADD COLUMN IF NOT EXISTS min_stock INT NOT NULL DEFAULT 0;
ALTER TABLE products ADD COLUMN IF NOT EXISTS is_service BOOLEAN NOT NULL DEFAULT false;

CREATE TABLE IF NOT EXISTS item_variants (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    barcode VARCHAR(100) NOT NULL DEFAULT '',
    additional_price NUMERIC(15,2) NOT NULL DEFAULT 0,
    stock INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS warehouses (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    branch_id BIGINT REFERENCES branches(id),
    address TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS stock (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    quantity NUMERIC(15,2) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(product_id, warehouse_id)
);

CREATE TABLE IF NOT EXISTS stock_movements (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id),
    warehouse_id BIGINT NOT NULL REFERENCES warehouses(id),
    quantity NUMERIC(15,2) NOT NULL DEFAULT 0,
    movement_type VARCHAR(50) NOT NULL DEFAULT 'other',
    reference_type VARCHAR(50) NOT NULL DEFAULT '',
    reference_id BIGINT,
    note TEXT NOT NULL DEFAULT '',
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS stock_transfers (
    id BIGSERIAL PRIMARY KEY,
    transfer_number VARCHAR(50) NOT NULL UNIQUE,
    from_warehouse_id BIGINT NOT NULL REFERENCES warehouses(id),
    to_warehouse_id BIGINT NOT NULL REFERENCES warehouses(id),
    transfer_date DATE NOT NULL DEFAULT CURRENT_DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    notes TEXT NOT NULL DEFAULT '',
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS stock_transfer_items (
    id BIGSERIAL PRIMARY KEY,
    transfer_id BIGINT NOT NULL REFERENCES stock_transfers(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    product_name VARCHAR(255) NOT NULL DEFAULT '',
    quantity NUMERIC(15,2) NOT NULL DEFAULT 0,
    unit VARCHAR(50) NOT NULL DEFAULT 'pcs',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS stock_adjustments (
    id BIGSERIAL PRIMARY KEY,
    adjustment_number VARCHAR(50) NOT NULL UNIQUE,
    warehouse_id BIGINT NOT NULL REFERENCES warehouses(id),
    adjustment_date DATE NOT NULL DEFAULT CURRENT_DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    notes TEXT NOT NULL DEFAULT '',
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS stock_adjustment_items (
    id BIGSERIAL PRIMARY KEY,
    adjustment_id BIGINT NOT NULL REFERENCES stock_adjustments(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    product_name VARCHAR(255) NOT NULL DEFAULT '',
    system_qty NUMERIC(15,2) NOT NULL DEFAULT 0,
    actual_qty NUMERIC(15,2) NOT NULL DEFAULT 0,
    difference NUMERIC(15,2) NOT NULL DEFAULT 0,
    reason TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS price_changes (
    id BIGSERIAL PRIMARY KEY,
    price_change_number VARCHAR(50) NOT NULL UNIQUE,
    change_date DATE NOT NULL DEFAULT CURRENT_DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'applied',
    notes TEXT NOT NULL DEFAULT '',
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS price_change_items (
    id BIGSERIAL PRIMARY KEY,
    price_change_id BIGINT NOT NULL REFERENCES price_changes(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    product_name VARCHAR(255) NOT NULL DEFAULT '',
    old_price NUMERIC(15,2) NOT NULL DEFAULT 0,
    new_price NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO permissions (name, description) VALUES
    ('inventory_manage', 'Kelola Persediaan & Gudang')
ON CONFLICT (name) DO NOTHING;

INSERT INTO warehouses (code, name, address)
SELECT 'WH-001', 'Gudang Utama', 'Gudang utama sistem'
WHERE NOT EXISTS (SELECT 1 FROM warehouses)
ON CONFLICT (code) DO NOTHING;

INSERT INTO stock (product_id, warehouse_id, quantity)
SELECT p.id, w.id, p.stock FROM products p, warehouses w
WHERE w.code = 'WH-001'
ON CONFLICT (product_id, warehouse_id) DO NOTHING;
