# Error Handling Implementation Plan

## 🎯 Objective
Replace all `log.Fatal` and `log.Panic` calls with proper error handling to prevent application crashes and improve system reliability.

## 🔍 Current State Analysis

### Critical Issues Found
The application currently has **8 critical error locations** that crash the entire application:

1. **cmd/telebot/main.go:70** - `log.Fatal(err)` on bot startup failure
2. **internal/adapters/gemini/gemini.go:43** - `log.Fatal(err)` on AI generation failure
3. **internal/adapters/telegram/handler.go:32** - `log.Panic(err)` on bot creation failure
4. **internal/adapters/telegram/handler.go:60** - `log.Panic()` on type assertion failure
5. **internal/adapters/google/spreadsheet/client.go:32** - `log.Fatalf()` on Sheets client creation failure
6. **internal/adapters/google/spreadsheet/client.go:45** - `log.Fatalf()` on timezone loading failure
7. **internal/adapters/google/spreadsheet/client.go:65** - `log.Fatalf()` on data insertion failure
8. **internal/adapters/google/spreadsheet/client.go:108** - `log.Fatalf()` on data retrieval failure

### Impact Assessment
- **Severity**: CRITICAL - Application crashes completely on any error
- **User Experience**: Poor - No graceful degradation or error recovery
- **Operational Impact**: High - Requires manual restart after any failure
- **Development Impact**: Hard to test error scenarios

## 📋 Implementation Strategy

### Phase 1: Create Error Infrastructure (Day 1-2)
**Estimated Time**: 1.5 days

#### 1.1 Create Custom Error Types Package
**Location**: `internal/errors/`

```go
// internal/errors/errors.go
package errors

import (
    "fmt"
)

// AppError represents application-specific errors with context
type AppError struct {
    Code    string
    Message string
    Cause   error
    Context map[string]interface{}
}

func (e *AppError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

func (e *AppError) Unwrap() error {
    return e.Cause
}

// Predefined error codes
const (
    ErrCodeConfig         = "CONFIG_ERROR"
    ErrCodeTelegram       = "TELEGRAM_ERROR"
    ErrCodeGemini         = "GEMINI_ERROR"
    ErrCodeSpreadsheet    = "SPREADSHEET_ERROR"
    ErrCodeFileOperation  = "FILE_ERROR"
    ErrCodeValidation     = "VALIDATION_ERROR"
)

// Error constructors
func NewConfigError(message string, cause error) *AppError {
    return &AppError{
        Code:    ErrCodeConfig,
        Message: message,
        Cause:   cause,
        Context: make(map[string]interface{}),
    }
}

func NewTelegramError(message string, cause error) *AppError {
    return &AppError{
        Code:    ErrCodeTelegram,
        Message: message,
        Cause:   cause,
        Context: make(map[string]interface{}),
    }
}

func NewGeminiError(message string, cause error) *AppError {
    return &AppError{
        Code:    ErrCodeGemini,
        Message: message,
        Cause:   cause,
        Context: make(map[string]interface{}),
    }
}

func NewSpreadsheetError(message string, cause error) *AppError {
    return &AppError{
        Code:    ErrCodeSpreadsheet,
        Message: message,
        Cause:   cause,
        Context: make(map[string]interface{}),
    }
}
```

#### 1.2 Create Error Handling Utilities
```go
// internal/errors/handling.go
package errors

import (
    "log"
)

// HandleCriticalError logs critical errors and provides fallback behavior
func HandleCriticalError(err error, fallbackMsg string) {
    log.Printf("CRITICAL ERROR: %v - %s", err, fallbackMsg)
}

// IsRetryableError determines if an error should trigger a retry
func IsRetryableError(err error) bool {
    if appErr, ok := err.(*AppError); ok {
        return appErr.Code == ErrCodeSpreadsheet || appErr.Code == ErrCodeGemini
    }
    return false
}

// LogError logs errors with appropriate level based on severity
func LogError(err error) {
    if appErr, ok := err.(*AppError); ok {
        log.Printf("[%s] %s", appErr.Code, appErr.Error())
        return
    }
    log.Printf("ERROR: %v", err)
}
```

### Phase 2: Fix Google Spreadsheet Adapter (Day 2-3)
**Estimated Time**: 1 day
**Priority**: HIGHEST (4 fatal errors)

