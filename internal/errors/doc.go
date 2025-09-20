// Package errors provides application-specific error types and handling utilities
// for the Money Tracker Bot application.
//
// # Overview
//
// This package replaces direct usage of log.Fatal and log.Panic throughout the
// application with structured error handling that enables graceful degradation
// and recovery.
//
// # Key Features
//
// - Structured error types with codes, context, and severity levels
// - Domain-specific error constructors (Telegram, Gemini, Spreadsheet)
// - Retry logic support with retryable error classification
// - Panic recovery and conversion to structured errors
// - Comprehensive logging with context information
// - Testing utilities and mock implementations
//
// # Usage Patterns
//
// Replace log.Fatal calls:
//
//	// Before:
//	if err != nil {
//	    log.Fatalf("Failed to create service: %v", err)
//	}
//
//	// After:
//	if err != nil {
//	    return nil, errors.NewConfigError("failed to create service", err)
//	}
//
// Add context to errors:
//
//	err := errors.NewTelegramError("message send failed", cause)
//	err.WithContext("user_id", userID)
//	err.WithContext("message_type", "photo")
//
// Handle errors gracefully:
//
//	if err := operation(); err != nil {
//	    if errors.IsRetryableError(err) {
//	        // Implement retry logic
//	    } else {
//	        errors.HandleError(err, "operation context")
//	    }
//	}
//
// # Error Codes
//
// The package defines standard error codes for categorization:
//   - CONFIG_ERROR: Configuration and startup issues
//   - TELEGRAM_ERROR: Telegram API issues
//   - GEMINI_ERROR: AI service issues
//   - SPREADSHEET_ERROR: Google Sheets issues
//   - FILE_ERROR: File operation issues
//   - VALIDATION_ERROR: Input validation issues
//   - NETWORK_ERROR: Network connectivity issues
//   - TIMEOUT_ERROR: Operation timeout issues
//
// # Severity Levels
//
// Errors are classified by severity:
//   - INFO: Informational messages
//   - WARNING: Warning conditions that don't affect operation
//   - ERROR: Error conditions that affect specific operations
//   - CRITICAL: Critical conditions that may require shutdown
package errors
