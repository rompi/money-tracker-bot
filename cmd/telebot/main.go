package main

import (
	"log"
	"os"
	"rompi/gobot/internal/adapters/gemini"
	"rompi/gobot/internal/adapters/google/spreadsheet"
	"rompi/gobot/internal/adapters/telegram"
	"rompi/gobot/internal/service/transactions"

	"github.com/joho/godotenv"
)

// startBot loads configuration, initializes services, and starts the Telegram bot.
// startBotWithDeps allows dependency injection for easier testing.

// Dependency interfaces for testability
type SpreadsheetService interface{}
type GeminiClient interface{}

func startBotWithDeps(telegramToken, apiKey string, spreadsheetService SpreadsheetService, geminiClient GeminiClient) error {
	if telegramToken == "" {
		return ErrEnvVarMissing("TELEGRAM_BOT_TOKEN")
	}
	if apiKey == "" {
		return ErrEnvVarMissing("GEMINI_API_KEY")
	}
	// Only run the real bot if using real implementations
	if s, ok := spreadsheetService.(*spreadsheet.SpreadsheetService); ok {
		if g, ok := geminiClient.(*gemini.GeminiClient); ok {
			transactionService := transactions.NewTransactionService(g, s)
			telegramHandler := telegram.NewTelegramHandler(telegramToken, transactionService)
			log.Println("Telegram bot started")
			telegramHandler.Start()
		}
	}
	return nil
}

func startBot() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load, proceeding with system env")
	}

	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	apiKey := os.Getenv("GEMINI_API_KEY")
	googleSpreadsheet := spreadsheet.NewSpreadsheetService()
	geminiClient := gemini.NewClient(apiKey)
	return startBotWithDeps(telegramToken, apiKey, googleSpreadsheet, geminiClient)
}

// ErrEnvVarMissing is returned when a required environment variable is missing.
type ErrEnvVarMissing string

func (e ErrEnvVarMissing) Error() string {
	return "required environment variable not set: " + string(e)
}

func main() {
	if err := startBot(); err != nil {
		log.Fatal(err)
	}
}
