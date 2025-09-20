# Money Tracker Bot - Development Documentation

## Project Overview
The Money Tracker Bot is a Telegram bot that helps users track their financial transactions using Google Spreadsheets and AI-powered transaction analysis. This documentation provides comprehensive details about the project's architecture, recent changes, and development context for AI-assisted development.

## Recent Changes Summary

### Phase 1: Error Infrastructure Implementation (Latest)
**Commit**: `e8bc1c7 feat: implement Phase 1 error infrastructure`
**Date**: Latest commit
**Status**: ✅ Complete

A comprehensive error handling infrastructure has been implemented to replace direct `log.Fatal` and `log.Panic` calls throughout the application with structured error handling that enables graceful degradation and recovery.

## 🏗️ Architecture Overview

### Current Architecture
```
cmd/telebot/                    # Main application entry point
├── main.go                     # Bootstrap and dependency injection

internal/
├── errors/                     # ✅ NEW: Centralized error handling
├── adapters/                   # External service adapters
│   ├── telegram/              # Telegram bot API
│   ├── google/spreadsheet/    # Google Sheets integration
│   └── gemini/                # AI service integration
├── service/transactions/      # Business logic layer
├── port/out/ai/              # AI interface definitions
├── domain/transactions/       # Domain models
└── common/                   # Shared utilities

plan/                          # Documentation
├── refactoring-plan.md       # Comprehensive refactoring roadmap
└── CONTRIBUTOR.md            # Development guidelines
```

## 🆕 Phase 1: Error Infrastructure Implementation

### New Error Handling Package (`internal/errors/`)

#### Core Files Added:
- **`doc.go`**: Package documentation and usage patterns
- **`codes.go`**: Error code constants and severity levels
- **`errors.go`**: Core error types and constructors
- **`handling.go`**: Error handling, logging, and recovery utilities
- **`README.md`**: Comprehensive usage guide and migration examples

#### Test Files Added:
- **`errors_test.go`**: Core error functionality tests
- **`handling_test.go`**: Error handling and logging tests
- **`examples_test.go`**: Usage examples and documentation tests
- **`integration_test.go`**: End-to-end error handling scenarios

### Key Features Implemented

#### 1. Structured Error Types
```go
type AppError struct {
    Code      string                 // Error categorization
    Message   string                 // Human-readable message
    Cause     error                  // Underlying error (if any)
    Context   map[string]interface{} // Additional debugging context
    Severity  Severity               // INFO, WARNING, ERROR, CRITICAL
    Timestamp time.Time              // When the error occurred
    Component string                 // Which component generated the error
}
```

#### 2. Domain-Specific Error Constructors
- **Configuration**: `NewConfigError()` - Critical startup issues
- **Telegram**: `NewTelegramError()`, `NewTelegramCriticalError()` - Bot API issues
- **Gemini AI**: `NewGeminiError()`, `NewGeminiTimeoutError()` - AI service issues
- **Spreadsheet**: `NewSpreadsheetError()`, `NewSpreadsheetCriticalError()` - Google Sheets issues
- **Validation**: `NewValidationError()` - Input validation issues
- **Transaction**: `NewTransactionError()` - Business logic issues
- **Network**: `NewNetworkError()` - Connectivity issues
- **File**: `NewFileError()` - File operation issues
- **Data**: `NewDataAccessError()` - Data persistence issues

#### 3. Error Classification System

**Error Codes:**
- `CONFIG_ERROR` - Configuration and startup issues (Critical, Non-retryable)
- `TELEGRAM_ERROR` - Telegram API issues (Error, Non-retryable)
- `GEMINI_ERROR` - AI service issues (Error, Retryable)
- `SPREADSHEET_ERROR` - Google Sheets issues (Error, Retryable)
- `FILE_ERROR` - File operation issues (Error, Non-retryable)
- `VALIDATION_ERROR` - Input validation issues (Error, Non-retryable)
- `NETWORK_ERROR` - Network connectivity issues (Warning, Retryable)
- `TIMEOUT_ERROR` - Operation timeout issues (Warning, Retryable)
- `TRANSACTION_ERROR` - Business logic issues (Error, Non-retryable)
- `DATA_ACCESS_ERROR` - Data persistence issues (Error, Non-retryable)

**Severity Levels:**
- `INFO` - Informational messages
- `WARNING` - Warning conditions that don't affect operation
- `ERROR` - Error conditions that affect specific operations
- `CRITICAL` - Critical conditions that may require shutdown

