package gemini

import (
	"context"
	"testing"

	"github.com/google/generative-ai-go/genai"
)

type mockModel struct {
	GenerateContentCalled bool
	ResponseText         string
}

func (m *mockModel) GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error) {
	m.GenerateContentCalled = true
	resp := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []genai.Part{
						genai.Text(m.ResponseText),
					},
				},
			},
		},
	}
	return resp, nil
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
	testCases := []struct {
		name           string
		input         string
		responseJSON  string
		expectedAmount string
	}{
		{
			name:           "Positive amount",
			input:         "spent 100 on groceries",
			responseJSON:  `{"amount": "100", "title": "Groceries", "notes": "test"}`,
			expectedAmount: "100",
		},
		{
			name:           "Negative amount",
			input:         "spent -100 on groceries",
			responseJSON:  `{"amount": "-100", "title": "Groceries", "notes": "test"}`,
			expectedAmount: "100",
		},
		{
			name:           "Amount with currency",
			input:         "spent $100 on groceries",
			responseJSON:  `{"amount": "100", "title": "Groceries", "notes": "test"}`,
			expectedAmount: "100",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockModel := &mockModel{
				ResponseText: tc.responseJSON,
			}
			client := &GeminiClient{
				GenAi: nil,
				Model: mockModel,
			}

			trx, err := client.TextToTransaction(context.Background(), tc.input)
			if err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
			if trx == nil {
				t.Errorf("expected transaction, got nil")
			}
			if trx.Amount != tc.expectedAmount {
				t.Errorf("expected amount %s, got %s", tc.expectedAmount, trx.Amount)
			}
		})
	}
}

func TestEnsurePositiveAmount(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Positive amount",
			input:    "100",
			expected: "100",
		},
		{
			name:     "Negative amount",
			input:    "-100",
			expected: "100",
		},
		{
			name:     "Zero amount",
			input:    "0",
			expected: "0",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Amount with decimals",
			input:    "-100.50",
			expected: "100.50",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ensurePositiveAmount(tc.input)
			if result != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}
