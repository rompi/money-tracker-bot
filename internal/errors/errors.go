package errors

import (
	"fmt"
	"time"
)

// AppError represents application-specific errors with rich context
type AppError struct {
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Cause     error                  `json:"-"`
	Context   map[string]interface{} `json:"context"`
	Severity  Severity               `json:"severity"`
	Timestamp time.Time              `json:"timestamp"`
	Component string                 `json:"component"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap implements the unwrapper interface for error chains
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithContext adds context information to the error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithComponent sets the component where the error occurred
func (e *AppError) WithComponent(component string) *AppError {
	e.Component = component
	return e
}

// IsRetryable determines if the error indicates a retryable condition
func (e *AppError) IsRetryable() bool {
	switch e.Code {
	case ErrCodeNetwork, ErrCodeTimeout, ErrCodeSpreadsheet, ErrCodeGemini:
		return true
	default:
		return false
	}
}

// IsCritical determines if the error is critical for application operation
func (e *AppError) IsCritical() bool {
	return e.Severity == SeverityCritical
}

// newAppError creates a new AppError with basic fields
func newAppError(code, message, component string, severity Severity, cause error) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Cause:     cause,
		Context:   make(map[string]interface{}),
		Severity:  severity,
		Timestamp: time.Now(),
		Component: component,
	}
}

// Error constructors for different domains

// Configuration errors
func NewConfigError(message string, cause error) *AppError {
	return newAppError(ErrCodeConfig, message, "config", SeverityCritical, cause)
}

// Telegram service errors
func NewTelegramError(message string, cause error) *AppError {
	return newAppError(ErrCodeTelegram, message, "telegram", SeverityError, cause)
}

func NewTelegramCriticalError(message string, cause error) *AppError {
	return newAppError(ErrCodeTelegram, message, "telegram", SeverityCritical, cause)
}

// Gemini AI service errors
func NewGeminiError(message string, cause error) *AppError {
	return newAppError(ErrCodeGemini, message, "gemini", SeverityError, cause)
}

func NewGeminiTimeoutError(message string, cause error) *AppError {
	return newAppError(ErrCodeTimeout, message, "gemini", SeverityWarning, cause)
}

// Google Spreadsheet errors
func NewSpreadsheetError(message string, cause error) *AppError {
	return newAppError(ErrCodeSpreadsheet, message, "spreadsheet", SeverityError, cause)
}

func NewSpreadsheetCriticalError(message string, cause error) *AppError {
	return newAppError(ErrCodeSpreadsheet, message, "spreadsheet", SeverityCritical, cause)
}

// File operation errors
func NewFileError(message string, cause error) *AppError {
	return newAppError(ErrCodeFileOperation, message, "file", SeverityError, cause)
}

// Validation errors
func NewValidationError(message string, cause error) *AppError {
	return newAppError(ErrCodeValidation, message, "validation", SeverityError, cause)
}

// Transaction processing errors
func NewTransactionError(message string, cause error) *AppError {
	return newAppError(ErrCodeTransaction, message, "transaction", SeverityError, cause)
}

// Network errors
func NewNetworkError(message string, cause error) *AppError {
	return newAppError(ErrCodeNetwork, message, "network", SeverityWarning, cause)
}

// Data access errors
func NewDataAccessError(message string, cause error) *AppError {
	return newAppError(ErrCodeDataAccess, message, "data", SeverityError, cause)
}