#### 2.1 Refactor NewSpreadsheetService
**File**: `internal/adapters/google/spreadsheet/client.go:32`

**Current Code**:
```go
func NewSpreadsheetService() *SpreadsheetService {
    srv, err := sheets.NewService(context.Background(), option.WithCredentialsFile("google-service-account.json"))
    if err != nil {
        log.Fatalf("Unable to retrieve Sheets client: %v", err)
    }
    return &SpreadsheetService{
        Sheet: srv,
    }
}
```

**Refactored Code**:
```go
func NewSpreadsheetService() (*SpreadsheetService, error) {
    srv, err := sheets.NewService(context.Background(), option.WithCredentialsFile("google-service-account.json"))
    if err != nil {
        return nil, errors.NewSpreadsheetError("unable to retrieve Sheets client", err)
    }
    return &SpreadsheetService{
        Sheet: srv,
    }, nil
}
```

#### 2.2 Refactor AppendRow Method
**File**: `internal/adapters/google/spreadsheet/client.go:45,65`

**Current Code**:
```go
func (s SpreadsheetService) AppendRow(ctx context.Context, spreadsheetId string, trx transaction_domain.Transaction) CategorySummary {
    loc, err := time.LoadLocation("Asia/Bangkok")
    if err != nil {
        log.Fatalf("Unable to load Asia/Bangkok timezone: %v", err)
    }
    // ... rest of method
    _, err = s.Sheet.Spreadsheets.Values.Append(spreadsheetId, "detailed!A:G", values).ValueInputOption("USER_ENTERED").Do()
    if err != nil {
        log.Fatalf("Unable to insert data to sheet: %v", err)
    }
}
```

**Refactored Code**:
```go
func (s SpreadsheetService) AppendRow(ctx context.Context, spreadsheetId string, trx transaction_domain.Transaction) (CategorySummary, error) {
    loc, err := time.LoadLocation("Asia/Bangkok")
    if err != nil {
        return CategorySummary{}, errors.NewSpreadsheetError("unable to load timezone", err)
    }

    // ... create values ...

    _, err = s.Sheet.Spreadsheets.Values.Append(spreadsheetId, "detailed!A:G", values).ValueInputOption("USER_ENTERED").Do()
    if err != nil {
        return CategorySummary{}, errors.NewSpreadsheetError("unable to insert data to sheet", err)
    }

    // ... fetch summary with error handling ...
    summary, err := s.fetchCategorySummary(ctx, spreadsheetId, trx.Category)
    if err != nil {
        // Return partial success - data was inserted but summary fetch failed
        errors.LogError(err)
        return CategorySummary{Category: trx.Category}, nil
    }

    return summary, nil
}

// Helper method for summary fetching
func (s SpreadsheetService) fetchCategorySummary(ctx context.Context, spreadsheetId, category string) (CategorySummary, error) {
    summaryRange := "summary!A2:F12"
    summaryValues, err := s.Sheet.Spreadsheets.Values.Get(spreadsheetId, summaryRange).Do()
    if err != nil {
        return CategorySummary{}, errors.NewSpreadsheetError("unable to get data from summary sheet", err)
    }

    // ... existing logic to find category summary ...

    return result, nil
}
```

#### 2.3 Update GetCellValue Method
**File**: `internal/adapters/google/spreadsheet/client.go:108`

**Current Code**:
```go
func (s SpreadsheetService) GetCellValue(ctx context.Context, spreadsheetId string) {
    values, err := s.Sheet.Spreadsheets.Values.Get(spreadsheetId, "Sheet1!A2:E7").Do()
    if err != nil {
        log.Fatalf("Unable to Get data from sheet: %v", err)
    }
    for _, value := range values.Values {
        fmt.Println(value)
    }
}
```

**Refactored Code**:
```go
func (s SpreadsheetService) GetCellValue(ctx context.Context, spreadsheetId string) error {
    values, err := s.Sheet.Spreadsheets.Values.Get(spreadsheetId, "Sheet1!A2:E7").Do()
    if err != nil {
        return errors.NewSpreadsheetError("unable to get data from sheet", err)
    }
    for _, value := range values.Values {
        fmt.Println(value)
    }
    return nil
}
```

