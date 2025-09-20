package errors

import (
	"fmt"
	"strings"
	"testing"
)

// Mock logger for testing
type MockLogger struct {
	Messages []string
}

func (m *MockLogger) Printf(format string, v ...interface{}) {
	m.Messages = append(m.Messages, fmt.Sprintf(format, v...))
}

func (m *MockLogger) Println(v ...interface{}) {
	m.Messages = append(m.Messages, fmt.Sprint(v...))
}

func (m *MockLogger) Reset() {
	m.Messages = nil
}

func TestHandleError(t *testing.T) {
	mockLogger := &MockLogger{}
	SetLogger(mockLogger)
	defer SetLogger(DefaultLogger{})

	err := NewTelegramError("connection failed", fmt.Errorf("network timeout"))
	err.WithContext("user_id", "12345")

	HandleError(err, "processing user message")

	if len(mockLogger.Messages) == 0 {
		t.Error("expected error to be logged")
	}

	message := mockLogger.Messages[0]
	if !strings.Contains(message, "TELEGRAM_ERROR") {
		t.Errorf("expected log message to contain error code, got: %s", message)
	}

	if !strings.Contains(message, "connection failed") {
		t.Errorf("expected log message to contain error message, got: %s", message)
	}
}

func TestHandleCriticalError(t *testing.T) {
	mockLogger := &MockLogger{}
	SetLogger(mockLogger)
	defer SetLogger(DefaultLogger{})

	err := NewConfigError("missing required configuration", fmt.Errorf("file not found"))

	HandleCriticalError(err, "application startup")

	if len(mockLogger.Messages) < 2 {
		t.Error("expected critical error to generate multiple log messages")
	}

	// Check that critical error is logged
	found := false
	for _, msg := range mockLogger.Messages {
		if strings.Contains(msg, "CRITICAL ERROR") {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected critical error message to be logged")
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "retryable network error",
			err:      NewNetworkError("connection failed", nil),
			expected: true,
		},
		{
			name:     "non-retryable config error",
			err:      NewConfigError("missing config", nil),
			expected: false,
		},
		{
			name:     "generic error",
			err:      fmt.Errorf("generic error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryableError(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsCriticalError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "critical config error",
			err:      NewConfigError("missing config", nil),
			expected: true,
		},
		{
			name:     "non-critical telegram error",
			err:      NewTelegramError("message failed", nil),
			expected: false,
		},
		{
			name:     "generic error",
			err:      fmt.Errorf("generic error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsCriticalError(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRecoverFromPanic(t *testing.T) {
	// Test that RecoverFromPanic returns nil when no panic occurred
	err := RecoverFromPanic()
	if err != nil {
		t.Errorf("expected nil when no panic, got: %v", err)
	}

	// Test will be covered by TestSafeExecute which actually triggers a panic
}

func TestSafeExecute(t *testing.T) {
	mockLogger := &MockLogger{}
	SetLogger(mockLogger)
	defer SetLogger(DefaultLogger{})

	// Test successful execution
	err := SafeExecute(func() error {
		return nil
	}, "test operation")

	if err != nil {
		t.Errorf("expected no error for successful execution, got: %v", err)
	}

	// Test execution with panic
	SafeExecute(func() error {
		panic("test panic")
	}, "test operation with panic")

	// Check that panic was logged
	found := false
	for _, msg := range mockLogger.Messages {
		if strings.Contains(msg, "panic") {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected panic to be logged")
	}
}

func TestToAppError(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected string
	}{
		{
			name:     "AppError input",
			input:    NewTelegramError("test error", nil),
			expected: "TELEGRAM_ERROR",
		},
		{
			name:     "generic error input",
			input:    fmt.Errorf("generic error"),
			expected: "GENERIC_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toAppError(tt.input)
			if result.Code != tt.expected {
				t.Errorf("expected code %s, got %s", tt.expected, result.Code)
			}
		})
	}
}

func TestFormatContext(t *testing.T) {
	tests := []struct {
		name     string
		context  map[string]interface{}
		expected string
	}{
		{
			name:     "empty context",
			context:  map[string]interface{}{},
			expected: "",
		},
		{
			name: "single context item",
			context: map[string]interface{}{
				"user_id": "12345",
			},
			expected: "user_id=12345",
		},
		{
			name: "multiple context items",
			context: map[string]interface{}{
				"user_id": "12345",
				"action":  "send_message",
			},
			// Note: map iteration order is not guaranteed, so we check both contain the items
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatContext(tt.context)
			if tt.name == "empty context" && result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			} else if tt.name == "single context item" && result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			} else if tt.name == "multiple context items" {
				// For multiple items, just check that both values are present
				if !strings.Contains(result, "user_id=12345") || !strings.Contains(result, "action=send_message") {
					t.Errorf("expected result to contain both context items, got: %q", result)
				}
			}
		})
	}
}
