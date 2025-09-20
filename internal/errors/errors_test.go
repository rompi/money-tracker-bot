package errors

import (
	"fmt"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appError *AppError
		expected string
	}{
		{
			name: "error with cause",
			appError: &AppError{
				Code:    ErrCodeConfig,
				Message: "configuration failed",
				Cause:   fmt.Errorf("file not found"),
			},
			expected: "[CONFIG_ERROR] configuration failed: file not found",
		},
		{
			name: "error without cause",
			appError: &AppError{
				Code:    ErrCodeTelegram,
				Message: "connection failed",
				Cause:   nil,
			},
			expected: "[TELEGRAM_ERROR] connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.appError.Error()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("underlying error")
	appErr := NewConfigError("test error", cause)

	if appErr.Unwrap() != cause {
		t.Errorf("expected unwrap to return cause error")
	}
}

func TestAppError_WithContext(t *testing.T) {
	appErr := NewTelegramError("test error", nil)

	appErr.WithContext("user_id", "12345")
	appErr.WithContext("chat_id", "67890")

	if appErr.Context["user_id"] != "12345" {
		t.Errorf("expected context user_id to be '12345'")
	}

	if appErr.Context["chat_id"] != "67890" {
		t.Errorf("expected context chat_id to be '67890'")
	}
}

func TestAppError_WithComponent(t *testing.T) {
	appErr := NewGeminiError("test error", nil)
	appErr.WithComponent("ai-processor")

	if appErr.Component != "ai-processor" {
		t.Errorf("expected component to be 'ai-processor', got '%s'", appErr.Component)
	}
}

func TestAppError_IsRetryable(t *testing.T) {
	tests := []struct {
		name        string
		errorCode   string
		shouldRetry bool
	}{
		{"network error", ErrCodeNetwork, true},
		{"timeout error", ErrCodeTimeout, true},
		{"spreadsheet error", ErrCodeSpreadsheet, true},
		{"gemini error", ErrCodeGemini, true},
		{"config error", ErrCodeConfig, false},
		{"validation error", ErrCodeValidation, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr := &AppError{Code: tt.errorCode}
			if appErr.IsRetryable() != tt.shouldRetry {
				t.Errorf("expected IsRetryable() to return %v for %s", tt.shouldRetry, tt.errorCode)
			}
		})
	}
}

func TestAppError_IsCritical(t *testing.T) {
	tests := []struct {
		name       string
		severity   Severity
		isCritical bool
	}{
		{"critical error", SeverityCritical, true},
		{"error", SeverityError, false},
		{"warning", SeverityWarning, false},
		{"info", SeverityInfo, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr := &AppError{Severity: tt.severity}
			if appErr.IsCritical() != tt.isCritical {
				t.Errorf("expected IsCritical() to return %v for %s", tt.isCritical, tt.name)
			}
		})
	}
}

func TestErrorConstructors(t *testing.T) {
	cause := fmt.Errorf("underlying error")

	tests := []struct {
		name             string
		constructor      func(string, error) *AppError
		expectedCode     string
		expectedSeverity Severity
		expectedComponent string
	}{
		{"NewConfigError", NewConfigError, ErrCodeConfig, SeverityCritical, "config"},
		{"NewTelegramError", NewTelegramError, ErrCodeTelegram, SeverityError, "telegram"},
		{"NewGeminiError", NewGeminiError, ErrCodeGemini, SeverityError, "gemini"},
		{"NewSpreadsheetError", NewSpreadsheetError, ErrCodeSpreadsheet, SeverityError, "spreadsheet"},
		{"NewFileError", NewFileError, ErrCodeFileOperation, SeverityError, "file"},
		{"NewValidationError", NewValidationError, ErrCodeValidation, SeverityError, "validation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr := tt.constructor("test message", cause)

			if appErr.Code != tt.expectedCode {
				t.Errorf("expected code %s, got %s", tt.expectedCode, appErr.Code)
			}

			if appErr.Severity != tt.expectedSeverity {
				t.Errorf("expected severity %v, got %v", tt.expectedSeverity, appErr.Severity)
			}

			if appErr.Component != tt.expectedComponent {
				t.Errorf("expected component %s, got %s", tt.expectedComponent, appErr.Component)
			}

			if appErr.Message != "test message" {
				t.Errorf("expected message 'test message', got '%s'", appErr.Message)
			}

			if appErr.Cause != cause {
				t.Errorf("expected cause to be set")
			}

			if appErr.Timestamp.IsZero() {
				t.Errorf("expected timestamp to be set")
			}
		})
	}
}

func TestSeverity_String(t *testing.T) {
	tests := []struct {
		severity Severity
		expected string
	}{
		{SeverityInfo, "INFO"},
		{SeverityWarning, "WARNING"},
		{SeverityError, "ERROR"},
		{SeverityCritical, "CRITICAL"},
		{Severity(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.severity.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}