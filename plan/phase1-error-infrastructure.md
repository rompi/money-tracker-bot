# Phase 1: Error Infrastructure Implementation Plan

## 🎯 Objective
Create robust error infrastructure foundation that will support all subsequent error handling improvements across the application.

## ⏱️ Timeline
**Estimated Duration**: 1.5 days (12 hours)
**Priority**: CRITICAL
**Dependencies**: None - can start immediately

## 📋 Task Breakdown

### Task 1.1: Project Structure Setup (30 minutes)
**Estimated Time**: 0.5 hours

#### Create Error Package Directory
```bash
mkdir -p internal/errors
```

#### Verify Package Structure
```
internal/
├── errors/          # NEW - Error handling infrastructure
│   ├── errors.go    # Core error types and constructors
│   ├── handling.go  # Error handling utilities
│   └── codes.go     # Error code constants
├── adapters/
├── domain/
├── service/
└── ...
```

#### Initial Package Documentation
Create `internal/errors/doc.go`:
```go
// Package errors provides application-specific error types and handling utilities.
//
// This package defines structured error types that replace log.Fatal and log.Panic
// calls throughout the application, enabling graceful error handling and recovery.
//
// Key Components:
// - AppError: Core error type with code, message, and context
// - Error constructors for different domains (Telegram, Gemini, Spreadsheet)
// - Utilities for error classification and handling
package errors
```

### Task 1.2: Core Error Types Implementation (2 hours)
**Estimated Time**: 2 hours

#### File: `internal/errors/codes.go`
```go
package errors

// Error codes for categorizing different types of application errors
const (
    // Configuration and startup errors
    ErrCodeConfig = "CONFIG_ERROR"

    // External service errors
    ErrCodeTelegram    = "TELEGRAM_ERROR"
    ErrCodeGemini      = "GEMINI_ERROR"
    ErrCodeSpreadsheet = "SPREADSHEET_ERROR"

    // Internal operation errors
    ErrCodeFileOperation = "FILE_ERROR"
    ErrCodeValidation    = "VALIDATION_ERROR"
    ErrCodeTransaction   = "TRANSACTION_ERROR"

    // Network and connectivity errors
    ErrCodeNetwork = "NETWORK_ERROR"
    ErrCodeTimeout = "TIMEOUT_ERROR"

    // Data and persistence errors
    ErrCodeDataAccess    = "DATA_ACCESS_ERROR"
    ErrCodeDataFormat    = "DATA_FORMAT_ERROR"
    ErrCodeDataIntegrity = "DATA_INTEGRITY_ERROR"
)

// Error severity levels
type Severity int

const (
    SeverityInfo Severity = iota
    SeverityWarning
    SeverityError
    SeverityCritical
)

func (s Severity) String() string {
    switch s {
    case SeverityInfo:
        return "INFO"
    case SeverityWarning:
        return "WARNING"
    case SeverityError:
        return "ERROR"
    case SeverityCritical:
        return "CRITICAL"
    default:
        return "UNKNOWN"
    }
}
```

#### File: `internal/errors/errors.go`
```go
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
    Severity  Severity              `json:"severity"`
    Timestamp time.Time             `json:"timestamp"`
    Component string                `json:"component"`
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
```

### Task 1.3: Error Handling Utilities (2 hours)
**Estimated Time**: 2 hours

