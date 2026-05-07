package repo

import (
	"accounting/internal/domain"
	"fmt"
	"sync"
)

type Store struct {
	mu       sync.RWMutex
	accounts map[string]*domain.Account
	entries  map[string]*domain.JournalEntry
}

func NewStore() *Store {
	return &Store{
		accounts: make(map[string]*domain.Account),
		entries:  make(map[string]*domain.JournalEntry),
	}
}

func (s *Store) SaveAccount(account *domain.Account) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accounts[account.ID] = account
	return nil
}

func (s *Store) GetAccount(id string) (*domain.Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	account, exists := s.accounts[id]
	if !exists {
		return nil, fmt.Errorf("account not found: %s", id)
	}
	return account, nil
}

func (s *Store) GetAllAccounts() []*domain.Account {
	s.mu.RLock()
	defer s.mu.RUnlock()
	accounts := make([]*domain.Account, 0, len(s.accounts))
	for _, acc := range s.accounts {
		accounts = append(accounts, acc)
	}
	return accounts
}

func (s *Store) SaveEntry(entry *domain.JournalEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[entry.ID] = entry
	return nil
}

func (s *Store) GetEntries() []*domain.JournalEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries := make([]*domain.JournalEntry, 0, len(s.entries))
	for _, entry := range s.entries {
		entries = append(entries, entry)
	}
	return entries
}
