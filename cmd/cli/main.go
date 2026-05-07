package main

import (
	"accounting/config"
	"accounting/internal/domain"
	"accounting/internal/repo"
	"accounting/internal/service"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountingCLI struct {
	accountSvc  *service.AccountService
	ledgerSvc   *service.LedgerService
	closingSvc  *service.ClosingService
	reporterSvc *service.ReporterService
	ctx         context.Context
}

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

	store := repo.NewStore(pool)
	accountSvc := service.NewAccountService(store)
	ledgerSvc := service.NewLedgerService(store)
	closingSvc := service.NewClosingService(store, ledgerSvc)
	reporterSvc := service.NewReporterService(store)

	cli := &AccountingCLI{
		accountSvc:  accountSvc,
		ledgerSvc:   ledgerSvc,
		closingSvc:  closingSvc,
		reporterSvc: reporterSvc,
		ctx:         ctx,
	}

	cli.run()
}

func (cli *AccountingCLI) run() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		cli.printMenu()
		fmt.Print("\nEnter choice: ")
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			cli.recordSales()
		case "2":
			cli.recordExpense()
		case "3":
			cli.recordPurchaseInventory()
		case "4":
			cli.recordOwnerDraw()
		case "5":
			cli.viewTrialBalance()
		case "6":
			cli.viewIncomeStatement()
		case "7":
			cli.viewBalanceSheet()
		case "8":
			cli.closeMonth()
		case "9":
			cli.viewAllEntries()
		case "0":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func (cli *AccountingCLI) printMenu() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("         ACCOUNTING SYSTEM CLI")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("1. Record Daily Sales")
	fmt.Println("2. Record Expense")
	fmt.Println("3. Record Inventory Purchase")
	fmt.Println("4. Record Owner's Draw")
	fmt.Println("5. View Trial Balance")
	fmt.Println("6. View Income Statement (P&L)")
	fmt.Println("7. View Balance Sheet")
	fmt.Println("8. Close Month (Run Closing Entries)")
	fmt.Println("9. View All Journal Entries")
	fmt.Println("0. Exit")
}

func (cli *AccountingCLI) recordSales() {
	fmt.Print("\nEnter sales amount: $")
	amount := cli.getFloatInput()
	fmt.Print("Enter payment method (cash/credit): ")
	paymentMethod := cli.getStringInput()

	entry := domain.NewJournalEntry(time.Now(), fmt.Sprintf("Sales - %s", paymentMethod))

	if strings.ToLower(paymentMethod) == "cash" {
		entry.Lines = append(entry.Lines, domain.EntryLine{AccountName: "Cash", Debit: amount})
	} else {
		entry.Lines = append(entry.Lines, domain.EntryLine{AccountName: "Accounts Receivable", Debit: amount})
	}

	entry.Lines = append(entry.Lines, domain.EntryLine{AccountName: "General Sales", Credit: amount})

	if err := cli.resolveAccountIDs(entry); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if err := cli.ledgerSvc.PostEntry(cli.ctx, entry); err != nil {
		fmt.Printf("Error posting entry: %v\n", err)
		return
	}

	fmt.Printf("\n✓ Recorded sales: $%.2f\n", amount)
}

func (cli *AccountingCLI) recordExpense() {
	fmt.Print("Enter expense account name: ")
	expenseAccount := cli.getStringInput()
	fmt.Print("Enter amount: $")
	amount := cli.getFloatInput()
	fmt.Print("Enter description: ")
	description := cli.getStringInput()

	entry := domain.NewJournalEntry(time.Now(), description)
	entry.Lines = append(entry.Lines, domain.EntryLine{AccountName: expenseAccount, Debit: amount})
	entry.Lines = append(entry.Lines, domain.EntryLine{AccountName: "Cash", Credit: amount})

	if err := cli.resolveAccountIDs(entry); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if err := cli.ledgerSvc.PostEntry(cli.ctx, entry); err != nil {
		fmt.Printf("Error posting entry: %v\n", err)
		return
	}

	fmt.Printf("\n✓ Recorded expense: $%.2f\n", amount)
}

func (cli *AccountingCLI) recordPurchaseInventory() {
	fmt.Print("\nEnter inventory description: ")
	desc := cli.getStringInput()
	fmt.Print("Enter total cost: $")
	totalCost := cli.getFloatInput()
	fmt.Print("Payment method (cash/credit): ")
	paymentMethod := cli.getStringInput()

	entry := domain.NewJournalEntry(time.Now(), fmt.Sprintf("Purchased %v - %s", desc, paymentMethod))
	entry.Lines = append(entry.Lines, domain.EntryLine{AccountName: "Inventory", Debit: totalCost})

	if strings.ToLower(paymentMethod) == "cash" {
		entry.Lines = append(entry.Lines, domain.EntryLine{AccountName: "Cash", Credit: totalCost})
	} else {
		entry.Lines = append(entry.Lines, domain.EntryLine{AccountName: "Accounts Payable", Credit: totalCost})
	}

	if err := cli.resolveAccountIDs(entry); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if err := cli.ledgerSvc.PostEntry(cli.ctx, entry); err != nil {
		fmt.Printf("Error posting entry: %v\n", err)
		return
	}

	fmt.Printf("\n✓ Recorded inventory purchase: $%.2f\n", totalCost)
}

