package main

import (
	"log"
	"money-tracker-bot/internal/adapters/gemini"
	"money-tracker-bot/internal/adapters/google/spreadsheet"
	"money-tracker-bot/internal/adapters/telegram"
	"money-tracker-bot/internal/errors"
	"money-tracker-bot/internal/service/transactions"
	"os"

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
			telegramHandler, err := telegram.NewTelegramHandler(telegramToken, transactionService)
			if err != nil {
				return err
			}
			log.Println("Telegram bot started")
			if err := telegramHandler.Start(); err != nil {
				return err
			}
		}
	}
	return nil
}

var testBotDeps struct {
       SpreadsheetService SpreadsheetService
       GeminiClient GeminiClient
       Override bool
}

func startBot() error {
       if err := godotenv.Load(); err != nil {
	       log.Println("No .env file found or failed to load, proceeding with system env")
       }

       telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
       apiKey := os.Getenv("GEMINI_API_KEY")
       if testBotDeps.Override {
	       return startBotWithDeps(telegramToken, apiKey, testBotDeps.SpreadsheetService, testBotDeps.GeminiClient)
       }
       googleSpreadsheet, err := spreadsheet.NewSpreadsheetService()
       if err != nil {
	       return err
       }
       geminiClient, err := gemini.NewClient(apiKey)
       if err != nil {
	       return err
       }
       return startBotWithDeps(telegramToken, apiKey, googleSpreadsheet, geminiClient)
}

// ErrEnvVarMissing is returned when a required environment variable is missing.
type ErrEnvVarMissing string

func (e ErrEnvVarMissing) Error() string {
	return "required environment variable not set: " + string(e)
}

func main() {
	if err := startBot(); err != nil {
		// Use structured error handling instead of log.Fatal
		appErr, ok := err.(*errors.AppError)
		if ok {
			errors.HandleCriticalError(appErr, "application startup")
			if errors.IsCriticalError(appErr) {
				errors.ExitGracefully(appErr, 1)
			}
		} else {
			// Handle non-AppError types
			criticalErr := errors.NewConfigError("application startup failed", err).
				WithComponent("main")
			errors.HandleCriticalError(criticalErr, "application startup")
			errors.ExitGracefully(criticalErr, 1)
		}
	}
}
