# Error Handling Package

This package provides structured error handling for the Money Tracker Bot application, replacing direct usage of `log.Fatal` and `log.Panic` with graceful error handling and recovery.

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

## Features

- **Structured Error Types**: Rich error objects with codes, context, and severity
- **Domain-Specific Constructors**: Specialized error types for different components
- **Context Support**: Add contextual information for better debugging
- **Retry Logic**: Automatic classification of retryable vs non-retryable errors
- **Panic Recovery**: Safe execution with automatic panic recovery
- **Graceful Logging**: Structured logging with appropriate severity levels

## Error Types

### Error Codes

| Code | Description | Retryable | Typical Severity |
|------|-------------|-----------|------------------|
| `CONFIG_ERROR` | Configuration and startup issues | ❌ | Critical |
| `TELEGRAM_ERROR` | Telegram API issues | ❌ | Error |
| `GEMINI_ERROR` | AI service issues | ✅ | Error |
| `SPREADSHEET_ERROR` | Google Sheets issues | ✅ | Error |
| `FILE_ERROR` | File operation issues | ❌ | Error |
| `VALIDATION_ERROR` | Input validation issues | ❌ | Error |
| `NETWORK_ERROR` | Network connectivity issues | ✅ | Warning |
| `TIMEOUT_ERROR` | Operation timeout issues | ✅ | Warning |

### Severity Levels

- **INFO**: Informational messages
- **WARNING**: Warning conditions that don't affect operation
- **ERROR**: Error conditions that affect specific operations
- **CRITICAL**: Critical conditions that may require shutdown

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

### Service Constructor Pattern

```go
// Before
func NewService() *Service {
    client, err := createClient()
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    return &Service{client: client}
}

// After
func NewService() (*Service, error) {
    client, err := createClient()
    if err != nil {
        return nil, errors.NewConfigError("failed to create client", err)
    }
    return &Service{client: client}, nil
}
```

## Usage Examples

### Basic Error Creation

```go
// Create a simple error
err := errors.NewTelegramError("message send failed", cause)

// Add context for debugging
err.WithContext("user_id", "12345")
err.WithContext("message_type", "photo")
err.WithComponent("message-handler")
```

### Error Handling

```go
// Handle errors with context
if err := sendMessage(); err != nil {
    errors.HandleError(err, "processing user command")

    // Check if we should retry
    if errors.IsRetryableError(err) {
        // Implement retry logic
        time.Sleep(time.Second)
        return sendMessage()
    }

    return err
}
```

### Safe Execution

```go
// Execute risky operations safely
err := errors.SafeExecute(func() error {
    // This operation might panic
    return processImage(imagePath)
}, "image processing")

if err != nil {
    // Handle the error (panic was converted to error)
    errors.LogError(err)
}
```

### Critical Error Handling

```go
// Handle critical errors
if err := initializeDatabase(); err != nil {
    errors.HandleCriticalError(err, "application startup")

    if errors.IsCriticalError(err) {
        // Decide whether to exit or continue with degraded functionality
        errors.ExitGracefully(err, 1)
    }
}
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./internal/errors -v

# Run with coverage
go test ./internal/errors -cover

# Run race detection
go test ./internal/errors -race

# Run static analysis
go vet ./internal/errors
```

### Mock Logger for Testing

```go
// Use mock logger in tests
mockLogger := &MockLogger{}
errors.SetLogger(mockLogger)
defer errors.SetLogger(errors.DefaultLogger{})

// ... perform operations that log errors ...

// Verify log messages
if len(mockLogger.Messages) == 0 {
    t.Error("expected error to be logged")
}
```

## Best Practices

1. **Use Specific Error Types**: Choose the most appropriate error constructor for your domain
2. **Add Context**: Always add relevant context information for debugging
3. **Handle Gracefully**: Avoid letting errors crash the application
4. **Check Retryability**: Use `IsRetryableError()` to implement smart retry logic
5. **Log Appropriately**: Use `HandleError()` for operational errors, `HandleCriticalError()` for critical issues
6. **Test Error Scenarios**: Write tests for both success and failure cases

## Integration with Existing Code

This package is designed to be gradually adopted:

1. **Phase 1**: Add the error package (✅ Complete)
2. **Phase 2**: Replace `log.Fatal` calls in adapters
3. **Phase 3**: Replace `log.Panic` calls in handlers
4. **Phase 4**: Update service constructors to return errors
5. **Phase 5**: Add retry logic and graceful degradation

Each phase can be implemented independently without breaking existing functionality.

## Performance Considerations

- Error creation is lightweight with minimal allocations
- Context maps are created lazily (only when needed)
- Stack traces are only generated for critical errors
- Logging is performed asynchronously where possible

## Contributing

When adding new error types:

1. Add the error code constant to `codes.go`
2. Add the constructor function to `errors.go`
3. Add appropriate tests to `errors_test.go`
4. Update this README with the new error type
5. Add usage examples if needed