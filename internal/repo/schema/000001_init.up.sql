CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    normal_balance TEXT NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0, -- Store in cents
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE journal_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE entry_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    journal_entry_id UUID NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id),
    debit BIGINT NOT NULL DEFAULT 0, -- Store in cents
    credit BIGINT NOT NULL DEFAULT 0 -- Store in cents
);

CREATE INDEX idx_accounts_code ON accounts(code);
CREATE INDEX idx_entry_lines_journal_entry_id ON entry_lines(journal_entry_id);
CREATE INDEX idx_entry_lines_account_id ON entry_lines(account_id);