func (cli *AccountingCLI) recordOwnerDraw() {
	fmt.Print("\nEnter draw amount: $")
	amount := cli.getFloatInput()
	fmt.Print("Enter purpose: ")
	purpose := cli.getStringInput()

	entry := domain.NewJournalEntry(time.Now(), fmt.Sprintf("Owner's draw - %s", purpose))
	// Generic Equity account search
	entry.Lines = append(entry.Lines, domain.EntryLine{AccountName: "Owner's Draw", Debit: amount})
	entry.Lines = append(entry.Lines, domain.EntryLine{AccountName: "Cash", Credit: amount})

	if err := cli.resolveAccountIDs(entry); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if err := cli.ledgerSvc.PostEntry(cli.ctx, entry); err != nil {
		fmt.Printf("Error posting entry: %v\n", err)
		return
	}

	fmt.Printf("\n✓ Recorded owner's draw: $%.2f\n", amount)
}

func (cli *AccountingCLI) viewTrialBalance() {
	tb, err := cli.ledgerSvc.GetTrialBalance(cli.ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Printf("%-15s %-30s %15s\n", "Code", "Account Name", "Balance")
	fmt.Println(strings.Repeat("-", 70))

	for _, acc := range tb["accounts"].([]map[string]interface{}) {
		fmt.Printf("%-15s %-30s %15.2f\n", acc["code"], acc["name"], acc["balance"])
	}

	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("%-45s %15.2f\n", "Total Debits:", tb["total_debits"])
	fmt.Printf("%-45s %15.2f\n", "Total Credits:", tb["total_credits"])
	fmt.Println(strings.Repeat("-", 70))

	if tb["is_balanced"].(bool) {
		fmt.Println("✓ Balanced")
	} else {
		fmt.Println("✗ Not Balanced")
	}
}

func (cli *AccountingCLI) viewIncomeStatement() {
	fmt.Print("\nEnter month (1-12): ")
	month := cli.getIntInput()
	startDate := time.Date(time.Now().Year(), time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	pnl, err := cli.reporterSvc.GenerateIncomeStatement(cli.ctx, startDate, endDate)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("\nINCOME STATEMENT for %s to %s\n", 
		pnl["period_start"].(time.Time).Format("2006-01-02"),
		pnl["period_end"].(time.Time).Format("2006-01-02"))
	
	fmt.Println("\nREVENUES:")
	for _, rev := range pnl["revenues"].([]map[string]interface{}) {
		fmt.Printf("  %-30s %15.2f\n", rev["account"], rev["amount"])
	}
	fmt.Printf("Total Revenue: %15.2f\n", pnl["total_revenue"])

	fmt.Println("\nEXPENSES:")
	for _, exp := range pnl["expenses"].([]map[string]interface{}) {
		fmt.Printf("  %-30s %15.2f\n", exp["account"], exp["amount"])
	}
	fmt.Printf("Total Expenses: %15.2f\n", pnl["total_expenses"])
	fmt.Printf("\nNET INCOME: %15.2f\n", pnl["net_income"])
}

func (cli *AccountingCLI) viewBalanceSheet() {
	bs, err := cli.reporterSvc.GenerateBalanceSheet(cli.ctx, time.Now())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("\nBALANCE SHEET as of %s\n", bs["as_of_date"].(time.Time).Format("2006-01-02"))
	
	fmt.Println("\nASSETS:")
	for _, asset := range bs["assets"].([]map[string]interface{}) {
		fmt.Printf("  %-30s %15.2f\n", asset["account"], asset["balance"])
	}
	fmt.Printf("Total Assets: %15.2f\n", bs["total_assets"])

	fmt.Println("\nLIABILITIES:")
	for _, lb := range bs["liabilities"].([]map[string]interface{}) {
		fmt.Printf("  %-30s %15.2f\n", lb["account"], lb["balance"])
	}
	fmt.Printf("Total Liabilities: %15.2f\n", bs["total_liabilities"])

	fmt.Println("\nEQUITY:")
	for _, eq := range bs["equity"].([]map[string]interface{}) {
		fmt.Printf("  %-30s %15.2f\n", eq["account"], eq["balance"])
	}
	fmt.Printf("Total Equity: %15.2f\n", bs["total_equity"])
}

func (cli *AccountingCLI) closeMonth() {
	monthEnd := time.Now().AddDate(0, 0, -time.Now().Day()).AddDate(0, 1, -1)

	entries, err := cli.closingSvc.CloseMonth(cli.ctx, monthEnd)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	for _, entry := range entries {
		fmt.Printf("\n✓ Month closed. Entry ID: %s\n", entry.ID)
	}
}

func (cli *AccountingCLI) viewAllEntries() {
	entries, err := cli.accountSvc.GetEntries(cli.ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	for _, entry := range entries {
		fmt.Printf("\n%s - %s\n", entry.Date.Time.Format("2006-01-02"), entry.Description)
		fmt.Printf("  (Lines output suppressed for brevity in CLI, use API for full details)\n")
	}
}

func (cli *AccountingCLI) resolveAccountIDs(entry *domain.JournalEntry) error {
	accounts, _ := cli.accountSvc.ListAccounts(cli.ctx)
	nameToID := make(map[string]string)
	for _, acc := range accounts {
		nameToID[acc.Name] = acc.ID.String()
		nameToID[acc.Code] = acc.ID.String()
	}
	for i := range entry.Lines {
		id, exists := nameToID[entry.Lines[i].AccountName]
		if !exists { return fmt.Errorf("account not found: %s", entry.Lines[i].AccountName) }
		entry.Lines[i].AccountID = id
	}
	return nil
}

func (cli *AccountingCLI) getFloatInput() float64 {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	val, _ := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)
	return val
}

func (cli *AccountingCLI) getIntInput() int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	val, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	return val
}

func (cli *AccountingCLI) getStringInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
