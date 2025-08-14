package spreadsheet

import (
	"context"
	transaction_domain "money-tracker-bot/internal/domain/transactions"
	"testing"
)

func TestCategorySummary_WithQuota(t *testing.T) {
	service := MockSpreadsheetService{}
	trx := transaction_domain.Transaction{Category: "MockCategory"}
	result := service.AppendRow(context.Background(), "dummy", trx)

	if result.Category != "MockCategory" {
		t.Errorf("expected Category 'MockCategory', got '%s'", result.Category)
	}
	if result.MonthlyExpenses != "1000" {
		t.Errorf("expected MonthlyExpenses '1000', got '%s'", result.MonthlyExpenses)
	}
	if result.MonthlyBudget != "5000" {
		t.Errorf("expected MonthlyBudget '5000', got '%s'", result.MonthlyBudget)
	}
	if result.BudgetLeft != "4000" {
		t.Errorf("expected BudgetLeft '4000', got '%s'", result.BudgetLeft)
	}
	if result.Quota != "2000" {
		t.Errorf("expected Quota '2000', got '%s'", result.Quota)
	}
	if result.QuotaLeft != "1500" {
		t.Errorf("expected QuotaLeft '1500', got '%s'", result.QuotaLeft)
	}
}

func TestCategorySummary_MissingQuota(t *testing.T) {
	// Simulate a row with only 4 columns (no quota info)
	row := []interface{}{"Food", "1000", "5000", "4000"}
	var quota, quotaLeft string
	if len(row) > 4 {
		quota = row[4].(string)
	}
	if len(row) > 5 {
		quotaLeft = row[5].(string)
	}
	result := CategorySummary{
		Category:        row[0].(string),
		MonthlyExpenses: row[1].(string),
		MonthlyBudget:   row[2].(string),
		BudgetLeft:      row[3].(string),
		Quota:           quota,
		QuotaLeft:       quotaLeft,
	}
	if result.Quota != "" {
		t.Errorf("expected empty Quota, got '%s'", result.Quota)
	}
	if result.QuotaLeft != "" {
		t.Errorf("expected empty QuotaLeft, got '%s'", result.QuotaLeft)
	}
}