#### 4. Panic Recovery and Safe Execution
```go
// Execute risky operations with automatic panic recovery
err := errors.SafeExecute(func() error {
    return riskyOperation()
}, "operation context")

// Manual panic recovery
func someFunction() (err error) {
    defer func() {
        if recovered := errors.RecoverFromPanic(); recovered != nil {
            err = recovered
        }
    }()
    // ... risky code ...
}
```

#### 5. Context-Rich Error Handling
```go
err := errors.NewTelegramError("message send failed", cause)
err.WithContext("user_id", "12345")
err.WithContext("message_type", "photo")
err.WithComponent("message-handler")

errors.HandleError(err, "processing user command")
```

#### 6. Intelligent Retry Logic
```go
if err := operation(); err != nil {
    if errors.IsRetryableError(err) {
        // Implement retry logic for network, timeout, and service errors
        time.Sleep(time.Second)
        return operation()
    }
    return err
}
```

#### 7. Graceful Degradation
```go
// Handle critical errors without crashing
if err := initializeService(); err != nil {
    errors.HandleCriticalError(err, "service initialization")

    if errors.IsCriticalError(err) {
        // Decide: exit gracefully or continue with reduced functionality
        errors.ExitGracefully(err, 1)
    }
}
```

### Migration Strategy

The error infrastructure was designed for gradual adoption:

1. **✅ Phase 1**: Add error package (Complete)
2. **Phase 2**: Replace `log.Fatal` calls in adapters
3. **Phase 3**: Replace `log.Panic` calls in handlers
4. **Phase 4**: Update service constructors to return errors
5. **Phase 5**: Add retry logic and graceful degradation

### Test Coverage
- **Unit Tests**: Core functionality, error creation, context management
- **Integration Tests**: End-to-end error handling scenarios
- **Mock Support**: Testable logging with dependency injection
- **Examples**: Comprehensive usage examples for documentation

## 🔍 Code Locations Requiring Migration

Based on the refactoring plan, these locations currently use `log.Fatal` or `log.Panic` and need migration:

### Critical `log.Fatal` Calls:
- `internal/adapters/google/spreadsheet/client.go:32,45,65,108`
- `internal/adapters/telegram/handler.go:32,60`
- `internal/adapters/gemini/gemini.go:43`

### `context.TODO()` Usage:
- `internal/adapters/telegram/handler.go:155,189`

## 📋 Development Guidelines

### Error Handling Best Practices

1. **Use Specific Error Types**: Choose the most appropriate constructor
   ```go
   // Good
   return errors.NewTelegramError("failed to send message", err)

   // Avoid
   return fmt.Errorf("telegram error: %w", err)
   ```

2. **Add Meaningful Context**:
   ```go
   err := errors.NewGeminiError("analysis failed", cause)
   err.WithContext("transaction_text", text)
   err.WithContext("user_id", userID)
   err.WithComponent("transaction-processor")
   ```

3. **Handle Errors Gracefully**:
   ```go
   // Replace log.Fatal
   if err != nil {
       return errors.NewConfigError("service initialization failed", err)
   }

   // Replace log.Panic
   if bot == nil {
       return errors.NewTelegramCriticalError("bot is nil", nil)
   }
   ```

4. **Use Retry Logic for Retryable Errors**:
   ```go
   if errors.IsRetryableError(err) {
       time.Sleep(backoffDuration)
       return retryOperation()
   }
   ```

5. **Log Appropriately**:
   ```go
   errors.HandleError(err, "processing user message")           // Normal errors
   errors.HandleCriticalError(err, "application startup")      // Critical errors
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

## 🧪 Testing

### Running Error Package Tests
```bash
# Run all tests
go test ./internal/errors -v

# Run with coverage
go test ./internal/errors -cover

# Run with race detection
go test ./internal/errors -race

# Static analysis
go vet ./internal/errors
```

### Mock Logger Usage
```go
// Test error logging
mockLogger := &MockLogger{}
errors.SetLogger(mockLogger)
defer errors.SetLogger(errors.DefaultLogger{})

// ... perform operations ...

