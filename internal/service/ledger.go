package service

import (
	"accounting/internal/domain"
	"accounting/internal/repo"
	"fmt"
)

type LedgerService struct {
	store *repo.Store
}

func NewLedgerService(store *repo.Store) *LedgerService {
	return &LedgerService{store: store}
}

func (ls *LedgerService) PostEntry(entry *domain.JournalEntry) error {
	if err := entry.Validate(); err != nil {
		return err
	}

	for _, line := range entry.Lines {
		account, err := ls.store.GetAccount(line.AccountID)
		if err != nil {
			return err
		}

		if line.Debit > 0 {
			account.Increase(line.Debit)
		} else if line.Credit > 0 {
			account.Increase(line.Credit)
		}
	}

	return ls.store.SaveEntry(entry)
}

func (ls *LedgerService) GetTrialBalance() map[string]interface{} {
	accounts := ls.store.GetAllAccounts()

	var totalDebits, totalCredits float64
	rows := []map[string]interface{}{}

	for _, acc := range accounts {
		var balance float64
		if acc.NormalBalance == "debit" {
			balance = acc.Balance
			totalDebits += balance
		} else {
			balance = acc.Balance
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
	}
}
