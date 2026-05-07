package service

import (
	"accounting/internal/domain"
	"accounting/internal/repo"
	"accounting/internal/repo/db"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/google/uuid"
)

type LedgerService struct {
	store *repo.Store
}

func NewLedgerService(store *repo.Store) *LedgerService {
	return &LedgerService{store: store}
}

func (ls *LedgerService) PostEntry(ctx context.Context, entry *domain.JournalEntry) error {
	if err := entry.Validate(); err != nil {
		return err
	}

	return ls.store.ExecTx(ctx, func(q *db.Queries) error {
		je, err := q.CreateJournalEntry(ctx, db.CreateJournalEntryParams{
			Date:        pgtype.Timestamptz{Time: entry.Date, Valid: true},
			Description: entry.Description,
		})
		if err != nil {
			return err
		}

		for _, line := range entry.Lines {
			accID, err := uuid.Parse(line.AccountID)
			if err != nil {
				acc, lookupErr := q.GetAccountByCode(ctx, line.AccountName)
				if lookupErr != nil {
					return fmt.Errorf("account not found: %s", line.AccountName)
				}
				accID = acc.ID
			}

			_, err = q.CreateEntryLine(ctx, db.CreateEntryLineParams{
				JournalEntryID: je.ID,
				AccountID:      accID,
				Debit:          toCents(line.Debit),
				Credit:         toCents(line.Credit),
			})
			if err != nil {
				return err
			}

			acc, err := q.GetAccount(ctx, accID)
			if err != nil {
				return err
			}

			delta := 0.0
			if acc.NormalBalance == "debit" {
				delta = line.Debit - line.Credit
			} else {
				delta = line.Credit - line.Debit
			}

			err = q.UpdateAccountBalance(ctx, db.UpdateAccountBalanceParams{
				ID:      accID,
				Balance: toCents(delta),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (ls *LedgerService) GetTrialBalance(ctx context.Context) (map[string]interface{}, error) {
	accounts, err := ls.store.Queries().GetTrialBalance(ctx)
	if err != nil {
		return nil, err
	}

	var totalDebits, totalCredits float64
	rows := []map[string]interface{}{}

	for _, acc := range accounts {
		balance := fromCents(acc.Balance)
		if acc.NormalBalance == "debit" {
			totalDebits += balance
		} else {
			totalCredits += balance
		}

		rows = append(rows, map[string]interface{}{
			"code":    acc.Code,
			"name":    acc.Name,
			"type":    acc.Type,
			"balance": balance,
		})
	}

	return map[string]interface{}{
		"accounts":      rows,
		"total_debits":  totalDebits,
		"total_credits": totalCredits,
		"is_balanced":   fmt.Sprintf("%.2f", totalDebits) == fmt.Sprintf("%.2f", totalCredits),
	}, nil
}
