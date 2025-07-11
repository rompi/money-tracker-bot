package spreadsheet

import (
	"context"
	transaction_domain "money-tracker-bot/internal/domain/transactions"
)

type MockSpreadsheetService struct{}

func (m *MockSpreadsheetService) AppendRow(ctx context.Context, spreadsheetId string, trx transaction_domain.Transaction) {
}
func (m *MockSpreadsheetService) GetCellValue(ctx context.Context, spreadsheetId string) {}

// Add more mock methods as needed for your tests
