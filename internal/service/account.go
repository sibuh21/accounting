package service

import (
	"accounting/internal/domain"
	"accounting/internal/repo"
	"accounting/internal/repo/db"
	"context"
)

type AccountService struct {
	store *repo.Store
}

func NewAccountService(store *repo.Store) *AccountService {
	return &AccountService{store: store}
}

func (s *AccountService) CreateAccount(ctx context.Context, req domain.CreateAccountRequest) (*db.Account, error) {
	// Determine normal balance
	normalBalance := "debit"
	if req.Type == domain.Liability || req.Type == domain.Equity || req.Type == domain.Revenue {
		normalBalance = "credit"
	}

	acc, err := s.store.Queries().CreateAccount(ctx, db.CreateAccountParams{
		Code:          req.Code,
		Name:          req.Name,
		Type:          string(req.Type),
		NormalBalance: normalBalance,
	})
	if err != nil {
		return nil, err
	}
	
	return &acc, nil
}

func (s *AccountService) ListAccounts(ctx context.Context) ([]db.Account, error) {
	return s.store.Queries().ListAccounts(ctx)
}

func (s *AccountService) GetEntries(ctx context.Context) ([]db.ListJournalEntriesRow, error) {
	return s.store.Queries().ListJournalEntries(ctx)
}
