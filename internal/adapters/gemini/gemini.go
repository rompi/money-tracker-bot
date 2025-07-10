package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	transaction_domain "rompi/gobot/internal/domain/transactions"
	"strings"
	"time"

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

// buildPrompt centralizes the construction of prompts for Gemini
func buildPrompt(params map[string]string) string {
	fileID := params["fileID"]
	message := params["message"]
	currentDate := params["currentDate"]
	isImage := params["isImage"] == "true"

	fields := `Fields:
- title (summary of the transaction notes)
- transaction_date (format always YYYY-MM-DD)
- amount (in rupiah, format 1,000,000 for transaction with amont 1 million. if it is 100k then output should be 100,000)
- notes (details of the transaction, containing items bought)
- category (Groceries / Utilities / Entertainment / Gifting / Household / Eating Out / Health / Transportation / Savings / Emergency / Rent House)`
	if isImage {
		fields += `
- destination_number
- source_account (only GOPAY / BCA / OVO / DANA / ISAKU / MANDIRI / BNI / BRI / CASH)
- file_id ` + fileID
	} else {
		fields += `
- file_id should be empty`
	}

	var inputDesc string
	if isImage {
		inputDesc = "from the image"
	} else {
		inputDesc = fmt.Sprintf("from the following message: %s", message)
	}

	var dateLine string
	if isImage {
		dateLine = ""
	} else {
		dateLine = fmt.Sprintf("- transaction_date should be %s (format always YYYY-MM-DD)\n", currentDate)
	}

	exampleFileID := fileID
	if !isImage {
		exampleFileID = ""
	}

	prompt := fmt.Sprintf(`Please extract the following data %s and return it as valid JSON.

%s
%sIMPORTANT:
Respond ONLY with raw JSON.
No explanation, no formatting, no code blocks.

Example:
{
	"title": "Transfer to ABC Cafe",
	"transaction_date": "2025-03-30",
	"amount": "150",
	"notes": "Lunch at ABC cafe",
	"destination_number": "0524012911",
	"source_account": "Gopay",
	"category": "Groceries",
	"file_id": "%s"
}`,
		inputDesc, fields, dateLine, exampleFileID)
	return prompt
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

	prompt := buildPrompt(map[string]string{
		"isImage": "true",
		"fileID":  fileID,
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
	}

	if err := os.Remove(imgPath); err != nil {
		log.Printf("Failed to remove file %s: %v", imgPath, err)
	}
	return &transaction, nil
}

func (c *GeminiClient) TextToTransaction(ctx context.Context, message string) (*transaction_domain.Transaction, error) {
	currentDate := time.Now().Format("2006-01-02")

	prompt := buildPrompt(map[string]string{
		"isImage":     "false",
		"message":     message,
		"currentDate": currentDate,
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
