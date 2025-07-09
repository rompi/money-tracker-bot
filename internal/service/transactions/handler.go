package transactions

import (
  "context"
  "os"
  transaction_domain "rompi/gobot/internal/domain/transactions"
  aiport "rompi/gobot/internal/port/out/ai"
)

type TransactionService struct {
  DefaultAiPort      aiport.AiPort
  SpreadsheetService SpreadsheetServicePort
}

// SpreadsheetServicePort abstracts spreadsheet operations for testability
type SpreadsheetServicePort interface {
  AppendRow(ctx context.Context, spreadsheetId string, trx transaction_domain.Transaction)
  GetCellValue(ctx context.Context, spreadsheetId string)
}

func NewTransactionService(ai aiport.AiPort, sheets SpreadsheetServicePort) *TransactionService {
   return &TransactionService{
	   DefaultAiPort:      ai,
	   SpreadsheetService: sheets,
   }
}

func (t *TransactionService) SaveTransaction(trx transaction_domain.Transaction) error {
	// Implement the logic to save the transaction
	// This is a placeholder implementation
	spreadsheetId := os.Getenv("GOOGLE_SPREADSHEET_ID")
	t.SpreadsheetService.AppendRow(context.Background(), spreadsheetId, trx)
	return nil

}

func (t *TransactionService) HandleImageInput(ctx context.Context, imagePath string, uploader string, aiPort aiport.AiPort) (*transaction_domain.Transaction, error) {
	ai := t.DefaultAiPort
	if aiPort != nil {
		ai = aiPort
	}

	trx, err := ai.ReadImageToTransaction(ctx, imagePath)
	if err != nil {
		return nil, err
	}
	trx.CreatedBy = uploader
	return trx, nil
}

func (t *TransactionService) HandleTextInput(ctx context.Context, imagePath string, uploader string, aiPort aiport.AiPort) (*transaction_domain.Transaction, error) {
	ai := t.DefaultAiPort
	if aiPort != nil {
		ai = aiPort
	}

	trx, err := ai.TextToTransaction(ctx, imagePath)
	if err != nil {
		return nil, err
	}
	trx.CreatedBy = uploader
	return trx, nil
}
