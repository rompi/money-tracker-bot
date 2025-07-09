package transactions

import (
  "context"
  "testing"
  transaction_domain "rompi/gobot/internal/domain/transactions"
  "rompi/gobot/internal/adapters/google/spreadsheet"
)


type mockAiPort struct{}

func (m *mockAiPort) GenerateContent(ctx context.Context, prompt string) {}
func (m *mockAiPort) ReadImageToTransaction(ctx context.Context, imagePath string) (*transaction_domain.Transaction, error) {
  return &transaction_domain.Transaction{Title: "mocked"}, nil
}
func (m *mockAiPort) TextToTransaction(ctx context.Context, message string) (*transaction_domain.Transaction, error) {
  return &transaction_domain.Transaction{Title: "mocked"}, nil
}

// DummySpreadsheetService implements only the methods needed for TransactionService
type DummySpreadsheetService struct{}
func (d *DummySpreadsheetService) AppendRow(ctx context.Context, spreadsheetId string, trx transaction_domain.Transaction) {}
func (d *DummySpreadsheetService) GetCellValue(ctx context.Context, spreadsheetId string) {}

func TestSaveTransaction(t *testing.T) {
  ts := &TransactionService{
    DefaultAiPort:      &mockAiPort{},
    SpreadsheetService: &spreadsheet.MockSpreadsheetService{},
  }
  trx := transaction_domain.Transaction{Title: "test"}
  err := ts.SaveTransaction(trx)
  if err != nil {
    t.Errorf("expected no error, got: %v", err)
  }
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
