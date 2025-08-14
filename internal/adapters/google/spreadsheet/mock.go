package spreadsheet

import (
	"context"
	transaction_domain "money-tracker-bot/internal/domain/transactions"
)

type MockSpreadsheetService struct{}

func (m *MockSpreadsheetService) AppendRow(ctx context.Context, spreadsheetId string, trx transaction_domain.Transaction) CategorySummary {
	// Return mock data with all fields populated for testing
	return CategorySummary{
		Category:        "MockCategory",
		MonthlyExpenses: "1000",
		MonthlyBudget:   "5000",
		BudgetLeft:      "4000",
		Quota:           "2000",
		QuotaLeft:       "1500",
	}
}
func (m *MockSpreadsheetService) GetCellValue(ctx context.Context, spreadsheetId string) {}

// Add more mock methods as needed for your tests
