-- =====================================================
-- Migration 015: Multi-Tenant SaaS Foundation
-- =====================================================

-- ============ Tenants Table ============
CREATE TABLE IF NOT EXISTS tenants (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    domain VARCHAR(255) NOT NULL DEFAULT '',
    logo_url TEXT NOT NULL DEFAULT '',
    subscription_plan VARCHAR(50) NOT NULL DEFAULT 'free',
    subscription_expires_at TIMESTAMPTZ,
    max_users INT NOT NULL DEFAULT 5,
    max_products INT NOT NULL DEFAULT 500,
    max_branches INT NOT NULL DEFAULT 1,
    settings JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============ Add tenant_id to ALL tenant-specific tables ============
-- Categories
ALTER TABLE categories ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);
ALTER TABLE categories DROP CONSTRAINT IF EXISTS categories_name_key;
ALTER TABLE categories ADD CONSTRAINT categories_name_tenant_unique UNIQUE (name, tenant_id);

-- Products
ALTER TABLE products ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Product Availability
ALTER TABLE product_availability ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Users (already has tenant_id as VARCHAR, convert to BIGINT reference)
-- First, create default tenant
INSERT INTO tenants (name, slug, subscription_plan, status) VALUES
    ('Default Tenant', 'default', 'enterprise', 'active')
ON CONFLICT (slug) DO NOTHING;

-- Update users.tenant_id to reference tenants table
-- Set default tenant_id for existing users
DO $$
DECLARE
    default_tenant_id BIGINT;
BEGIN
    SELECT id INTO default_tenant_id FROM tenants WHERE slug = 'default' LIMIT 1;
    UPDATE users SET tenant_id = default_tenant_id::TEXT WHERE tenant_id = '' OR tenant_id IS NULL;
END $$;

-- Shifts
ALTER TABLE shifts ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Orders
ALTER TABLE orders ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Order Items
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Transactions
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Branches
ALTER TABLE branches ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Departments
ALTER TABLE departments ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Employees
ALTER TABLE employees ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Commissions
ALTER TABLE commissions ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Salaries
ALTER TABLE salaries ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Cash Transactions
ALTER TABLE cash_transactions ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Debts
ALTER TABLE debts ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Receivables
ALTER TABLE receivables ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Withdrawals
ALTER TABLE withdrawals ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Suppliers
ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Purchase Orders
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Purchase Items
ALTER TABLE purchase_items ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Purchase Invoices
ALTER TABLE purchase_invoices ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Purchase Payments
ALTER TABLE purchase_payments ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Purchase Returns
ALTER TABLE purchase_returns ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Purchase Return Items
ALTER TABLE purchase_return_items ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Customer Categories
ALTER TABLE customer_categories ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);
ALTER TABLE customer_categories DROP CONSTRAINT IF EXISTS customer_categories_name_key;
ALTER TABLE customer_categories ADD CONSTRAINT customer_categories_name_tenant_unique UNIQUE (name, tenant_id);

-- Sales Categories
ALTER TABLE sales_categories ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);
ALTER TABLE sales_categories DROP CONSTRAINT IF EXISTS sales_categories_name_key;
ALTER TABLE sales_categories ADD CONSTRAINT sales_categories_name_tenant_unique UNIQUE (name, tenant_id);