#### File: `internal/errors/handling.go`
```go
package errors

import (
    "fmt"
    "log"
    "os"
    "runtime"
    "strings"
)

// Logger interface for dependency injection
type Logger interface {
    Printf(format string, v ...interface{})
    Println(v ...interface{})
}

// DefaultLogger uses the standard log package
type DefaultLogger struct{}

func (d DefaultLogger) Printf(format string, v ...interface{}) {
    log.Printf(format, v...)
}

func (d DefaultLogger) Println(v ...interface{}) {
    log.Println(v...)
}

var (
    // Global logger instance (can be replaced for testing)
    logger Logger = DefaultLogger{}
)

// SetLogger allows replacing the default logger
func SetLogger(l Logger) {
    logger = l
}

// HandleCriticalError handles critical errors that may require application shutdown
func HandleCriticalError(err error, context string) {
    appErr := toAppError(err)

    // Log the critical error with full context
    logErrorWithContext(appErr, context, true)

    // If it's truly critical, we may need to exit
    if appErr.IsCritical() {
        logger.Printf("CRITICAL ERROR: Application may need to shutdown - %s", context)
        // Don't exit immediately - let the caller decide
    }
}

// HandleError handles non-critical errors with appropriate logging
func HandleError(err error, context string) {
    if err == nil {
        return
    }

    appErr := toAppError(err)
    logErrorWithContext(appErr, context, false)
}

// LogError logs an error with appropriate formatting
func LogError(err error) {
    if err == nil {
        return
    }

    appErr := toAppError(err)
    logErrorWithContext(appErr, "", false)
}

// logErrorWithContext logs an error with full context information
func logErrorWithContext(err *AppError, context string, includeStackTrace bool) {
    // Build log message
    var parts []string

    // Add severity
    parts = append(parts, fmt.Sprintf("[%s]", err.Severity.String()))

    // Add component
    if err.Component != "" {
        parts = append(parts, fmt.Sprintf("[%s]", err.Component))
    }

    // Add code
    parts = append(parts, fmt.Sprintf("[%s]", err.Code))

    // Add message
    parts = append(parts, err.Message)

    // Add context if provided
    if context != "" {
        parts = append(parts, fmt.Sprintf("Context: %s", context))
    }

    // Add error context
    if len(err.Context) > 0 {
        contextStr := formatContext(err.Context)
        parts = append(parts, fmt.Sprintf("Details: %s", contextStr))
    }

    // Add underlying error
    if err.Cause != nil {
        parts = append(parts, fmt.Sprintf("Cause: %v", err.Cause))
    }

    // Log the formatted message
    message := strings.Join(parts, " | ")
    logger.Println(message)

    // Add stack trace for critical errors
    if includeStackTrace || err.IsCritical() {
        logger.Printf("Stack trace: %s", getStackTrace())
    }
}

// formatContext converts context map to readable string
func formatContext(context map[string]interface{}) string {
    if len(context) == 0 {
        return ""
    }

    var parts []string
    for key, value := range context {
        parts = append(parts, fmt.Sprintf("%s=%v", key, value))
    }
    return strings.Join(parts, ", ")
}

// getStackTrace returns a formatted stack trace
func getStackTrace() string {
    buf := make([]byte, 1024)
    n := runtime.Stack(buf, false)
    return string(buf[:n])
}

// toAppError converts any error to AppError
func toAppError(err error) *AppError {
    if appErr, ok := err.(*AppError); ok {
        return appErr
    }

    // Create generic AppError for non-AppError types
    return &AppError{
        Code:      "GENERIC_ERROR",
        Message:   err.Error(),
        Cause:     err,
        Context:   make(map[string]interface{}),
        Severity:  SeverityError,
        Component: "unknown",
    }
}

// IsRetryableError determines if an error should trigger a retry
func IsRetryableError(err error) bool {
    if appErr, ok := err.(*AppError); ok {
        return appErr.IsRetryable()
    }
    return false
}

// IsCriticalError determines if an error is critical
func IsCriticalError(err error) bool {
    if appErr, ok := err.(*AppError); ok {
        return appErr.IsCritical()
    }
    return false
}

// RecoverFromPanic recovers from panics and converts them to errors
func RecoverFromPanic() error {
    if r := recover(); r != nil {
        var err error
        switch x := r.(type) {
        case string:
            err = fmt.Errorf("panic: %s", x)
        case error:
            err = fmt.Errorf("panic: %w", x)
        default:
            err = fmt.Errorf("panic: %v", x)
        }

        appErr := NewTransactionError("recovered from panic", err)
        appErr.WithContext("panic_value", r)
        appErr.WithContext("stack_trace", getStackTrace())

        return appErr
    }
    return nil
}

// SafeExecute executes a function with panic recovery
func SafeExecute(fn func() error, context string) error {
    defer func() {
        if err := RecoverFromPanic(); err != nil {
            HandleError(err, context)
        }
    }()

    return fn()
}

// ExitGracefully performs graceful shutdown
func ExitGracefully(err error, exitCode int) {
    if err != nil {
        HandleCriticalError(err, "Application shutdown")
    }

    // Perform any cleanup here
    logger.Printf("Application exiting with code %d", exitCode)
    os.Exit(exitCode)
}
```

