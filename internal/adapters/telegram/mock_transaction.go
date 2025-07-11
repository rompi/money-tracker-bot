package telegram

import (
	"context"
	transaction_domain "money-tracker-bot/internal/domain/transactions"
	aiport "money-tracker-bot/internal/port/out/ai"
)

type MockTransactionService struct {
	HandleTextInputCalled  bool
	HandleImageInputCalled bool
	SaveTransactionCalled  bool
}

func (m *MockTransactionService) HandleTextInput(ctx context.Context, text, user string, ai aiport.AiPort) (*transaction_domain.Transaction, error) {
	m.HandleTextInputCalled = true
	return &transaction_domain.Transaction{Notes: "test notes", Amount: "1000"}, nil
}
func (m *MockTransactionService) HandleImageInput(ctx context.Context, path, user string, ai aiport.AiPort) (*transaction_domain.Transaction, error) {
	m.HandleImageInputCalled = true
	return &transaction_domain.Transaction{Notes: "img notes", Amount: "2000"}, nil
}
func (m *MockTransactionService) SaveTransaction(tx transaction_domain.Transaction) error {
	m.SaveTransactionCalled = true
	return nil
}
