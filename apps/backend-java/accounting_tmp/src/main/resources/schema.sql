-- Arda Core Accounting Schema
-- Version: 1.0 (Banking Standard)

CREATE TABLE IF NOT EXISTS accounts (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_by VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS balances (
    account_id VARCHAR(36) PRIMARY KEY REFERENCES accounts(id),
    available_balance DECIMAL(20, 2) NOT NULL DEFAULT 0.00,
    locked_balance DECIMAL(20, 2) NOT NULL DEFAULT 0.00,
    version BIGINT NOT NULL DEFAULT 0, -- For Optimistic Locking
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS journals (
    id VARCHAR(36) PRIMARY KEY,
    description TEXT,
    reference_id VARCHAR(100), -- ID từ EPAS hoặc External System
    type VARCHAR(50) NOT NULL, -- TRANSFER, DEPOSIT, WITHDRAW, FEE
    status VARCHAR(20) NOT NULL DEFAULT 'POSTED',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS entries (
    id BIGSERIAL PRIMARY KEY,
    journal_id VARCHAR(36) NOT NULL REFERENCES journals(id),
    account_id VARCHAR(36) NOT NULL REFERENCES accounts(id),
    type VARCHAR(10) NOT NULL, -- DEBIT, CREDIT
    amount DECIMAL(20, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexing for performance
CREATE INDEX idx_entries_account_id ON entries(account_id);
CREATE INDEX idx_entries_journal_id ON entries(journal_id);
CREATE INDEX idx_journals_reference_id ON journals(reference_id);
