package main

import (
	"accounting/config"
	"accounting/internal/domain"
	"accounting/internal/handler"
	"accounting/internal/repo"
	"accounting/internal/service"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.LoadConfig()
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		cfg.Database.URL = dbURL
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer pool.Close()

	// Initialize Store
	store := repo.NewStore(pool)

	// Run Migrations
	if err := store.Migrate(ctx, "internal/repo/schema/000001_init.up.sql"); err != nil {
		log.Printf("Migration warning (may already exist): %v", err)
	}

	// Initialize Services
	accountSvc := service.NewAccountService(store)
	ledgerSvc := service.NewLedgerService(store)
	closingSvc := service.NewClosingService(store, ledgerSvc)
	reporterSvc := service.NewReporterService(store)

	// Seed initial accounts (Generic)
	seedAccounts(ctx, accountSvc)

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

func seedAccounts(ctx context.Context, svc *service.AccountService) {
	// Check if already seeded
	accs, _ := svc.ListAccounts(ctx)
	if len(accs) > 0 {
		return
	}

	accounts := []domain.CreateAccountRequest{
		{Code: "1000", Name: "Cash", Type: domain.Asset},
		{Code: "1100", Name: "Accounts Receivable", Type: domain.Asset},
		{Code: "2000", Name: "Accounts Payable", Type: domain.Liability},
		{Code: "3000", Name: "Owner's Capital", Type: domain.Equity},
		{Code: "3100", Name: "Retained Earnings", Type: domain.Equity},
		{Code: "4000", Name: "General Sales", Type: domain.Revenue},
		{Code: "5000", Name: "COGS", Type: domain.Expense},
		{Code: "6000", Name: "General Expenses", Type: domain.Expense},
	}

	for _, acc := range accounts {
		svc.CreateAccount(ctx, acc)
	}
}