### Task 1.4: Testing Infrastructure (3 hours)
**Estimated Time**: 3 hours

#### File: `internal/errors/errors_test.go`
```go
package errors

import (
    "fmt"
    "testing"
    "time"
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
        name        string
        constructor func(string, error) *AppError
        expectedCode string
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
```

#### File: `internal/errors/handling_test.go`
```go
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
    var recoveredErr error

    func() {
        defer func() {
            recoveredErr = RecoverFromPanic()
        }()
        panic("test panic")
    }()

    if recoveredErr == nil {
        t.Error("expected panic to be recovered as error")
    }

    appErr, ok := recoveredErr.(*AppError)
    if !ok {
        t.Error("expected recovered error to be AppError")
    }

    if appErr.Code != ErrCodeTransaction {
        t.Errorf("expected error code %s, got %s", ErrCodeTransaction, appErr.Code)
    }

    if !strings.Contains(appErr.Error(), "panic") {
        t.Errorf("expected error message to contain 'panic', got: %s", appErr.Error())
    }
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
```

### Task 1.5: Integration and Documentation (3 hours)
**Estimated Time**: 3 hours

#### Update go.mod if needed
Check if any new dependencies are required (none expected for this phase).

#### Create Usage Examples
File: `internal/errors/examples_test.go`
```go
package errors_test

import (
    "fmt"
    "money-tracker-bot/internal/errors"
)

func ExampleNewTelegramError() {
    // Create a Telegram-specific error
    err := errors.NewTelegramError(
        "failed to send message",
        fmt.Errorf("network timeout"),
    )

    // Add context information
    err.WithContext("user_id", "12345")
    err.WithContext("chat_id", "67890")
    err.WithComponent("message-sender")

    // Handle the error
    errors.HandleError(err, "processing user command")

    // Check if error is retryable
    if errors.IsRetryableError(err) {
        fmt.Println("This error can be retried")
    }

    // Output: This error can be retried
}

func ExampleNewSpreadsheetError() {
    // Create a Spreadsheet-specific error
    err := errors.NewSpreadsheetError(
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
    errors.LogError(err)
}

func ExampleSafeExecute() {
    // Execute a potentially dangerous operation safely
    err := errors.SafeExecute(func() error {
        // This might panic
        someRiskyOperation()
        return nil
    }, "risky operation")

    if err != nil {
        fmt.Printf("Operation failed: %v", err)
    }
}

func someRiskyOperation() {
    // Simulate a function that might panic
}
```

#### Update Package Documentation
Update `internal/errors/doc.go` with comprehensive documentation:
```go
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
//   // Before:
//   if err != nil {
//       log.Fatalf("Failed to create service: %v", err)
//   }
//
//   // After:
//   if err != nil {
//       return nil, errors.NewConfigError("failed to create service", err)
//   }
//
// Add context to errors:
//   err := errors.NewTelegramError("message send failed", cause)
//   err.WithContext("user_id", userID)
//   err.WithContext("message_type", "photo")
//
// Handle errors gracefully:
//   if err := operation(); err != nil {
//       if errors.IsRetryableError(err) {
//           // Implement retry logic
//       } else {
//           errors.HandleError(err, "operation context")
//       }
//   }
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
```

#### Create README for the errors package
File: `internal/errors/README.md`
```markdown
# Error Handling Package

This package provides structured error handling for the Money Tracker Bot application.

## Quick Start

```go
// Create domain-specific errors
err := errors.NewTelegramError("failed to send message", cause)
err.WithContext("user_id", "12345")

// Handle errors gracefully
errors.HandleError(err, "processing user message")

// Check error properties
if errors.IsRetryableError(err) {
    // Implement retry logic
}
```

## Migration Guide

### Replacing log.Fatal

```go
// Before
if err != nil {
    log.Fatalf("Service failed: %v", err)
}

// After
if err != nil {
    return errors.NewConfigError("service failed", err)
}
```

### Replacing log.Panic

```go
// Before
if bot == nil {
    log.Panic("Bot is nil")
}

