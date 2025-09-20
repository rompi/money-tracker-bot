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
func SafeExecute(fn func() error, context string) (returnErr error) {
	defer func() {
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
			appErr.WithContext("function_context", context)

			HandleError(appErr, context)
			returnErr = appErr
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
