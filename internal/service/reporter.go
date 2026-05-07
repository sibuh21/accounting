package service

import (
	"accounting/internal/repo"
	"time"
)

type ReporterService struct {
	store *repo.Store
}

func NewReporterService(store *repo.Store) *ReporterService {
	return &ReporterService{store: store}
}

func (rs *ReporterService) GenerateIncomeStatement(startDate, endDate time.Time) map[string]interface{} {
	accounts := rs.store.GetAllAccounts()

	var totalRevenue float64
	var totalExpenses float64
	revenues := []map[string]interface{}{}
	expenses := []map[string]interface{}{}

	for _, acc := range accounts {
		if acc.Type == "Revenue" {
			totalRevenue += acc.Balance
			revenues = append(revenues, map[string]interface{}{
				"account": acc.Name,
				"amount":  acc.Balance,
			})
		} else if acc.Type == "Expense" {
			totalExpenses += acc.Balance
			expenses = append(expenses, map[string]interface{}{
				"account": acc.Name,
				"amount":  acc.Balance,
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
	}
}

func (rs *ReporterService) GenerateBalanceSheet(asOfDate time.Time) map[string]interface{} {
	accounts := rs.store.GetAllAccounts()

	var totalAssets float64
	var totalLiabilities float64
	var totalEquity float64

	assets := []map[string]interface{}{}
	liabilities := []map[string]interface{}{}
	equity := []map[string]interface{}{}

	for _, acc := range accounts {
		switch acc.Type {
		case "Asset":
			totalAssets += acc.Balance
			assets = append(assets, map[string]interface{}{
				"account": acc.Name,
				"balance": acc.Balance,
			})
		case "Liability":
			totalLiabilities += acc.Balance
			liabilities = append(liabilities, map[string]interface{}{
				"account": acc.Name,
				"balance": acc.Balance,
			})
		case "Equity":
			totalEquity += acc.Balance
			equity = append(equity, map[string]interface{}{
				"account": acc.Name,
				"balance": acc.Balance,
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
	}
}
