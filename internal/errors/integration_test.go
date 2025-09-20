package errors_test

import (
	"errors"
	"fmt"
	apperrors "money-tracker-bot/internal/errors"
	"testing"
)

func TestErrorPackageIntegration(t *testing.T) {
	// Test that error package integrates well with existing patterns

	// Simulate service creation failure
	err := simulateServiceCreation()
	if err == nil {
		t.Error("expected service creation to fail for test")
	}

	// Verify it's an AppError
	appErr, ok := err.(*apperrors.AppError)
	if !ok {
		t.Errorf("expected AppError, got %T", err)
	}

	// Verify error properties
	if appErr.Code != apperrors.ErrCodeSpreadsheet {
		t.Errorf("expected spreadsheet error code, got %s", appErr.Code)
	}

	if !appErr.IsRetryable() {
		t.Error("expected spreadsheet error to be retryable")
	}
}

func simulateServiceCreation() error {
	// Simulate the pattern we'll use in real services
	return apperrors.NewSpreadsheetError(
		"failed to create spreadsheet service",
		fmt.Errorf("credentials file not found"),
	)
}

func TestErrorChaining(t *testing.T) {
	// Test error unwrapping works correctly
	rootCause := fmt.Errorf("network connection failed")
	appErr := apperrors.NewNetworkError("operation failed", rootCause)

	// Test unwrap
	if appErr.Unwrap() != rootCause {
		t.Error("error unwrapping failed")
	}

	// Test error chain traversal
	if !errors.Is(appErr, rootCause) {
		t.Error("error chain traversal failed")
	}
}

func TestContextualErrorHandling(t *testing.T) {
	// Test pattern that will be used in real application
	err := simulateTransactionProcessing("user123", "photo_upload")

	if err == nil {
		t.Error("expected transaction processing to fail for test")
	}

	// Verify context was preserved
	appErr, ok := err.(*apperrors.AppError)
	if !ok {
		t.Errorf("expected AppError, got %T", err)
	}

	if appErr.Context["user_id"] != "user123" {
		t.Error("expected user_id context to be preserved")
	}

	if appErr.Context["operation"] != "photo_upload" {
		t.Error("expected operation context to be preserved")
	}
}

func simulateTransactionProcessing(userID, operation string) error {
	// Simulate the pattern for adding context to errors
	err := apperrors.NewGeminiError("AI processing failed", fmt.Errorf("quota exceeded"))
	err.WithContext("user_id", userID)
	err.WithContext("operation", operation)
	err.WithContext("timestamp", "2025-01-20T10:30:00Z")

	return err
}

func TestServiceIntegrationPattern(t *testing.T) {
	// Test the pattern we'll use when refactoring services
	spreadsheetService, err := createMockSpreadsheetService()
	if err != nil {
		// This demonstrates the pattern: service creation returns error instead of panicking
		if !apperrors.IsRetryableError(err) {
			t.Error("expected spreadsheet service creation error to be retryable")
		}
	}

	if spreadsheetService != nil {
		t.Error("expected service creation to fail for test")
	}
}

func createMockSpreadsheetService() (interface{}, error) {
	// Simulate the new pattern for service creation
	return nil, apperrors.NewSpreadsheetError(
		"unable to retrieve Sheets client",
		fmt.Errorf("credentials not found"),
	)
}

func TestErrorRecoveryPattern(t *testing.T) {
	// Test the pattern for recovering from operations that might panic
	var finalErr error

	finalErr = apperrors.SafeExecute(func() error {
		// Simulate operation that might panic
		simulateRiskyOperation()
		return nil
	}, "risky operation")

	// The function should not panic, and error should be handled gracefully
	if finalErr != nil {
		t.Logf("Operation failed gracefully: %v", finalErr)
	}
}

func simulateRiskyOperation() {
	// This would normally cause a panic, but SafeExecute should handle it
	panic("simulated panic for testing")
}