package gemini

import (
	"context"
	"testing"

	"github.com/google/generative-ai-go/genai"
)

type mockModel struct {
	GenerateContentCalled bool
}

func (m *mockModel) GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error) {
	m.GenerateContentCalled = true
	return &genai.GenerateContentResponse{}, nil
}

type mockGenAi struct{}

func (m *mockGenAi) GenerativeModel(name string) *mockModel {
	return &mockModel{}
}

func TestGeminiClient_GenerateContent(t *testing.T) {
	client := &GeminiClient{
		GenAi: nil,
		Model: &mockModel{},
	}
	client.GenerateContent(context.Background(), "test prompt")
	// No panic means pass for this stub
}

func TestGeminiClient_TextToTransaction(t *testing.T) {
	client := &GeminiClient{
		GenAi: nil,
		Model: &mockModel{},
	}
	trx, err := client.TextToTransaction(context.Background(), "test message")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if trx == nil {
		t.Errorf("expected transaction, got nil")
	}
}