### Phase 3: Fix Telegram Adapter (Day 3-4)
**Estimated Time**: 1 day
**Priority**: HIGH (2 panic errors)

#### 3.1 Refactor NewTelegramHandler
**File**: `internal/adapters/telegram/handler.go:32`

**Current Code**:
```go
func NewTelegramHandler(token string, transactionService transactions.ITransaction) *TelegramHandler {
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        log.Panic(err)
    }
    return &TelegramHandler{
        Telebot:            bot,
        TransactionService: transactionService,
    }
}
```

**Refactored Code**:
```go
func NewTelegramHandler(token string, transactionService transactions.ITransaction) (*TelegramHandler, error) {
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return nil, errors.NewTelegramError("failed to create Telegram bot", err)
    }
    return &TelegramHandler{
        Telebot:            bot,
        TransactionService: transactionService,
    }, nil
}
```

#### 3.2 Refactor Start Method
**File**: `internal/adapters/telegram/handler.go:60`

**Current Code**:
```go
func (t *TelegramHandler) Start() {
    realBot, ok := t.Telebot.(*tgbotapi.BotAPI)
    if !ok {
        log.Panic("Telebot is not a *tgbotapi.BotAPI")
    }
    // ... rest of method
}
```

**Refactored Code**:
```go
func (t *TelegramHandler) Start() error {
    realBot, ok := t.Telebot.(*tgbotapi.BotAPI)
    if !ok {
        return errors.NewTelegramError("telebot is not a *tgbotapi.BotAPI", nil)
    }

    realBot.Debug = true
    log.Printf("Authorized on account %s", realBot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := realBot.GetUpdatesChan(u)

    for update := range updates {
        // Wrap message handling in error recovery
        if err := t.handleUpdate(update); err != nil {
            errors.LogError(err)
            // Continue processing other messages instead of crashing
        }
    }
    return nil
}

// New helper method for safe update handling
func (t *TelegramHandler) handleUpdate(update tgbotapi.Update) error {
    if update.Message == nil {
        return nil
    }

    if update.Message.IsCommand() {
        return t.handleCommand(update.Message)
    }

    if update.Message.Document != nil {
        return t.handleDocumentSafe(update.Message)
    } else if update.Message.Photo != nil {
        return t.handlePhotoSafe(update.Message)
    } else {
        return t.handleMessageSafe(update.Message)
    }
}
```

### Phase 4: Fix Gemini Adapter (Day 4)
**Estimated Time**: 0.5 days
**Priority**: MEDIUM (1 fatal error)

#### 4.1 Refactor GenerateContent Method
**File**: `internal/adapters/gemini/gemini.go:43`

**Current Code**:
```go
func (c *GeminiClient) GenerateContent(ctx context.Context, prompt string) {
    _, err := c.Model.GenerateContent(ctx, genai.Text(prompt))
    if err != nil {
        log.Fatal(err)
    }
}
```

**Refactored Code**:
```go
func (c *GeminiClient) GenerateContent(ctx context.Context, prompt string) error {
    _, err := c.Model.GenerateContent(ctx, genai.Text(prompt))
    if err != nil {
        return errors.NewGeminiError("failed to generate content", err)
    }
    return nil
}
```

### Phase 5: Fix Main Application (Day 4-5)
**Estimated Time**: 0.5 days
**Priority**: CRITICAL (1 fatal error)

#### 5.1 Refactor Main Function
**File**: `cmd/telebot/main.go:70`

**Current Code**:
```go
func main() {
    if err := startBot(); err != nil {
        log.Fatal(err)
    }
}
```

**Refactored Code**:
```go
func main() {
    if err := startBot(); err != nil {
        errors.HandleCriticalError(err, "Bot failed to start - check configuration and dependencies")
        os.Exit(1) // Controlled exit instead of fatal crash
    }
}
```

#### 5.2 Update startBot and Dependencies
Update all service initialization calls to handle the new error returns:

