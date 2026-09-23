-- --------------------------------------------------------
-- Migration 010: Add General Ledger tables
-- --------------------------------------------------------

-- Chart of Accounts (Akun Perkiraan)
CREATE TABLE IF NOT EXISTS chart_of_accounts (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(20) NOT NULL,
    normal_balance VARCHAR(10) NOT NULL DEFAULT 'debit',
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Journal Entries (Jurnal Umum)
CREATE TABLE IF NOT EXISTS journal_entries (
    id SERIAL PRIMARY KEY,
    entry_number VARCHAR(20) NOT NULL,
    entry_date DATE DEFAULT CURRENT_DATE,
    description TEXT,
    reference VARCHAR(100) DEFAULT 'Umum',
    status VARCHAR(20) DEFAULT 'posted',
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Journal Entry Items
CREATE TABLE IF NOT EXISTS journal_entry_items (
    id SERIAL PRIMARY KEY,
    journal_entry_id INTEGER REFERENCES journal_entries(id) ON DELETE CASCADE,
    account_id INTEGER REFERENCES chart_of_accounts(id),
    debit NUMERIC(15,2) NOT NULL DEFAULT 0,
    credit NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Fixed Assets (Aset Tetap untuk perhitungan depresiasi)
CREATE TABLE IF NOT EXISTS fixed_assets (
    id SERIAL PRIMARY KEY,
    asset_code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    purchase_date DATE DEFAULT CURRENT_DATE,
    cost NUMERIC(15,2) NOT NULL DEFAULT 0,
    salvage_value NUMERIC(15,2) NOT NULL DEFAULT 0,
    useful_life_months INTEGER NOT NULL DEFAULT 12,
    accumulated_depreciation NUMERIC(15,2) NOT NULL DEFAULT 0,
    depreciation_account_id INTEGER REFERENCES chart_of_accounts(id),
    accumulated_account_id INTEGER REFERENCES chart_of_accounts(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Period Closings (Proses Akhir Bulan)
CREATE TABLE IF NOT EXISTS period_closings (
    id SERIAL PRIMARY KEY,
    period VARCHAR(7) NOT NULL,
    status VARCHAR(20) DEFAULT 'open',
    depreciation_total NUMERIC(15,2) NOT NULL DEFAULT 0,
    fx_gain NUMERIC(15,2) NOT NULL DEFAULT 0,
    fx_loss NUMERIC(15,2) NOT NULL DEFAULT 0,
    notes TEXT,
    closed_by BIGINT,
    closed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- FX Differences (Selisih Kurs Valas)
CREATE TABLE IF NOT EXISTS fx_differences (
    id SERIAL PRIMARY KEY,
    period VARCHAR(7) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    description TEXT,
    gain_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    loss_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE UNIQUE INDEX IF NOT EXISTS idx_chart_of_accounts_code ON chart_of_accounts(code);
CREATE INDEX IF NOT EXISTS idx_chart_of_accounts_category ON chart_of_accounts(category);
CREATE UNIQUE INDEX IF NOT EXISTS idx_journal_entries_number ON journal_entries(entry_number);
CREATE INDEX IF NOT EXISTS idx_journal_entries_date ON journal_entries(entry_date);
CREATE INDEX IF NOT EXISTS idx_journal_entry_items_entry ON journal_entry_items(journal_entry_id);
CREATE INDEX IF NOT EXISTS idx_journal_entry_items_account ON journal_entry_items(account_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_period_closings_period ON period_closings(period);
CREATE INDEX IF NOT EXISTS idx_fx_differences_period ON fx_differences(period);
