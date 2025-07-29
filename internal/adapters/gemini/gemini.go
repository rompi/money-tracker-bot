package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	transaction_domain "money-tracker-bot/internal/domain/transactions"
	"os"
	"strings"
	"time"

	"money-tracker-bot/internal/common"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GenerativeModelPort abstracts the generative model for testability
type GenerativeModelPort interface {
	GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error)
}

// GeminiClient is a client for communicating with the Gemini API
type GeminiClient struct {
	GenAi *genai.Client
	Model GenerativeModelPort
}

// NewClient creates a new GeminiClient
func NewClient(apiKey string) *GeminiClient {
	client, _ := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	return &GeminiClient{
		GenAi: client,
		Model: client.GenerativeModel("gemini-2.0-flash"),
	}
}

// GenerateContent sends a prompt to Gemini and returns the response text
func (c *GeminiClient) GenerateContent(ctx context.Context, prompt string) {
	_, err := c.Model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}
	// For testability, we do not process the response here
}

func (c *GeminiClient) ReadImageToTransaction(ctx context.Context, imgPath string) (*transaction_domain.Transaction, error) {
	imgData, err := os.ReadFile(imgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}

	fileID := ""
	if parts := strings.Split(imgPath, "/"); len(parts) > 0 {
		fileID = parts[len(parts)-1]
	}

	prompt := common.BuildPrompt(common.PromptParams{
		IsImage: true,
		FileID:  fileID,
	})

	req := []genai.Part{
		genai.ImageData("jpeg", imgData),
		genai.Text(prompt),
	}

	resp, err := c.Model.GenerateContent(ctx, req...)
	if err != nil {
		return nil, fmt.Errorf("gemini generate content error: %w", err)
	}

	var transaction transaction_domain.Transaction
	for _, cand := range resp.Candidates {
		if cand.Content == nil || len(cand.Content.Parts) == 0 {
			continue
		}
		var jsonText string
		for _, part := range cand.Content.Parts {
			if textPart, ok := part.(genai.Text); ok {
				jsonText += fmt.Sprintf("%s", textPart)
			}
		}
		jsonText = trimJson(jsonText)
		if err := json.Unmarshal([]byte(jsonText), &transaction); err != nil {
			log.Printf("Failed to parse JSON: %v\nResponse:\n%s", err, jsonText)
			continue
		}
		// Ensure amount is positive
		transaction.Amount = ensurePositiveAmount(transaction.Amount)
	}

	if err := os.Remove(imgPath); err != nil {
		log.Printf("Failed to remove file %s: %v", imgPath, err)
	}
	return &transaction, nil
}

func (c *GeminiClient) TextToTransaction(ctx context.Context, message string) (*transaction_domain.Transaction, error) {
	currentDate := time.Now().Format("2006-01-02")

	prompt := common.BuildPrompt(common.PromptParams{
		IsImage:     false,
		Message:     message,
		CurrentDate: currentDate,
	})

	req := []genai.Part{
		genai.Text(prompt),
	}

	resp, err := c.Model.GenerateContent(ctx, req...)
	if err != nil {
		return nil, fmt.Errorf("gemini generate content error: %w", err)
	}

	var transaction transaction_domain.Transaction
	for _, cand := range resp.Candidates {
		if cand.Content == nil || len(cand.Content.Parts) == 0 {
			continue
		}
		var jsonText string
		for _, part := range cand.Content.Parts {
			if textPart, ok := part.(genai.Text); ok {
				jsonText += fmt.Sprintf("%s", textPart)
			}
		}
		jsonText = trimJson(jsonText)
		if err := json.Unmarshal([]byte(jsonText), &transaction); err != nil {
			log.Printf("Failed to parse JSON: %v\nResponse:\n%s", err, jsonText)
			continue
		}
		// Ensure amount is positive
		transaction.Amount = ensurePositiveAmount(transaction.Amount)
	}
	return &transaction, nil
}

func trimJson(jsonText string) string {
	jsonText = strings.TrimSpace(jsonText)
	jsonText = strings.TrimPrefix(jsonText, "```json")
	jsonText = strings.TrimSuffix(jsonText, "```")
	jsonText = strings.TrimSpace(jsonText)
	return jsonText
}

// ensurePositiveAmount ensures the amount string is positive by removing any negative signs
func ensurePositiveAmount(amount string) string {
	// Remove any negative signs from the amount
	return strings.TrimPrefix(amount, "-")
}
