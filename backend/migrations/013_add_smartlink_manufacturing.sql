-- --------------------------------------------------------
-- Migration 013: Store & SmartLink + Modul Manufaktur
-- e-Commerce (marketplace), e-Banking (mutasi bank),
-- e-Faktur / Pajak, dan Manufaktur (BOM, Work Order, HPP)
-- --------------------------------------------------------

-- NPWP untuk keperluan ekspor e-Faktur / PPh
ALTER TABLE customers ADD COLUMN IF NOT EXISTS npwp VARCHAR(30) NOT NULL DEFAULT '';
ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS npwp VARCHAR(30) NOT NULL DEFAULT '';

-- ============ SmartLink e-Commerce (Marketplace) ============
CREATE TABLE IF NOT EXISTS marketplace_connections (
    id BIGSERIAL PRIMARY KEY,
    platform VARCHAR(50) NOT NULL,
    shop_name VARCHAR(255) NOT NULL,
    api_token TEXT NOT NULL DEFAULT '',
    account_name VARCHAR(100) NOT NULL DEFAULT '',
    customer_id BIGINT REFERENCES customers(id),
    sales_category_id BIGINT REFERENCES sales_categories(id),
    shipping_fee_account VARCHAR(100) NOT NULL DEFAULT 'Biaya Ongkir',
    commission_account VARCHAR(100) NOT NULL DEFAULT 'Komisi Platform',
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_sync_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS marketplace_orders (
    id BIGSERIAL PRIMARY KEY,
    connection_id BIGINT NOT NULL REFERENCES marketplace_connections(id) ON DELETE CASCADE,
    marketplace_order_id VARCHAR(100) NOT NULL,
    order_date DATE NOT NULL DEFAULT CURRENT_DATE,
    customer_name VARCHAR(255) NOT NULL DEFAULT '',
    product_id BIGINT REFERENCES products(id),
    product_name VARCHAR(255) NOT NULL DEFAULT '',
    quantity NUMERIC(15,2) NOT NULL DEFAULT 0,
    unit_price NUMERIC(15,2) NOT NULL DEFAULT 0,
    subtotal NUMERIC(15,2) NOT NULL DEFAULT 0,
    shipping_fee NUMERIC(15,2) NOT NULL DEFAULT 0,
    platform_fee NUMERIC(15,2) NOT NULL DEFAULT 0,
    grand_total NUMERIC(15,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'imported',
    sales_order_id BIGINT REFERENCES sales_orders(id),
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(connection_id, marketplace_order_id)
);

-- ============ SmartLink e-Banking (Impor Mutasi Bank) ============
CREATE TABLE IF NOT EXISTS bank_statement_imports (
    id BIGSERIAL PRIMARY KEY,
    import_number VARCHAR(50) NOT NULL UNIQUE,
    account_name VARCHAR(100) NOT NULL,
    account_type VARCHAR(20) NOT NULL DEFAULT 'bank',
    statement_date DATE NOT NULL DEFAULT CURRENT_DATE,
    source VARCHAR(20) NOT NULL DEFAULT 'manual',
    file_name VARCHAR(255) NOT NULL DEFAULT '',
    total_rows INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'imported',
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bank_statement_lines (
    id BIGSERIAL PRIMARY KEY,
    import_id BIGINT NOT NULL REFERENCES bank_statement_imports(id) ON DELETE CASCADE,
    transaction_date DATE NOT NULL DEFAULT CURRENT_DATE,
    description VARCHAR(255) NOT NULL DEFAULT '',
    reference VARCHAR(100) NOT NULL DEFAULT '',
    amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    type VARCHAR(10) NOT NULL DEFAULT 'debit',
    is_reconciled BOOLEAN NOT NULL DEFAULT false,
    reconciliation_id BIGINT REFERENCES bank_reconciliations(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============ SmartLink Tax / e-Faktur ============
CREATE TABLE IF NOT EXISTS efaktur_exports (
    id BIGSERIAL PRIMARY KEY,
    export_number VARCHAR(50) NOT NULL UNIQUE,
    export_type VARCHAR(20) NOT NULL,
    period VARCHAR(7) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'generated',
    total_rows INT NOT NULL DEFAULT 0,
    total_dpp NUMERIC(15,2) NOT NULL DEFAULT 0,
    total_tax NUMERIC(15,2) NOT NULL DEFAULT 0,
    file_path VARCHAR(255) NOT NULL DEFAULT '',
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS efaktur_lines (
    id BIGSERIAL PRIMARY KEY,
    export_id BIGINT NOT NULL REFERENCES efaktur_exports(id) ON DELETE CASCADE,
    reference_number VARCHAR(50) NOT NULL DEFAULT '',
    party_name VARCHAR(255) NOT NULL DEFAULT '',
    npwp VARCHAR(30) NOT NULL DEFAULT '',
    line_date DATE,
    taxable_base NUMERIC(15,2) NOT NULL DEFAULT 0,
    tax_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'valid'
);

-- ============ Modul Manufaktur (BOM & Work Order) ============
CREATE TABLE IF NOT EXISTS bill_of_materials (
    id BIGSERIAL PRIMARY KEY,
    bom_number VARCHAR(50) NOT NULL UNIQUE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    product_name VARCHAR(255) NOT NULL DEFAULT '',
    quantity_output NUMERIC(15,2) NOT NULL DEFAULT 1,
    unit VARCHAR(50) NOT NULL DEFAULT 'pcs',
    overhead_cost NUMERIC(15,2) NOT NULL DEFAULT 0,
    cost_per_unit NUMERIC(15,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    notes TEXT NOT NULL DEFAULT '',
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bom_items (
    id BIGSERIAL PRIMARY KEY,
    bom_id BIGINT NOT NULL REFERENCES bill_of_materials(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    product_name VARCHAR(255) NOT NULL DEFAULT '',
    quantity_required NUMERIC(15,2) NOT NULL DEFAULT 0,
    unit VARCHAR(50) NOT NULL DEFAULT 'pcs',
    cost_per_unit NUMERIC(15,2) NOT NULL DEFAULT 0,
    estimated_cost NUMERIC(15,2) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS work_orders (
    id BIGSERIAL PRIMARY KEY,
    work_order_number VARCHAR(50) NOT NULL UNIQUE,
    bom_id BIGINT NOT NULL REFERENCES bill_of_materials(id),
    bom_number VARCHAR(50) NOT NULL DEFAULT '',
    product_id BIGINT NOT NULL REFERENCES products(id),
    product_name VARCHAR(255) NOT NULL DEFAULT '',
    quantity NUMERIC(15,2) NOT NULL DEFAULT 0,
    warehouse_id BIGINT REFERENCES warehouses(id),
    scheduled_date DATE NOT NULL DEFAULT CURRENT_DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'planned',
    actual_cost NUMERIC(15,2) NOT NULL DEFAULT 0,
    notes TEXT NOT NULL DEFAULT '',
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS work_order_items (
    id BIGSERIAL PRIMARY KEY,
    work_order_id BIGINT NOT NULL REFERENCES work_orders(id) ON DELETE CASCADE,
    bom_item_id BIGINT,
    product_id BIGINT NOT NULL REFERENCES products(id),
    product_name VARCHAR(255) NOT NULL DEFAULT '',
    quantity_required NUMERIC(15,2) NOT NULL DEFAULT 0,
    unit VARCHAR(50) NOT NULL DEFAULT 'pcs',
    cost_per_unit NUMERIC(15,2) NOT NULL DEFAULT 0,
    subtotal NUMERIC(15,2) NOT NULL DEFAULT 0
);

-- ============ Indexes ============
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_connection ON marketplace_orders(connection_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_status ON marketplace_orders(status);
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_sales_order ON marketplace_orders(sales_order_id);
CREATE INDEX IF NOT EXISTS idx_bank_statement_lines_import ON bank_statement_lines(import_id);
CREATE INDEX IF NOT EXISTS idx_bank_statement_lines_reconciled ON bank_statement_lines(is_reconciled);
CREATE INDEX IF NOT EXISTS idx_efaktur_exports_period ON efaktur_exports(period);
CREATE INDEX IF NOT EXISTS idx_efaktur_lines_export ON efaktur_lines(export_id);
CREATE INDEX IF NOT EXISTS idx_bom_items_bom ON bom_items(bom_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_status ON work_orders(status);
CREATE INDEX IF NOT EXISTS idx_work_order_items_work_order ON work_order_items(work_order_id);

-- ============ Permissions ============
INSERT INTO permissions (name, description) VALUES
    ('smartlink_manage', 'Kelola Store & SmartLink (e-Commerce, e-Banking, e-Faktur)'),
    ('manufacturing_manage', 'Kelola Modul Manufaktur (BOM & Perintah Kerja)')
ON CONFLICT (name) DO NOTHING;
