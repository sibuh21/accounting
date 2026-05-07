package service

import (
	"accounting/internal/repo"
	"context"
	"time"
)

type ReporterService struct {
	store *repo.Store
}

func NewReporterService(store *repo.Store) *ReporterService {
	return &ReporterService{store: store}
}

func (rs *ReporterService) GenerateIncomeStatement(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
	accounts, err := rs.store.Queries().ListAccounts(ctx)
	if err != nil {
		return nil, err
	}

	var totalRevenue float64
	var totalExpenses float64
	revenues := []map[string]interface{}{}
	expenses := []map[string]interface{}{}

	for _, acc := range accounts {
		balance := fromCents(acc.Balance)
		if acc.Type == "Revenue" {
			totalRevenue += balance
			revenues = append(revenues, map[string]interface{}{
				"account": acc.Name,
				"amount":  balance,
			})
		} else if acc.Type == "Expense" {
			totalExpenses += balance
			expenses = append(expenses, map[string]interface{}{
				"account": acc.Name,
				"amount":  balance,
			})
		}
	}

	netIncome := totalRevenue - totalExpenses

	return map[string]interface{}{
		"period_start":   startDate,
		"period_end":     endDate,
		"revenues":       revenues,
		"total_revenue":  totalRevenue,
		"expenses":       expenses,
		"total_expenses": totalExpenses,
		"net_income":     netIncome,
	}, nil
}

func (rs *ReporterService) GenerateBalanceSheet(ctx context.Context, asOfDate time.Time) (map[string]interface{}, error) {
	accounts, err := rs.store.Queries().ListAccounts(ctx)
	if err != nil {
		return nil, err
	}

	var totalAssets float64
	var totalLiabilities float64
	var totalEquity float64

	assets := []map[string]interface{}{}
	liabilities := []map[string]interface{}{}
	equity := []map[string]interface{}{}

	for _, acc := range accounts {
		balance := fromCents(acc.Balance)
		switch acc.Type {
		case "Asset":
			totalAssets += balance
			assets = append(assets, map[string]interface{}{
				"account": acc.Name,
				"balance": balance,
			})
		case "Liability":
			totalLiabilities += balance
			liabilities = append(liabilities, map[string]interface{}{
				"account": acc.Name,
				"balance": balance,
			})
		case "Equity":
			totalEquity += balance
			equity = append(equity, map[string]interface{}{
				"account": acc.Name,
				"balance": balance,
			})
		}
	}

	return map[string]interface{}{
		"as_of_date":        asOfDate,
		"assets":            assets,
		"total_assets":      totalAssets,
		"liabilities":       liabilities,
		"total_liabilities": totalLiabilities,
		"equity":            equity,
		"total_equity":      totalEquity,
		"check":             totalAssets == (totalLiabilities + totalEquity),
	}, nil
}
