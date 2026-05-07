package service

import (
	"accounting/internal/domain"
	"accounting/internal/repo"
	"accounting/internal/repo/db"
	"context"
	"fmt"
	"strings"
	"time"
)

type ClosingService struct {
	store         *repo.Store
	ledgerService *LedgerService
}

func NewClosingService(store *repo.Store, ledgerService *LedgerService) *ClosingService {
	return &ClosingService{
		store:         store,
		ledgerService: ledgerService,
	}
}

func (cs *ClosingService) CloseMonth(ctx context.Context, monthEnd time.Time) ([]*domain.JournalEntry, error) {
	accounts, err := cs.store.Queries().ListAccounts(ctx)
	if err != nil {
		return nil, err
	}

	var totalRevenue, totalExpenses float64
	var revenueAccounts, expenseAccounts []db.Account

	for _, acc := range accounts {
		balance := fromCents(acc.Balance)
		switch acc.Type {
		case string(domain.Revenue):
			revenueAccounts = append(revenueAccounts, acc)
			totalRevenue += balance
		case string(domain.Expense):
			expenseAccounts = append(expenseAccounts, acc)
			totalExpenses += balance
		}
	}

	netProfit := totalRevenue - totalExpenses
	if netProfit == 0 && len(revenueAccounts) == 0 && len(expenseAccounts) == 0 {
		return nil, fmt.Errorf("no temporary accounts to close")
	}

	closingEntry := domain.NewJournalEntry(monthEnd, "Monthly closing entry - Revenue & Expenses")

	for _, acc := range revenueAccounts {
		balance := fromCents(acc.Balance)
		if balance > 0 {
			closingEntry.Lines = append(closingEntry.Lines, domain.EntryLine{
				AccountID:   acc.ID.String(),
				AccountName: acc.Name,
				Debit:       balance,
			})
		}
	}

	for _, acc := range expenseAccounts {
		balance := fromCents(acc.Balance)
		if balance > 0 {
			closingEntry.Lines = append(closingEntry.Lines, domain.EntryLine{
				AccountID:   acc.ID.String(),
				AccountName: acc.Name,
				Credit:      balance,
			})
		}
	}

	var retainedEarningsAccount *db.Account
	for i, acc := range accounts {
		if strings.Contains(strings.ToLower(acc.Name), "retained earnings") {
			retainedEarningsAccount = &accounts[i]
			break
		}
	}

	if retainedEarningsAccount == nil {
		return nil, fmt.Errorf("retained earnings account not found")
	}

	if netProfit > 0 {
		closingEntry.Lines = append(closingEntry.Lines, domain.EntryLine{
			AccountID:   retainedEarningsAccount.ID.String(),
			AccountName: retainedEarningsAccount.Name,
			Credit:      netProfit,
		})
	} else if netProfit < 0 {
		closingEntry.Lines = append(closingEntry.Lines, domain.EntryLine{
			AccountID:   retainedEarningsAccount.ID.String(),
			AccountName: retainedEarningsAccount.Name,
			Debit:       -netProfit,
		})
	}

	if err := cs.ledgerService.PostEntry(ctx, closingEntry); err != nil {
		return nil, err
	}

	return []*domain.JournalEntry{closingEntry}, nil
}