// Verify logged messages
assert.NotEmpty(t, mockLogger.Messages)
```

## 🚀 Next Development Steps

### Immediate Priorities (Phase 2)
1. **Replace `log.Fatal` in Google Spreadsheet Client** (`internal/adapters/google/spreadsheet/client.go`)
2. **Replace `log.Fatal` in Telegram Handler** (`internal/adapters/telegram/handler.go`)
3. **Replace `log.Fatal` in Gemini Client** (`internal/adapters/gemini/gemini.go`)

### Medium-term Goals (Phase 3-4)
1. **Update service constructors** to return errors instead of using `log.Fatal`
2. **Replace `context.TODO()`** with proper context propagation
3. **Add timeout and cancellation** support for AI operations

### Long-term Goals (Phase 5)
1. **Implement retry logic** for transient failures
2. **Add circuit breaker patterns** for external services
3. **Implement graceful degradation** modes
4. **Add metrics and observability** for error tracking

## Agent Instructions
**IMPORTANT**: Agents should ignore all files under the `plan/` directory unless explicitly instructed by the user to refer to those files. The planning documents are for reference only and should not influence regular development tasks.

## Project Structure & Module Organization
This Go module targets Go 1.23 and follows a hexagonal architecture. The executable entry point lives in `cmd/telebot`, while shared logic sits in `internal/` with adapters (Telegram, Gemini, Google Sheets), domain models, services, and ports kept isolated. Reference material is in `docs/`, and configuration templates live in `.env.example` and `google-service-account.json`.

## Build, Test, and Development Commands
- `make run` (or `go run ./cmd/telebot`) starts the Telegram bot; export `TELEGRAM_BOT_TOKEN`, `GEMINI_API_KEY`, and Sheets credentials first.
- `make build` emits a Linux AMD64 binary at `./bot`.
- `make fmt` and `make lint` wrap `go fmt ./...` and `go vet ./...`.
- `make test` runs `go test -cover ./...`.

## Coding Guidelines (Go)
Format code with `gofmt` (or `make fmt`) and organise imports consistently. Use tabs for indentation, `MixedCaps` for exports, and return wrapped errors with `%w`. Align package layout with the hexagonal layers: domain contracts in `internal/domain`, ports as interfaces under `internal/port`, and adapters implementing those contracts. Prefer dependency injection via constructors and avoid reading environment variables outside `cmd/`.

## Architecture Overview
Hexagonal boundaries must stay intact: the domain layer owns business rules and stays free of external SDKs; ports declare what the domain expects; adapters wrap Telegram, Gemini, and Sheets clients. New integrations should add a port interface, place the adapter under `internal/adapters/<system>`, and wire dependencies in `cmd/telebot/main.go`. Keep DTO translations within adapters to shield the domain from API drift.

## Testing Guidelines
Unit and integration tests reside alongside code as `*_test.go`. Mirror the structure of the code under test (e.g., tests for `internal/service/transactions` sit in the same folder). New features need `go test ./...` to pass without reducing coverage. Prefer table-driven tests and concise fixtures.

## Commit & Pull Request Guidelines
Current history favours prefixing subjects with the feature area, such as `Feature/shopping quota` or `Fix/...`, optionally referencing issues (`(#ID)`). Keep commits focused and written in imperative mood. Pull requests should summarise behaviour changes, cite relevant issues, include configuration or migration notes, and attach test evidence (`make test`) or bot run logs when altering runtime behaviour.

## Security & Configuration Tips
Store all secrets in a local `.env`; do not commit filled `.env` or service account files. When testing Google integrations, prefer stubbed clients from `internal/adapters/google` and redact spreadsheet IDs in logs. Rotate tokens regularly before sharing builds.

## 📚 References

### Documentation Files
- `internal/errors/README.md` - Complete error package guide
- `plan/refactoring-plan.md` - Comprehensive refactoring roadmap
- `plan/CONTRIBUTOR.md` - Development guidelines

### Key Components
- `internal/errors/` - Complete error handling infrastructure
- `cmd/telebot/main.go` - Application entry point and dependency injection
- Individual adapter AI.md files in each component directory

### External Dependencies
- **Telegram Bot API** - User interaction
- **Google Sheets API** - Data persistence
- **Google Gemini AI** - Transaction analysis
- **Go Testing Framework** - Comprehensive test suite

## 💡 Development Context for AI

### Current State
- ✅ **Error Infrastructure**: Complete and tested error handling system
- ✅ **Documentation**: Comprehensive guides and examples
- ✅ **Testing**: Full test coverage with mocks
- ⏳ **Migration**: Ready to begin replacing `log.Fatal` calls

### When Working on This Codebase
1. **Always use error constructors** from `internal/errors` package
2. **Add context** to errors for better debugging
3. **Check retryability** before implementing retry logic
4. **Use structured logging** through error handlers
5. **Write tests** for both success and error scenarios
6. **Follow the migration phases** outlined in the refactoring plan

### Common Patterns to Recognize
- Service initialization with dependency injection
- Error handling with structured logging
- Context propagation for cancellation
- Retry logic for external service calls
- Graceful degradation for non-critical failures

This documentation will be continuously updated as development progresses through the remaining phases of the refactoring plan.
