package aiport

import "testing"

import (
	"context"
	transaction_domain "money-tracker-bot/internal/domain/transactions"
)

type DummyAiPort struct{}

func (d *DummyAiPort) GenerateContent(ctx context.Context, prompt string) error { return nil }
func (d *DummyAiPort) ReadImageToTransaction(ctx context.Context, imagePath string) (*transaction_domain.Transaction, error) {
	return nil, nil
}
func (d *DummyAiPort) TextToTransaction(ctx context.Context, message string) (*transaction_domain.Transaction, error) {
	return nil, nil
}

func TestDummyAiPort_ImplementsAiPort(t *testing.T) {
	var _ AiPort = &DummyAiPort{}
}

func TestDummyAiPort_ReadImageToTransaction_ReturnsNil(t *testing.T) {
	dummy := &DummyAiPort{}
	tx, err := dummy.ReadImageToTransaction(context.Background(), "dummy_path.jpg")
	if tx != nil {
		t.Errorf("Expected nil transaction, got %+v", tx)
	}
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// Test with canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	tx, err = dummy.ReadImageToTransaction(ctx, "dummy_path.jpg")
	if tx != nil {
		t.Errorf("Expected nil transaction with canceled context, got %+v", tx)
	}
	if err != nil {
		t.Errorf("Expected nil error with canceled context, got %v", err)
	}
}

func TestDummyAiPort_TextToTransaction_ReturnsNil(t *testing.T) {
	dummy := &DummyAiPort{}
	tx, err := dummy.TextToTransaction(context.Background(), "dummy message")
	if tx != nil {
		t.Errorf("Expected nil transaction, got %+v", tx)
	}
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// Test with canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	tx, err = dummy.TextToTransaction(ctx, "dummy message")
	if tx != nil {
		t.Errorf("Expected nil transaction with canceled context, got %+v", tx)
	}
	if err != nil {
		t.Errorf("Expected nil error with canceled context, got %v", err)
	}
}

func TestDummyAiPort_GenerateContent_NoPanic(t *testing.T) {
	dummy := &DummyAiPort{}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GenerateContent panicked: %v", r)
		}
	}()
	dummy.GenerateContent(context.Background(), "test prompt")
}