-- Customers
ALTER TABLE customers ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Sales Quotations
ALTER TABLE sales_quotations ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Sales Quotation Items
ALTER TABLE sales_quotation_items ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Sales Orders
ALTER TABLE sales_orders ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Sales Order Items
ALTER TABLE sales_order_items ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Delivery Orders
ALTER TABLE delivery_orders ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Delivery Order Items
ALTER TABLE delivery_order_items ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Sales Invoices
ALTER TABLE sales_invoices ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Sales Receipts
ALTER TABLE sales_receipts ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Sales Returns
ALTER TABLE sales_returns ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Sales Return Items
ALTER TABLE sales_return_items ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Item Variants
ALTER TABLE item_variants ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Warehouses
ALTER TABLE warehouses ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Stock
ALTER TABLE stock ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Stock Movements
ALTER TABLE stock_movements ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Stock Transfers
ALTER TABLE stock_transfers ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Stock Transfer Items
ALTER TABLE stock_transfer_items ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Stock Adjustments
ALTER TABLE stock_adjustments ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Stock Adjustment Items
ALTER TABLE stock_adjustment_items ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Price Changes
ALTER TABLE price_changes ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Price Change Items
ALTER TABLE price_change_items ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Bank Transfers
ALTER TABLE bank_transfers ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Bank Reconciliations
ALTER TABLE bank_reconciliations ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Bank Reconciliation Items
ALTER TABLE bank_reconciliation_items ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Chart of Accounts
ALTER TABLE chart_of_accounts ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Journal Entries
ALTER TABLE journal_entries ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Journal Entry Items
ALTER TABLE journal_entry_items ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Fixed Assets
ALTER TABLE fixed_assets ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Period Closings
ALTER TABLE period_closings ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- FX Differences
ALTER TABLE fx_differences ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Daily Revenue
ALTER TABLE daily_revenue ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);
ALTER TABLE daily_revenue DROP CONSTRAINT IF EXISTS daily_revenue_pkey;
ALTER TABLE daily_revenue ADD PRIMARY KEY (date, tenant_id);

-- Marketplace Connections
ALTER TABLE marketplace_connections ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Marketplace Orders
ALTER TABLE marketplace_orders ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Bank Statement Imports
ALTER TABLE bank_statement_imports ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Bank Statement Lines
ALTER TABLE bank_statement_lines ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- EFaktur Exports
ALTER TABLE efaktur_exports ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- EFaktur Lines
ALTER TABLE efaktur_lines ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Bill of Materials
ALTER TABLE bill_of_materials ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- BOM Items
ALTER TABLE bom_items ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Work Orders
ALTER TABLE work_orders ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- Work Order Items
ALTER TABLE work_order_items ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id);

-- ============ Backfill tenant_id for existing data ============
DO $$
DECLARE
    default_tenant_id BIGINT;
BEGIN
    SELECT id INTO default_tenant_id FROM tenants WHERE slug = 'default' LIMIT 1;

    UPDATE categories SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE products SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE product_availability SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE shifts SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE orders SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE order_items SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE transactions SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE branches SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE departments SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE employees SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE commissions SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE salaries SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE cash_transactions SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE debts SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE receivables SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE withdrawals SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE suppliers SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE purchase_orders SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE purchase_items SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE purchase_invoices SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE purchase_payments SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE purchase_returns SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE purchase_return_items SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE customer_categories SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE sales_categories SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE customers SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE sales_quotations SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE sales_quotation_items SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE sales_orders SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE sales_order_items SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE delivery_orders SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE delivery_order_items SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE sales_invoices SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE sales_receipts SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE sales_returns SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE sales_return_items SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE item_variants SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE warehouses SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE stock SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE stock_movements SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE stock_transfers SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE stock_transfer_items SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE stock_adjustments SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE stock_adjustment_items SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE price_changes SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE price_change_items SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE bank_transfers SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE bank_reconciliations SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE bank_reconciliation_items SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE chart_of_accounts SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE journal_entries SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE journal_entry_items SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE fixed_assets SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE period_closings SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE fx_differences SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE daily_revenue SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE marketplace_connections SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE marketplace_orders SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE bank_statement_imports SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE bank_statement_lines SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE efaktur_exports SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE efaktur_lines SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE bill_of_materials SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE bom_items SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE work_orders SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
    UPDATE work_order_items SET tenant_id = default_tenant_id WHERE tenant_id IS NULL;
END $$;

