package transactions

import (
	"context"
	transaction_domain "rompi/gobot/internal/domain/transactions"
	aiport "rompi/gobot/internal/port/out/ai"
)

type ITransaction interface {
	// SaveTransactions saves the transactions to the database
	SaveTransaction(trx transaction_domain.Transaction) error
	HandleImageInput(context.Context, string, string, aiport.AiPort) (*transaction_domain.Transaction, error)
}
