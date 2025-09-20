package transactions

import (
	"context"
	spreadsheet "money-tracker-bot/internal/adapters/google/spreadsheet"
	transaction_domain "money-tracker-bot/internal/domain/transactions"
	aiport "money-tracker-bot/internal/port/out/ai"
)

type ITransaction interface {
	// SaveTransactions saves the transactions to the database
	SaveTransaction(trx transaction_domain.Transaction) (spreadsheet.CategorySummary, error)
	HandleImageInput(context.Context, string, string, aiport.AiPort) (*transaction_domain.Transaction, error)
	HandleTextInput(context.Context, string, string, aiport.AiPort) (*transaction_domain.Transaction, error)
}
