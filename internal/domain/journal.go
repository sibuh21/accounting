package domain

import (
	"fmt"
	"time"
)

type JournalEntry struct {
	ID          string      `json:"id"`
	Date        time.Time   `json:"date"`
	Description string      `json:"description"`
	Lines       []EntryLine `json:"lines"`
	CreatedAt   time.Time   `json:"created_at"`
}

type EntryLine struct {
	AccountID   string  `json:"account_id"`
	AccountName string  `json:"account_name"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
}

func (je *JournalEntry) Validate() error {
	var totalDebits, totalCredits float64

	for _, line := range je.Lines {
		totalDebits += line.Debit
		totalCredits += line.Credit
	}

	// Use a small epsilon for float comparison
	if fmt.Sprintf("%.2f", totalDebits) != fmt.Sprintf("%.2f", totalCredits) {
		return fmt.Errorf("entry not balanced: debits=%.2f, credits=%.2f", totalDebits, totalCredits)
	}

	return nil
}

func NewJournalEntry(date time.Time, description string) *JournalEntry {
	return &JournalEntry{
		ID:          fmt.Sprintf("je_%d", time.Now().UnixNano()),
		Date:        date,
		Description: description,
		Lines:       []EntryLine{},
		CreatedAt:   time.Now(),
	}
}

type CreateJournalEntryRequest struct {
	Date        time.Time          `json:"date"`
	Description string             `json:"description"`
	Lines       []CreateEntryLine  `json:"lines"`
}

type CreateEntryLine struct {
	AccountID   string  `json:"account_id"`
	AccountName string  `json:"account_name"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
}
