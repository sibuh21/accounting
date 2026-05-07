package service

import (
	"accounting/internal/domain"
	"accounting/internal/repo"
	"fmt"
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

func (cs *ClosingService) CloseMonth(monthEnd time.Time, ownerDrawAmount float64) (*domain.JournalEntry, error) {
	accounts := cs.store.GetAllAccounts()

	var totalRevenue, totalExpenses float64
	var revenueAccounts, expenseAccounts []*domain.Account

	for _, acc := range accounts {
		switch acc.Type {
		case domain.Revenue:
			revenueAccounts = append(revenueAccounts, acc)
			totalRevenue += acc.Balance
		case domain.Expense:
			expenseAccounts = append(expenseAccounts, acc)
			totalExpenses += acc.Balance
		}
	}

	netProfit := totalRevenue - totalExpenses
	closingEntry := domain.NewJournalEntry(monthEnd, "Monthly closing entry")

	for _, acc := range revenueAccounts {
		if acc.Balance > 0 {
			closingEntry.Lines = append(closingEntry.Lines, domain.EntryLine{
				AccountID:   acc.ID,
				AccountName: acc.Name,
				Debit:       acc.Balance,
			})
		}
	}

	for _, acc := range expenseAccounts {
		if acc.Balance > 0 {
			closingEntry.Lines = append(closingEntry.Lines, domain.EntryLine{
				AccountID:   acc.ID,
				AccountName: acc.Name,
				Credit:      acc.Balance,
			})
		}
	}

	for _, acc := range accounts {
		if acc.Name == "Owner's Draw" && acc.Balance > 0 {
			closingEntry.Lines = append(closingEntry.Lines, domain.EntryLine{
				AccountID:   acc.ID,
				AccountName: acc.Name,
				Credit:      acc.Balance,
			})
			netProfit -= acc.Balance
		}
	}

	var retainedEarningsAccount *domain.Account
	for _, acc := range accounts {
		if acc.Name == "Retained Earnings" {
			retainedEarningsAccount = acc
			break
		}
	}

	if retainedEarningsAccount == nil {
		return nil, fmt.Errorf("retained earnings account not found")
	}

	if netProfit > 0 {
		closingEntry.Lines = append(closingEntry.Lines, domain.EntryLine{
			AccountID:   retainedEarningsAccount.ID,
			AccountName: retainedEarningsAccount.Name,
			Credit:      netProfit,
		})
	} else if netProfit < 0 {
		closingEntry.Lines = append(closingEntry.Lines, domain.EntryLine{
			AccountID:   retainedEarningsAccount.ID,
			AccountName: retainedEarningsAccount.Name,
			Debit:       -netProfit,
		})
	}

	if err := cs.ledgerService.PostEntry(closingEntry); err != nil {
		return nil, err
	}

	return closingEntry, nil
}
