CREATE TABLE IF NOT EXISTS cash_transactions (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(20) NOT NULL DEFAULT 'in',
    amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    description TEXT NOT NULL DEFAULT '',
    category VARCHAR(100) NOT NULL DEFAULT '',
    reference VARCHAR(100) NOT NULL DEFAULT '',
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS debts (
    id BIGSERIAL PRIMARY KEY,
    creditor_name VARCHAR(255) NOT NULL,
    amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    paid_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    description TEXT NOT NULL DEFAULT '',
    due_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS receivables (
    id BIGSERIAL PRIMARY KEY,
    debtor_name VARCHAR(255) NOT NULL,
    amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    paid_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    description TEXT NOT NULL DEFAULT '',
    due_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS withdrawals (
    id BIGSERIAL PRIMARY KEY,
    amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    description TEXT NOT NULL DEFAULT '',
    method VARCHAR(50) NOT NULL DEFAULT 'cash',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    requested_by BIGINT REFERENCES users(id),
    approved_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
