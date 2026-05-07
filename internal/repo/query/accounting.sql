-- name: CreateAccount :one
INSERT INTO accounts (code, name, type, normal_balance)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts WHERE id = $1;

-- name: GetAccountByCode :one
SELECT * FROM accounts WHERE code = $1;

-- name: ListAccounts :many
SELECT * FROM accounts ORDER BY code;

-- name: UpdateAccountBalance :exec
UPDATE accounts 
SET balance = balance + $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: CreateJournalEntry :one
INSERT INTO journal_entries (date, description)
VALUES ($1, $2)
RETURNING *;

-- name: CreateEntryLine :one
INSERT INTO entry_lines (journal_entry_id, account_id, debit, credit)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListJournalEntries :many
SELECT je.*, array_to_json(array_agg(el.*)) as lines
FROM journal_entries je
LEFT JOIN entry_lines el ON je.id = el.journal_entry_id
GROUP BY je.id
ORDER BY je.date DESC;

-- name: GetTrialBalance :many
SELECT code, name, type, normal_balance, balance
FROM accounts
ORDER BY code;
