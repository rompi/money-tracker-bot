package aiport

import (
	"context"
	transaction_domain "money-tracker-bot/internal/domain/transactions"
)

type AiPort interface {
	GenerateContent(ctx context.Context, prompt string)
	ReadImageToTransaction(ctx context.Context, imgPath string) (*transaction_domain.Transaction, error)
	TextToTransaction(ctx context.Context, message string) (*transaction_domain.Transaction, error)
}
