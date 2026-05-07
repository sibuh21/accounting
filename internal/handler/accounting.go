package handler

import (
	"accounting/internal/domain"
	"accounting/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type AccountingHandler struct {
	accountSvc  *service.AccountService
	ledgerSvc   *service.LedgerService
	closingSvc  *service.ClosingService
	reporterSvc *service.ReporterService
}

func NewAccountingHandler(
	g *echo.Group,
	accountSvc *service.AccountService,
	ledgerSvc *service.LedgerService,
	closingSvc *service.ClosingService,
	reporterSvc *service.ReporterService,
) {
	h := &AccountingHandler{
		accountSvc:  accountSvc,
		ledgerSvc:   ledgerSvc,
		closingSvc:  closingSvc,
		reporterSvc: reporterSvc,
	}

	g.POST("/accounts", h.CreateAccount)
	g.GET("/accounts", h.ListAccounts)
	g.POST("/entries", h.PostEntry)
	g.GET("/entries", h.ListEntries)
	g.GET("/reports/trial-balance", h.GetTrialBalance)
	g.GET("/reports/income-statement", h.GetIncomeStatement)
	g.GET("/reports/balance-sheet", h.GetBalanceSheet)
	g.POST("/closing", h.CloseMonth)
}

func (h *AccountingHandler) CreateAccount(c echo.Context) error {
	var req domain.CreateAccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	acc, err := h.accountSvc.CreateAccount(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, acc)
}

func (h *AccountingHandler) ListAccounts(c echo.Context) error {
	accounts, err := h.accountSvc.ListAccounts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, accounts)
}

func (h *AccountingHandler) PostEntry(c echo.Context) error {
	var req domain.CreateJournalEntryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	entry := domain.NewJournalEntry(req.Date, req.Description)
	if req.Date.IsZero() {
		entry.Date = time.Now()
	}

	for _, l := range req.Lines {
		entry.Lines = append(entry.Lines, domain.EntryLine{
			AccountID:   l.AccountID,
			AccountName: l.AccountName,
			Debit:       l.Debit,
			Credit:      l.Credit,
		})
	}

	if err := h.ledgerSvc.PostEntry(c.Request().Context(), entry); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, entry)
}

func (h *AccountingHandler) ListEntries(c echo.Context) error {
	entries, err := h.accountSvc.GetEntries(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, entries)
}

func (h *AccountingHandler) GetTrialBalance(c echo.Context) error {
	tb, err := h.ledgerSvc.GetTrialBalance(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tb)
}

func (h *AccountingHandler) GetIncomeStatement(c echo.Context) error {
	monthQuery := c.QueryParam("month")
	month := int(time.Now().Month())
	if monthQuery != "" {
		if m, err := strconv.Atoi(monthQuery); err == nil && m >= 1 && m <= 12 {
			month = m
		}
	}

	startDate := time.Date(time.Now().Year(), time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	report, err := h.reporterSvc.GenerateIncomeStatement(c.Request().Context(), startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, report)
}

func (h *AccountingHandler) GetBalanceSheet(c echo.Context) error {
	report, err := h.reporterSvc.GenerateBalanceSheet(c.Request().Context(), time.Now())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, report)
}

func (h *AccountingHandler) CloseMonth(c echo.Context) error {
	monthEnd := time.Now().AddDate(0, 0, -time.Now().Day()).AddDate(0, 1, -1)
	entry, err := h.closingSvc.CloseMonth(c.Request().Context(), monthEnd)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, entry)
}
