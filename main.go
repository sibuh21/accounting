package main

import (
	"accounting/config"
	"accounting/internal/domain"
	"accounting/internal/handler"
	"accounting/internal/repo"
	"accounting/internal/service"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize Store
	store := repo.NewStore()

	// Initialize Services
	accountSvc := service.NewAccountService(store)
	ledgerSvc := service.NewLedgerService(store)
	closingSvc := service.NewClosingService(store, ledgerSvc)
	reporterSvc := service.NewReporterService(store)

	// Seed initial accounts
	seedAccounts(accountSvc)

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	g := e.Group("/v1")
	handler.NewAccountingHandler(g, accountSvc, ledgerSvc, closingSvc, reporterSvc)

	// Start Server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	log.Printf("Accounting service starting on :%s", cfg.Server.Port)
	if err := e.StartServer(srv); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %v", err)
	}
}

func seedAccounts(svc *service.AccountService) {
	accounts := []domain.CreateAccountRequest{
		{Code: "1000", Name: "Cash", Type: domain.Asset},
		{Code: "1100", Name: "Accounts Receivable", Type: domain.Asset},
		{Code: "1200", Name: "Inventory", Type: domain.Asset},
		{Code: "1300", Name: "Equipment", Type: domain.Asset},
		{Code: "2000", Name: "Accounts Payable", Type: domain.Liability},
		{Code: "2100", Name: "Bank Loan", Type: domain.Liability},
		{Code: "3000", Name: "Owner's Capital", Type: domain.Equity},
		{Code: "3100", Name: "Retained Earnings", Type: domain.Equity},
		{Code: "3200", Name: "Owner's Draw", Type: domain.Equity},
		{Code: "4000", Name: "Coffee Sales", Type: domain.Revenue},
		{Code: "4100", Name: "Pastry Sales", Type: domain.Revenue},
		{Code: "5000", Name: "COGS - Coffee Beans", Type: domain.Expense},
		{Code: "5100", Name: "COGS - Milk", Type: domain.Expense},
		{Code: "6000", Name: "Wages Expense", Type: domain.Expense},
		{Code: "6100", Name: "Rent Expense", Type: domain.Expense},
		{Code: "6200", Name: "Utilities Expense", Type: domain.Expense},
		{Code: "6300", Name: "Supplies Expense", Type: domain.Expense},
	}

	for _, acc := range accounts {
		svc.CreateAccount(acc)
	}
}
