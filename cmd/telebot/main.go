package main

import (
	"log"
	"os"
	"rompi/gobot/internal/adapters/gemini"
	"rompi/gobot/internal/adapters/google/spreadsheet"
	"rompi/gobot/internal/adapters/telegram"
	aiport "rompi/gobot/internal/port/out/ai"
	"rompi/gobot/internal/service/transactions"

	"github.com/joho/godotenv"
)

type Config struct {
	AiModel         aiport.AiPort
	TelegramHandler telegram.StoredFile
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load, proceeding with system env")
	}

	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	googleSpreadsheet := spreadsheet.NewSpreadsheetService()

	apiKey := os.Getenv("GEMINI_API_KEY")
	geminiClient := gemini.NewClient(apiKey)

	transactionService := transactions.NewTransactionService(geminiClient, googleSpreadsheet)
	telegramHandler := telegram.NewTelegramHandler(telegramToken, transactionService)
	log.Println("Telegram bot started")
	telegramHandler.Start()

}
