package gemini

import (
	"context"
)

type MockGeminiClient struct{}

func (m *MockGeminiClient) GenerateContent(ctx context.Context, prompt string) {}
func (m *MockGeminiClient) ReadImageToTransaction(ctx context.Context, imgPath string) (interface{}, error) {
	return nil, nil
}
func (m *MockGeminiClient) TextToTransaction(ctx context.Context, message string) (interface{}, error) {
	return nil, nil
}