-- ============ Set NOT NULL after backfill ============
ALTER TABLE categories ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE products ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE product_availability ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE shifts ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE orders ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE order_items ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE transactions ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE branches ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE departments ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE employees ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE commissions ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE salaries ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE cash_transactions ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE debts ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE receivables ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE withdrawals ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE suppliers ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE purchase_orders ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE purchase_items ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE purchase_invoices ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE purchase_payments ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE purchase_returns ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE purchase_return_items ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE customer_categories ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE sales_categories ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE customers ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE sales_quotations ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE sales_quotation_items ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE sales_orders ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE sales_order_items ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE delivery_orders ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE delivery_order_items ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE sales_invoices ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE sales_receipts ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE sales_returns ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE sales_return_items ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE item_variants ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE warehouses ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE stock ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE stock_movements ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE stock_transfers ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE stock_transfer_items ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE stock_adjustments ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE stock_adjustment_items ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE price_changes ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE price_change_items ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE bank_transfers ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE bank_reconciliations ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE bank_reconciliation_items ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE chart_of_accounts ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE journal_entries ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE journal_entry_items ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE fixed_assets ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE period_closings ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE fx_differences ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE daily_revenue ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE marketplace_connections ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE marketplace_orders ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE bank_statement_imports ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE bank_statement_lines ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE efaktur_exports ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE efaktur_lines ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE bill_of_materials ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE bom_items ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE work_orders ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE work_order_items ALTER COLUMN tenant_id SET NOT NULL;

