package transactions

// Package transactions provides services for handling transactions, including saving and processing them from images or text inputs.

import (
	"context"
	spreadsheet "money-tracker-bot/internal/adapters/google/spreadsheet"
	transaction_domain "money-tracker-bot/internal/domain/transactions"
	aiport "money-tracker-bot/internal/port/out/ai"
	"os"
)

type TransactionService struct {
	DefaultAiPort      aiport.AiPort
	SpreadsheetService SpreadsheetServicePort
}

// SpreadsheetServicePort abstracts spreadsheet operations for testability
type SpreadsheetServicePort interface {
	AppendRow(ctx context.Context, spreadsheetId string, trx transaction_domain.Transaction) spreadsheet.CategorySummary
	GetCellValue(ctx context.Context, spreadsheetId string)
}

func NewTransactionService(ai aiport.AiPort, sheets SpreadsheetServicePort) *TransactionService {
	return &TransactionService{
		DefaultAiPort:      ai,
		SpreadsheetService: sheets,
	}
}

func (t *TransactionService) SaveTransaction(trx transaction_domain.Transaction) (spreadsheet.CategorySummary, error) {
	spreadsheetId := os.Getenv("GOOGLE_SPREADSHEET_ID")
	summary := t.SpreadsheetService.AppendRow(context.Background(), spreadsheetId, trx)
	return summary, nil

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
