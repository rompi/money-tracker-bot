package errors_test

import (
	"fmt"
	apperrors "money-tracker-bot/internal/errors"
)

func ExampleNewTelegramError() {
	// Create a Network error (which is retryable)
	err := apperrors.NewNetworkError(
		"failed to send message",
		fmt.Errorf("connection timeout"),
	)

	// Add context information
	err.WithContext("user_id", "12345")
	err.WithContext("chat_id", "67890")
	err.WithComponent("message-sender")

	// Check if error is retryable
	if apperrors.IsRetryableError(err) {
		fmt.Println("This error can be retried")
	}

	// Output: This error can be retried
}

func ExampleNewSpreadsheetError() {
	// Create a Spreadsheet-specific error
	err := apperrors.NewSpreadsheetError(
		"failed to insert data",
		fmt.Errorf("quota exceeded"),
	)

	// Add operation context
	err.WithContext("spreadsheet_id", "abc123")
	err.WithContext("operation", "append_row")
	err.WithContext("row_data", map[string]string{
		"date":   "2025-01-20",
		"amount": "150000",
	})

	// Log the error
	apperrors.LogError(err)
}

func ExampleSafeExecute() {
	// Execute a potentially dangerous operation safely
	err := apperrors.SafeExecute(func() error {
		// This might panic
		someRiskyOperation()
		return nil
	}, "risky operation")

	if err != nil {
		fmt.Printf("Operation failed: %v", err)
	}
}

func ExampleHandleCriticalError() {
	// Handle a critical error that might require shutdown
	err := apperrors.NewConfigError(
		"missing required configuration",
		fmt.Errorf("DATABASE_URL not set"),
	)

	// This will log the error with full context and stack trace
	apperrors.HandleCriticalError(err, "application startup")

	// Application can decide whether to exit or continue with degraded functionality
}

func ExampleAppError_WithContext() {
	// Create an error and add rich context
	err := apperrors.NewGeminiError("AI processing failed", fmt.Errorf("rate limit exceeded"))

	// Add contextual information
	err.WithContext("user_id", "user_12345")
	err.WithContext("request_id", "req_abc123")
	err.WithContext("model", "gemini-2.0-flash")
	err.WithContext("prompt_length", 1024)

	// Set the component where error occurred
	err.WithComponent("ai-processor")

	// The error now contains rich context for debugging
	fmt.Printf("Error: %v", err)
	// Output: Error: [GEMINI_ERROR] AI processing failed: rate limit exceeded
}

func ExampleIsRetryableError() {
	// Check if different types of errors are retryable
	networkErr := apperrors.NewNetworkError("connection failed", nil)
	configErr := apperrors.NewConfigError("missing config", nil)

	fmt.Printf("Network error retryable: %v\n", apperrors.IsRetryableError(networkErr))
	fmt.Printf("Config error retryable: %v\n", apperrors.IsRetryableError(configErr))

	// Output:
	// Network error retryable: true
	// Config error retryable: false
}

func ExampleIsCriticalError() {
	// Check if errors are critical
	configErr := apperrors.NewConfigError("missing config", nil)
	telegramErr := apperrors.NewTelegramError("message failed", nil)

	fmt.Printf("Config error critical: %v\n", apperrors.IsCriticalError(configErr))
	fmt.Printf("Telegram error critical: %v\n", apperrors.IsCriticalError(telegramErr))

	// Output:
	// Config error critical: true
	// Telegram error critical: false
}

func someRiskyOperation() {
	// Simulate a function that might panic
	// This is just for example purposes
}