-- ============ Performance Indexes ============
CREATE INDEX IF NOT EXISTS idx_categories_tenant ON categories(tenant_id);
CREATE INDEX IF NOT EXISTS idx_products_tenant ON products(tenant_id);
CREATE INDEX IF NOT EXISTS idx_product_availability_tenant ON product_availability(tenant_id);
CREATE INDEX IF NOT EXISTS idx_shifts_tenant ON shifts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_orders_tenant ON orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_order_items_tenant ON order_items(tenant_id);
CREATE INDEX IF NOT EXISTS idx_transactions_tenant ON transactions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_branches_tenant ON branches(tenant_id);
CREATE INDEX IF NOT EXISTS idx_departments_tenant ON departments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_employees_tenant ON employees(tenant_id);
CREATE INDEX IF NOT EXISTS idx_commissions_tenant ON commissions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_salaries_tenant ON salaries(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cash_transactions_tenant ON cash_transactions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_debts_tenant ON debts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_receivables_tenant ON receivables(tenant_id);
CREATE INDEX IF NOT EXISTS idx_withdrawals_tenant ON withdrawals(tenant_id);
CREATE INDEX IF NOT EXISTS idx_suppliers_tenant ON suppliers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_tenant ON purchase_orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_purchase_items_tenant ON purchase_items(tenant_id);
CREATE INDEX IF NOT EXISTS idx_purchase_invoices_tenant ON purchase_invoices(tenant_id);
CREATE INDEX IF NOT EXISTS idx_purchase_payments_tenant ON purchase_payments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_purchase_returns_tenant ON purchase_returns(tenant_id);
CREATE INDEX IF NOT EXISTS idx_customer_categories_tenant ON customer_categories(tenant_id);
CREATE INDEX IF NOT EXISTS idx_sales_categories_tenant ON sales_categories(tenant_id);
CREATE INDEX IF NOT EXISTS idx_customers_tenant ON customers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_sales_quotations_tenant ON sales_quotations(tenant_id);
CREATE INDEX IF NOT EXISTS idx_sales_orders_tenant ON sales_orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_delivery_orders_tenant ON delivery_orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_tenant ON sales_invoices(tenant_id);
CREATE INDEX IF NOT EXISTS idx_sales_receipts_tenant ON sales_receipts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_sales_returns_tenant ON sales_returns(tenant_id);
CREATE INDEX IF NOT EXISTS idx_item_variants_tenant ON item_variants(tenant_id);
CREATE INDEX IF NOT EXISTS idx_warehouses_tenant ON warehouses(tenant_id);
CREATE INDEX IF NOT EXISTS idx_stock_tenant ON stock(tenant_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_tenant ON stock_movements(tenant_id);
CREATE INDEX IF NOT EXISTS idx_stock_transfers_tenant ON stock_transfers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_tenant ON stock_adjustments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_price_changes_tenant ON price_changes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_bank_transfers_tenant ON bank_transfers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_bank_reconciliations_tenant ON bank_reconciliations(tenant_id);
CREATE INDEX IF NOT EXISTS idx_chart_of_accounts_tenant ON chart_of_accounts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_journal_entries_tenant ON journal_entries(tenant_id);
CREATE INDEX IF NOT EXISTS idx_fixed_assets_tenant ON fixed_assets(tenant_id);
CREATE INDEX IF NOT EXISTS idx_period_closings_tenant ON period_closings(tenant_id);
CREATE INDEX IF NOT EXISTS idx_daily_revenue_tenant ON daily_revenue(tenant_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_connections_tenant ON marketplace_connections(tenant_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_tenant ON marketplace_orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_bank_statement_imports_tenant ON bank_statement_imports(tenant_id);
CREATE INDEX IF NOT EXISTS idx_efaktur_exports_tenant ON efaktur_exports(tenant_id);
CREATE INDEX IF NOT EXISTS idx_bill_of_materials_tenant ON bill_of_materials(tenant_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_tenant ON work_orders(tenant_id);

-- ============ RLS (Row Level Security) ============
-- Enable RLS on all tenant-specific tables
ALTER TABLE categories ENABLE ROW LEVEL SECURITY;
ALTER TABLE products ENABLE ROW LEVEL SECURITY;
ALTER TABLE product_availability ENABLE ROW LEVEL SECURITY;
ALTER TABLE shifts ENABLE ROW LEVEL SECURITY;
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE order_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE transactions ENABLE ROW LEVEL SECURITY;
ALTER TABLE branches ENABLE ROW LEVEL SECURITY;
ALTER TABLE departments ENABLE ROW LEVEL SECURITY;
ALTER TABLE employees ENABLE ROW LEVEL SECURITY;
ALTER TABLE commissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE salaries ENABLE ROW LEVEL SECURITY;
ALTER TABLE cash_transactions ENABLE ROW LEVEL SECURITY;
ALTER TABLE debts ENABLE ROW LEVEL SECURITY;
ALTER TABLE receivables ENABLE ROW LEVEL SECURITY;
ALTER TABLE withdrawals ENABLE ROW LEVEL SECURITY;
ALTER TABLE suppliers ENABLE ROW LEVEL SECURITY;
ALTER TABLE purchase_orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE purchase_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE purchase_invoices ENABLE ROW LEVEL SECURITY;
ALTER TABLE purchase_payments ENABLE ROW LEVEL SECURITY;
ALTER TABLE purchase_returns ENABLE ROW LEVEL SECURITY;
ALTER TABLE purchase_return_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE customer_categories ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales_categories ENABLE ROW LEVEL SECURITY;
ALTER TABLE customers ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales_quotations ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales_quotation_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales_orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales_order_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE delivery_orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE delivery_order_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales_invoices ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales_receipts ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales_returns ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales_return_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE item_variants ENABLE ROW LEVEL SECURITY;
ALTER TABLE warehouses ENABLE ROW LEVEL SECURITY;
ALTER TABLE stock ENABLE ROW LEVEL SECURITY;
ALTER TABLE stock_movements ENABLE ROW LEVEL SECURITY;
ALTER TABLE stock_transfers ENABLE ROW LEVEL SECURITY;
ALTER TABLE stock_transfer_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE stock_adjustments ENABLE ROW LEVEL SECURITY;
ALTER TABLE stock_adjustment_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE price_changes ENABLE ROW LEVEL SECURITY;
ALTER TABLE price_change_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE bank_transfers ENABLE ROW LEVEL SECURITY;
ALTER TABLE bank_reconciliations ENABLE ROW LEVEL SECURITY;
ALTER TABLE bank_reconciliation_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE chart_of_accounts ENABLE ROW LEVEL SECURITY;
ALTER TABLE journal_entries ENABLE ROW LEVEL SECURITY;
ALTER TABLE journal_entry_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE fixed_assets ENABLE ROW LEVEL SECURITY;
ALTER TABLE period_closings ENABLE ROW LEVEL SECURITY;
ALTER TABLE fx_differences ENABLE ROW LEVEL SECURITY;
ALTER TABLE daily_revenue ENABLE ROW LEVEL SECURITY;
ALTER TABLE marketplace_connections ENABLE ROW LEVEL SECURITY;
ALTER TABLE marketplace_orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE bank_statement_imports ENABLE ROW LEVEL SECURITY;
ALTER TABLE bank_statement_lines ENABLE ROW LEVEL SECURITY;
ALTER TABLE efaktur_exports ENABLE ROW LEVEL SECURITY;
ALTER TABLE efaktur_lines ENABLE ROW LEVEL SECURITY;
ALTER TABLE bill_of_materials ENABLE ROW LEVEL SECURITY;
ALTER TABLE bom_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE work_orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE work_order_items ENABLE ROW LEVEL SECURITY;

-- ============ FORCE RLS for table owners (even superuser bypasses) ============
ALTER TABLE categories FORCE ROW LEVEL SECURITY;
ALTER TABLE products FORCE ROW LEVEL SECURITY;
ALTER TABLE product_availability FORCE ROW LEVEL SECURITY;
ALTER TABLE shifts FORCE ROW LEVEL SECURITY;
ALTER TABLE orders FORCE ROW LEVEL SECURITY;
ALTER TABLE order_items FORCE ROW LEVEL SECURITY;
ALTER TABLE transactions FORCE ROW LEVEL SECURITY;
ALTER TABLE branches FORCE ROW LEVEL SECURITY;
ALTER TABLE departments FORCE ROW LEVEL SECURITY;
ALTER TABLE employees FORCE ROW LEVEL SECURITY;
ALTER TABLE commissions FORCE ROW LEVEL SECURITY;
ALTER TABLE salaries FORCE ROW LEVEL SECURITY;
ALTER TABLE cash_transactions FORCE ROW LEVEL SECURITY;
ALTER TABLE debts FORCE ROW LEVEL SECURITY;
ALTER TABLE receivables FORCE ROW LEVEL SECURITY;
ALTER TABLE withdrawals FORCE ROW LEVEL SECURITY;
ALTER TABLE suppliers FORCE ROW LEVEL SECURITY;
ALTER TABLE purchase_orders FORCE ROW LEVEL SECURITY;
ALTER TABLE purchase_items FORCE ROW LEVEL SECURITY;
ALTER TABLE purchase_invoices FORCE ROW LEVEL SECURITY;
ALTER TABLE purchase_payments FORCE ROW LEVEL SECURITY;
ALTER TABLE purchase_returns FORCE ROW LEVEL SECURITY;
ALTER TABLE purchase_return_items FORCE ROW LEVEL SECURITY;
ALTER TABLE customer_categories FORCE ROW LEVEL SECURITY;
ALTER TABLE sales_categories FORCE ROW LEVEL SECURITY;
ALTER TABLE customers FORCE ROW LEVEL SECURITY;
ALTER TABLE sales_quotations FORCE ROW LEVEL SECURITY;
ALTER TABLE sales_quotation_items FORCE ROW LEVEL SECURITY;
ALTER TABLE sales_orders FORCE ROW LEVEL SECURITY;
ALTER TABLE sales_order_items FORCE ROW LEVEL SECURITY;
ALTER TABLE delivery_orders FORCE ROW LEVEL SECURITY;
ALTER TABLE delivery_order_items FORCE ROW LEVEL SECURITY;
ALTER TABLE sales_invoices FORCE ROW LEVEL SECURITY;
ALTER TABLE sales_receipts FORCE ROW LEVEL SECURITY;
ALTER TABLE sales_returns FORCE ROW LEVEL SECURITY;
ALTER TABLE sales_return_items FORCE ROW LEVEL SECURITY;
ALTER TABLE item_variants FORCE ROW LEVEL SECURITY;
ALTER TABLE warehouses FORCE ROW LEVEL SECURITY;
ALTER TABLE stock FORCE ROW LEVEL SECURITY;
ALTER TABLE stock_movements FORCE ROW LEVEL SECURITY;
ALTER TABLE stock_transfers FORCE ROW LEVEL SECURITY;
ALTER TABLE stock_transfer_items FORCE ROW LEVEL SECURITY;
ALTER TABLE stock_adjustments FORCE ROW LEVEL SECURITY;
ALTER TABLE stock_adjustment_items FORCE ROW LEVEL SECURITY;
ALTER TABLE price_changes FORCE ROW LEVEL SECURITY;
ALTER TABLE price_change_items FORCE ROW LEVEL SECURITY;
ALTER TABLE bank_transfers FORCE ROW LEVEL SECURITY;
ALTER TABLE bank_reconciliations FORCE ROW LEVEL SECURITY;
ALTER TABLE bank_reconciliation_items FORCE ROW LEVEL SECURITY;
ALTER TABLE chart_of_accounts FORCE ROW LEVEL SECURITY;
ALTER TABLE journal_entries FORCE ROW LEVEL SECURITY;
ALTER TABLE journal_entry_items FORCE ROW LEVEL SECURITY;
ALTER TABLE fixed_assets FORCE ROW LEVEL SECURITY;
ALTER TABLE period_closings FORCE ROW LEVEL SECURITY;
ALTER TABLE fx_differences FORCE ROW LEVEL SECURITY;
ALTER TABLE daily_revenue FORCE ROW LEVEL SECURITY;
ALTER TABLE marketplace_connections FORCE ROW LEVEL SECURITY;
ALTER TABLE marketplace_orders FORCE ROW LEVEL SECURITY;
ALTER TABLE bank_statement_imports FORCE ROW LEVEL SECURITY;
ALTER TABLE bank_statement_lines FORCE ROW LEVEL SECURITY;
ALTER TABLE efaktur_exports FORCE ROW LEVEL SECURITY;
ALTER TABLE efaktur_lines FORCE ROW LEVEL SECURITY;
ALTER TABLE bill_of_materials FORCE ROW LEVEL SECURITY;
ALTER TABLE bom_items FORCE ROW LEVEL SECURITY;
ALTER TABLE work_orders FORCE ROW LEVEL SECURITY;
ALTER TABLE work_order_items FORCE ROW LEVEL SECURITY;

-- ============ RLS Policies ============
-- Policy: Users can only see data belonging to their tenant
-- Uses session variable: SET app.current_tenant = '<tenant_id>'

-- Categories
CREATE POLICY tenant_isolation ON categories FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Products
CREATE POLICY tenant_isolation ON products FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Product Availability
CREATE POLICY tenant_isolation ON product_availability FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Shifts
CREATE POLICY tenant_isolation ON shifts FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Orders
CREATE POLICY tenant_isolation ON orders FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Order Items
CREATE POLICY tenant_isolation ON order_items FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Transactions
CREATE POLICY tenant_isolation ON transactions FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Branches
CREATE POLICY tenant_isolation ON branches FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Departments
CREATE POLICY tenant_isolation ON departments FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Employees
CREATE POLICY tenant_isolation ON employees FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Commissions
CREATE POLICY tenant_isolation ON commissions FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Salaries
CREATE POLICY tenant_isolation ON salaries FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Cash Transactions
CREATE POLICY tenant_isolation ON cash_transactions FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Debts
CREATE POLICY tenant_isolation ON debts FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Receivables
CREATE POLICY tenant_isolation ON receivables FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Withdrawals
CREATE POLICY tenant_isolation ON withdrawals FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Suppliers
CREATE POLICY tenant_isolation ON suppliers FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Purchase Orders
CREATE POLICY tenant_isolation ON purchase_orders FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Purchase Items
CREATE POLICY tenant_isolation ON purchase_items FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Purchase Invoices
CREATE POLICY tenant_isolation ON purchase_invoices FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Purchase Payments
CREATE POLICY tenant_isolation ON purchase_payments FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Purchase Returns
CREATE POLICY tenant_isolation ON purchase_returns FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Purchase Return Items
CREATE POLICY tenant_isolation ON purchase_return_items FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Customer Categories
CREATE POLICY tenant_isolation ON customer_categories FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Sales Categories
CREATE POLICY tenant_isolation ON sales_categories FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Customers
CREATE POLICY tenant_isolation ON customers FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Sales Quotations
CREATE POLICY tenant_isolation ON sales_quotations FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Sales Quotation Items
CREATE POLICY tenant_isolation ON sales_quotation_items FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Sales Orders
CREATE POLICY tenant_isolation ON sales_orders FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Sales Order Items
CREATE POLICY tenant_isolation ON sales_order_items FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Delivery Orders
CREATE POLICY tenant_isolation ON delivery_orders FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Delivery Order Items
CREATE POLICY tenant_isolation ON delivery_order_items FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Sales Invoices
CREATE POLICY tenant_isolation ON sales_invoices FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Sales Receipts
CREATE POLICY tenant_isolation ON sales_receipts FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Sales Returns
CREATE POLICY tenant_isolation ON sales_returns FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Sales Return Items
CREATE POLICY tenant_isolation ON sales_return_items FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Item Variants
CREATE POLICY tenant_isolation ON item_variants FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Warehouses
CREATE POLICY tenant_isolation ON warehouses FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Stock
CREATE POLICY tenant_isolation ON stock FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Stock Movements
CREATE POLICY tenant_isolation ON stock_movements FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Stock Transfers
CREATE POLICY tenant_isolation ON stock_transfers FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Stock Transfer Items
CREATE POLICY tenant_isolation ON stock_transfer_items FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Stock Adjustments
CREATE POLICY tenant_isolation ON stock_adjustments FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Stock Adjustment Items
CREATE POLICY tenant_isolation ON stock_adjustment_items FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Price Changes
CREATE POLICY tenant_isolation ON price_changes FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Price Change Items
CREATE POLICY tenant_isolation ON price_change_items FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Bank Transfers
CREATE POLICY tenant_isolation ON bank_transfers FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Bank Reconciliations
CREATE POLICY tenant_isolation ON bank_reconciliations FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Bank Reconciliation Items
CREATE POLICY tenant_isolation ON bank_reconciliation_items FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Chart of Accounts
CREATE POLICY tenant_isolation ON chart_of_accounts FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Journal Entries
CREATE POLICY tenant_isolation ON journal_entries FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Journal Entry Items
CREATE POLICY tenant_isolation ON journal_entry_items FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Fixed Assets
CREATE POLICY tenant_isolation ON fixed_assets FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Period Closings
CREATE POLICY tenant_isolation ON period_closings FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- FX Differences
CREATE POLICY tenant_isolation ON fx_differences FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Daily Revenue
CREATE POLICY tenant_isolation ON daily_revenue FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Marketplace Connections
CREATE POLICY tenant_isolation ON marketplace_connections FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Marketplace Orders
CREATE POLICY tenant_isolation ON marketplace_orders FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Bank Statement Imports
CREATE POLICY tenant_isolation ON bank_statement_imports FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Bank Statement Lines
CREATE POLICY tenant_isolation ON bank_statement_lines FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- EFaktur Exports
CREATE POLICY tenant_isolation ON efaktur_exports FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- EFaktur Lines
CREATE POLICY tenant_isolation ON efaktur_lines FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Bill of Materials
CREATE POLICY tenant_isolation ON bill_of_materials FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- BOM Items
CREATE POLICY tenant_isolation ON bom_items FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Work Orders
CREATE POLICY tenant_isolation ON work_orders FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- Work Order Items
CREATE POLICY tenant_isolation ON work_order_items FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::bigint);

-- ============ Tenant Management Permissions ============
INSERT INTO permissions (name, description) VALUES
    ('tenant_manage', 'Kelola Tenant (Super Admin)'),
    ('tenant_view', 'Lihat Data Tenant')
ON CONFLICT (name) DO NOTHING;