// After
if bot == nil {
    return errors.NewTelegramError("bot is nil", nil)
}
```

## Testing

```bash
go test ./internal/errors -v
go test ./internal/errors -cover
```
```

### Task 1.6: Validation and Integration Testing (2 hours)
**Estimated Time**: 2 hours

#### Run All Tests
```bash
cd internal/errors
go test -v -cover
go test -race
go vet
```

#### Integration with Existing Code
Create a simple integration test to ensure the error package works with the existing codebase:

File: `internal/errors/integration_test.go`
```go
package errors_test

import (
    "fmt"
    "testing"
    "money-tracker-bot/internal/errors"
)

func TestErrorPackageIntegration(t *testing.T) {
    // Test that error package integrates well with existing patterns

    // Simulate service creation failure
    err := simulateServiceCreation()
    if err == nil {
        t.Error("expected service creation to fail for test")
    }

    // Verify it's an AppError
    appErr, ok := err.(*errors.AppError)
    if !ok {
        t.Errorf("expected AppError, got %T", err)
    }

    // Verify error properties
    if appErr.Code != errors.ErrCodeSpreadsheet {
        t.Errorf("expected spreadsheet error code, got %s", appErr.Code)
    }

    if !appErr.IsRetryable() {
        t.Error("expected spreadsheet error to be retryable")
    }
}

func simulateServiceCreation() error {
    // Simulate the pattern we'll use in real services
    return errors.NewSpreadsheetError(
        "failed to create spreadsheet service",
        fmt.Errorf("credentials file not found"),
    )
}

func TestErrorChaining(t *testing.T) {
    // Test error unwrapping works correctly
    rootCause := fmt.Errorf("network connection failed")
    appErr := errors.NewNetworkError("operation failed", rootCause)

    // Test unwrap
    if appErr.Unwrap() != rootCause {
        t.Error("error unwrapping failed")
    }

    // Test error chain traversal
    if !fmt.Is(appErr, rootCause) {
        t.Error("error chain traversal failed")
    }
}
```

## ✅ Deliverables Checklist

### Core Implementation
- [ ] `internal/errors/` directory created
- [ ] `codes.go` - Error codes and severity definitions
- [ ] `errors.go` - Core AppError type and constructors
- [ ] `handling.go` - Error handling utilities
- [ ] `doc.go` - Package documentation

### Testing
- [ ] `errors_test.go` - Unit tests for error types
- [ ] `handling_test.go` - Unit tests for handling utilities
- [ ] `examples_test.go` - Usage examples
- [ ] `integration_test.go` - Integration tests
- [ ] All tests pass with >90% coverage

### Documentation
- [ ] Package documentation complete
- [ ] README.md with migration guide
- [ ] Code examples and usage patterns
- [ ] Integration instructions for next phases

## 🧪 Testing Strategy

### Unit Tests
- Test all error constructors
- Test error properties (retryable, critical)
- Test context addition and formatting
- Test panic recovery
- Test logging functionality

### Integration Tests
- Test error package with existing patterns
- Test error unwrapping and chaining
- Test logger replacement
- Test with concurrent access

### Coverage Goals
- Aim for >90% test coverage
- Cover all error codes and constructors
- Test both success and failure scenarios
- Include edge cases and boundary conditions

## 📊 Success Metrics

### Immediate Validation
- [ ] All tests pass
- [ ] No breaking changes to existing code
- [ ] Error package can be imported and used
- [ ] Documentation is complete and accurate

### Quality Gates
- [ ] Code review completed
- [ ] Test coverage >90%
- [ ] No lint warnings or vet issues
- [ ] Integration tests demonstrate usage patterns

### Readiness for Phase 2
- [ ] Error infrastructure is stable
- [ ] Error constructors cover all domains
- [ ] Logging and handling utilities work correctly
- [ ] Ready to replace log.Fatal/log.Panic calls

---

**Phase 1 Total Time**: 1.5 days (12 hours)
**Critical Path**: Core error types → Testing → Documentation
**Dependencies**: None
**Risks**: Low - Pure additive changes
**Next Phase**: Ready to replace fatal errors in Google Spreadsheet adapter