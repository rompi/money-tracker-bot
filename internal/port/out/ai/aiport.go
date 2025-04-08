package aiport

import (
	"context"
	transaction_domain "rompi/gobot/internal/domain/transactions"
)

type AiPort interface {
	GenerateContent(ctx context.Context, prompt string)
	ReadImageToTransaction(ctx context.Context, imgPath string) (*transaction_domain.Transaction, error)
}
