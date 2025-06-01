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

// GeminiClient is a client for communicating with the Gemini API
type GeminiClient struct {
	GenAi *genai.Client
	Model *genai.GenerativeModel
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
	resp, err := c.Model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		fmt.Println(part)
	}
}

func (c *GeminiClient) ReadImageToTransaction(ctx context.Context, imgPath string) (*transaction_domain.Transaction, error) {

	imgData, err := os.ReadFile(imgPath)
	if err != nil {
		log.Fatal(err)
	}

	fileID := strings.Split(imgPath, "/")
	// split imgPath with /

	// Create the request.
	req := []genai.Part{
		genai.ImageData("jpeg", imgData),
		genai.Text(`Please extract the following data from the image and return it as valid JSON.

			Fields:
			- title (summary of the transaction notes)
			- transaction_date (format always YYYY-MM-DD)
			- amount (in rupiah, format 1,000,000 for transaction with amont 1 million. if it is 100k then output should be 100,000)
			- notes (details of the transaction, containing items bought)
			- destination_number
			- source_account (only GOPAY / BCA / OVO / DANA / ISAKU / MANDIRI / BNI / BRI / CASH)
			- category (Groceries / Utilities / Entertainment / Gifting / Household / Eating Out / Health / Transportation / Savings / Emergency / Rent House)
			- file_id ` + fileID[1] + `

			IMPORTANT:
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
				"file_id": "photo_1743586322.jpg"
			}`),
	}

	// Generate content.
	resp, err := c.Model.GenerateContent(ctx, req...)
	if err != nil {
		panic(err)
	}

	var transaction transaction_domain.Transaction

	for _, cand := range resp.Candidates {
		if cand.Content == nil || len(cand.Content.Parts) == 0 {
			continue
		}

		var jsonText string
		for _, part := range cand.Content.Parts {
			if textPart, ok := part.(genai.Text); ok {

				fmt.Println("textPart:", textPart)
				jsonText += fmt.Sprintf("%s", textPart)
			}
		}
		jsonText = trimJson(jsonText)
		errr := json.Unmarshal([]byte(jsonText), &transaction)
		if errr != nil {
			log.Printf("Failed to parse JSON: %v\nResponse:\n%s", errr, jsonText)
			continue
		}
	}

	fmt.Print("Transaction: ", transaction)
	r := os.Remove(imgPath)
	if r != nil {
		log.Printf("Failed to remove file %s: %v", imgPath, r)
	} else {
		log.Printf("File %s removed successfully", imgPath)
	}
	return &transaction, nil
}

func (c *GeminiClient) TextToTransaction(ctx context.Context, message string) (*transaction_domain.Transaction, error) {

	// get time.now in date format
	currentDate := time.Now().Format("2006-01-02")

	// Create the request.
	req := []genai.Part{
		genai.Text(`Please extract the message ` + message + ` and return it as valid JSON.

			Fields:
			- title (summary of the transaction notes)
			- transaction_date should be ` + currentDate + ` (format always YYYY-MM-DD)
			- amount (in rupiah, format 1,000,000 for transaction with amont 1 million. if it is 100k then output should be 100,000)
			- notes (details of the transaction, containing items bought)
			- category (Groceries / Utilities / Entertainment / Gifting / Household / Eating Out / Health / Transportation / Savings / Emergency / Rent House)
			- file_id should be empty

			IMPORTANT:
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
				"file_id": "photo_1743586322.jpg"
			}`),
	}

	// Generate content.
	resp, err := c.Model.GenerateContent(ctx, req...)
	if err != nil {
		panic(err)
	}

	var transaction transaction_domain.Transaction

	for _, cand := range resp.Candidates {
		if cand.Content == nil || len(cand.Content.Parts) == 0 {
			continue
		}

		var jsonText string
		for _, part := range cand.Content.Parts {
			if textPart, ok := part.(genai.Text); ok {

				fmt.Println("textPart:", textPart)
				jsonText += fmt.Sprintf("%s", textPart)
			}
		}
		jsonText = trimJson(jsonText)
		errr := json.Unmarshal([]byte(jsonText), &transaction)
		if errr != nil {
			log.Printf("Failed to parse JSON: %v\nResponse:\n%s", errr, jsonText)
			continue
		}
	}

	fmt.Print("Transaction: ", transaction)
	return &transaction, nil
}

func trimJson(jsonText string) string {
	jsonText = strings.TrimSpace(jsonText)
	jsonText = strings.TrimPrefix(jsonText, "```json")
	jsonText = strings.TrimSuffix(jsonText, "```")
	jsonText = strings.TrimSpace(jsonText)
	return jsonText
}
