package domain

import (
	"fmt"
	"time"
)

type AccountType string

const (
	Asset     AccountType = "Asset"
	Liability AccountType = "Liability"
	Equity    AccountType = "Equity"
	Revenue   AccountType = "Revenue"
	Expense   AccountType = "Expense"
)

type Account struct {
	ID            string      `json:"id"`
	Code          string      `json:"code"`
	Name          string      `json:"name"`
	Type          AccountType `json:"type"`
	NormalBalance string      `json:"normal_balance"`
	Balance       float64     `json:"balance"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

func NewAccount(code, name string, accountType AccountType) *Account {
	normalBalance := "debit"
	if accountType == Liability || accountType == Equity || accountType == Revenue {
		normalBalance = "credit"
	}

	return &Account{
		ID:            fmt.Sprintf("acc_%d", time.Now().UnixNano()),
		Code:          code,
		Name:          name,
		Type:          accountType,
		NormalBalance: normalBalance,
		Balance:       0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (a *Account) IsDebit() bool {
	return a.NormalBalance == "debit"
}

func (a *Account) Increase(amount float64) {
	if a.IsDebit() {
		a.Balance += amount
	} else {
		a.Balance -= amount
	}
	a.UpdatedAt = time.Now()
}

func (a *Account) Decrease(amount float64) {
	if a.IsDebit() {
		a.Balance -= amount
	} else {
		a.Balance += amount
	}
	a.UpdatedAt = time.Now()
}

type CreateAccountRequest struct {
	Code string      `json:"code"`
	Name string      `json:"name"`
	Type AccountType `json:"type"`
}
