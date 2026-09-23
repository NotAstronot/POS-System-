-- --------------------------------------------------------
-- Migration 009: Add Bank Transfer and Bank Reconciliation tables
-- --------------------------------------------------------

-- Create bank_transfers table
CREATE TABLE IF NOT EXISTS bank_transfers (
    id SERIAL PRIMARY KEY,
    transfer_number VARCHAR(20) NOT NULL,
    from_account_name VARCHAR(100) NOT NULL,
    from_account_type VARCHAR(20) NOT NULL,
    to_account_name VARCHAR(100) NOT NULL,
    to_account_type VARCHAR(20) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    transfer_date DATE DEFAULT CURRENT_DATE,
    notes TEXT,
    status VARCHAR(20) DEFAULT 'completed',
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create bank_reconciliations table
CREATE TABLE IF NOT EXISTS bank_reconciliations (
    id SERIAL PRIMARY KEY,
    reconciliation_number VARCHAR(20) NOT NULL,
    account_name VARCHAR(100) NOT NULL,
    account_type VARCHAR(20) NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    book_balance DECIMAL(15,2) NOT NULL DEFAULT 0,
    statement_balance DECIMAL(15,2) NOT NULL DEFAULT 0,
    difference DECIMAL(15,2) NOT NULL DEFAULT 0,
    notes TEXT,
    status VARCHAR(20) DEFAULT 'completed',
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create bank_reconciliation_items table
CREATE TABLE IF NOT EXISTS bank_reconciliation_items (
    id SERIAL PRIMARY KEY,
    reconciliation_id INTEGER REFERENCES bank_reconciliations(id) ON DELETE CASCADE,
    description VARCHAR(255),
    transaction_date DATE,
    amount DECIMAL(15,2) NOT NULL,
    type VARCHAR(20) NOT NULL,
    is_matched BOOLEAN DEFAULT FALSE
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_bank_transfers_status ON bank_transfers(status);
CREATE INDEX IF NOT EXISTS idx_bank_transfers_from_account ON bank_transfers(from_account_type, from_account_name);
CREATE INDEX IF NOT EXISTS idx_bank_transfers_to_account ON bank_transfers(to_account_type, to_account_name);
CREATE INDEX IF NOT EXISTS idx_bank_reconciliations_status ON bank_reconciliations(status);
CREATE INDEX IF NOT EXISTS idx_bank_reconciliations_period ON bank_reconciliations(period_start, period_end);
CREATE INDEX IF NOT EXISTS idx_bank_reconciliation_items_reconciliation ON bank_reconciliation_items(reconciliation_id);