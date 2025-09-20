package transactions

import (
	"context"
	"money-tracker-bot/internal/adapters/google/spreadsheet"
	transaction_domain "money-tracker-bot/internal/domain/transactions"
	"testing"
)

type mockAiPort struct{}

func (m *mockAiPort) GenerateContent(ctx context.Context, prompt string) error { return nil }
func (m *mockAiPort) ReadImageToTransaction(ctx context.Context, imagePath string) (*transaction_domain.Transaction, error) {
	return &transaction_domain.Transaction{Title: "mocked"}, nil
}
func (m *mockAiPort) TextToTransaction(ctx context.Context, message string) (*transaction_domain.Transaction, error) {
	return &transaction_domain.Transaction{Title: "mocked"}, nil
}

// DummySpreadsheetService implements only the methods needed for TransactionService
type DummySpreadsheetService struct{}

func (d *DummySpreadsheetService) AppendRow(ctx context.Context, spreadsheetId string, trx transaction_domain.Transaction) (spreadsheet.CategorySummary, error) {
	return spreadsheet.CategorySummary{}, nil
}
func (d *DummySpreadsheetService) GetCellValue(ctx context.Context, spreadsheetId string) error {
	return nil
}

func TestSaveTransaction(t *testing.T) {
	ts := &TransactionService{
		DefaultAiPort:      &mockAiPort{},
		SpreadsheetService: &DummySpreadsheetService{},
	}
	trx := transaction_domain.Transaction{Title: "test"}
	summary, err := ts.SaveTransaction(trx)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	// CategorySummary is empty in test but that's ok for this test
	_ = summary // we don't need to validate the summary contents in this test
}

func TestHandleImageInput(t *testing.T) {
	ts := &TransactionService{DefaultAiPort: &mockAiPort{}}
	trx, err := ts.HandleImageInput(context.Background(), "img.jpg", "user", nil)
	if err != nil || trx.Title != "mocked" {
		t.Errorf("unexpected result: %v, %v", trx, err)
	}
}

func TestHandleTextInput(t *testing.T) {
	ts := &TransactionService{DefaultAiPort: &mockAiPort{}}
	trx, err := ts.HandleTextInput(context.Background(), "img.jpg", "user", nil)
	if err != nil || trx.Title != "mocked" {
		t.Errorf("unexpected result: %v, %v", trx, err)
	}
}