```go
func startBot() error {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found or failed to load, proceeding with system env")
    }

    telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
    apiKey := os.Getenv("GEMINI_API_KEY")

    if testBotDeps.Override {
        return startBotWithDeps(telegramToken, apiKey, testBotDeps.SpreadsheetService, testBotDeps.GeminiClient)
    }

    googleSpreadsheet, err := spreadsheet.NewSpreadsheetService()
    if err != nil {
        return fmt.Errorf("failed to create spreadsheet service: %w", err)
    }

    geminiClient := gemini.NewClient(apiKey)
    return startBotWithDeps(telegramToken, apiKey, googleSpreadsheet, geminiClient)
}

func startBotWithDeps(telegramToken, apiKey string, spreadsheetService SpreadsheetService, geminiClient GeminiClient) error {
    if telegramToken == "" {
        return ErrEnvVarMissing("TELEGRAM_BOT_TOKEN")
    }
    if apiKey == "" {
        return ErrEnvVarMissing("GEMINI_API_KEY")
    }

    if s, ok := spreadsheetService.(*spreadsheet.SpreadsheetService); ok {
        if g, ok := geminiClient.(*gemini.GeminiClient); ok {
            transactionService := transactions.NewTransactionService(g, s)
            telegramHandler, err := telegram.NewTelegramHandler(telegramToken, transactionService)
            if err != nil {
                return fmt.Errorf("failed to create telegram handler: %w", err)
            }

            log.Println("Telegram bot started")
            return telegramHandler.Start()
        }
    }
    return nil
}
```

## 🧪 Testing Strategy

### Phase 6: Add Error Handling Tests (Day 5)
**Estimated Time**: 1 day

#### 6.1 Unit Tests for Error Types
```go
// internal/errors/errors_test.go
func TestAppError(t *testing.T) {
    cause := fmt.Errorf("original error")
    err := NewSpreadsheetError("test message", cause)

    assert.Equal(t, ErrCodeSpreadsheet, err.Code)
    assert.Equal(t, "test message", err.Message)
    assert.Equal(t, cause, err.Cause)
    assert.Contains(t, err.Error(), "test message")
    assert.Contains(t, err.Error(), "original error")
}
```

#### 6.2 Integration Tests for Error Scenarios
```go
// Test that services handle errors gracefully
func TestSpreadsheetService_AppendRow_HandlesErrors(t *testing.T) {
    // Test with invalid credentials
    // Test with network errors
    // Test with invalid spreadsheet ID
}

func TestTelegramHandler_HandlesErrors(t *testing.T) {
    // Test with invalid token
    // Test with network failures
    // Test message processing errors
}
```

#### 6.3 Error Recovery Tests
```go
func TestErrorRecovery(t *testing.T) {
    // Test that application continues after recoverable errors
    // Test that critical errors are logged properly
    // Test fallback behaviors
}
```

## ✅ Success Criteria

### Immediate Goals (End of Implementation)
- [ ] Zero `log.Fatal` or `log.Panic` calls in codebase
- [ ] All service constructors return errors instead of panicking
- [ ] All operations have proper error propagation
- [ ] Application continues running after non-critical errors
- [ ] Comprehensive error logging with context

### Quality Gates
- [ ] All existing tests still pass
- [ ] New error handling tests achieve >90% coverage
- [ ] Manual testing shows graceful error handling
- [ ] No application crashes during error scenarios

### Operational Improvements
- [ ] Errors are logged with context for debugging
- [ ] Application provides meaningful error messages
- [ ] System can recover from transient failures
- [ ] Monitoring can detect and alert on error patterns

## 📊 Risk Assessment & Mitigation

### High Risk Areas
1. **Spreadsheet Operations**: Most critical for data persistence
   - **Mitigation**: Implement retry logic and partial success handling
2. **Service Initialization**: Affects application startup
   - **Mitigation**: Validate configuration before service creation
3. **Message Processing**: Core functionality
   - **Mitigation**: Isolate error handling per message

### Rollback Plan
- Keep original implementations commented out until testing is complete
- Use feature flags to toggle between old and new error handling
- Maintain backward compatibility in service interfaces

## 📈 Expected Benefits

### Immediate Benefits
- **Zero Downtime**: Application won't crash on errors
- **Better Debugging**: Structured error messages with context
- **Improved Reliability**: Graceful degradation instead of crashes

### Long-term Benefits
- **Easier Maintenance**: Clear error patterns and handling
- **Better User Experience**: Partial functionality during issues
- **Operational Excellence**: Proper monitoring and alerting

---

**Total Estimated Time**: 5 days
**Priority**: CRITICAL
**Dependencies**: None (can start immediately)
**Impact**: High - Eliminates all application crashes