package service

import (
	"accounting/internal/domain"
	"accounting/internal/repo"
)

type AccountService struct {
	store *repo.Store
}

func NewAccountService(store *repo.Store) *AccountService {
	return &AccountService{store: store}
}

func (s *AccountService) CreateAccount(req domain.CreateAccountRequest) (*domain.Account, error) {
	acc := domain.NewAccount(req.Code, req.Name, req.Type)
	err := s.store.SaveAccount(acc)
	return acc, err
}

func (s *AccountService) ListAccounts() ([]*domain.Account, error) {
	return s.store.GetAllAccounts(), nil
}

func (s *AccountService) GetEntries() ([]*domain.JournalEntry, error) {
	return s.store.GetEntries(), nil
}
