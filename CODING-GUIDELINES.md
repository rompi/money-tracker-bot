# Coding Guidelines - Money Tracker Bot

This document establishes coding standards and best practices for the Money Tracker Bot project to ensure code consistency, maintainability, and quality.

## 📋 Table of Contents

- [Go Language Standards](#go-language-standards)
- [Architecture Principles](#architecture-principles)
- [Error Handling](#error-handling)
- [Testing Standards](#testing-standards)
- [Documentation Requirements](#documentation-requirements)
- [Code Organization](#code-organization)
- [Security Guidelines](#security-guidelines)
- [Performance Considerations](#performance-considerations)
- [Git and Versioning](#git-and-versioning)

## 🔧 Go Language Standards

### Code Formatting
```go
// Use gofmt/goimports for consistent formatting
// Tabs for indentation, not spaces
// Line length: prefer 80-100 characters, max 120

// Good: Clear, readable function
func ProcessTransaction(ctx context.Context, text string) (*Transaction, error) {
    if text == "" {
        return nil, errors.NewValidationError("transaction text cannot be empty", nil)
    }

    // Implementation...
    return &Transaction{
        Amount:   amount,
        Category: category,
    }, nil
}

// Bad: Poor formatting and naming
func proc(t string)(*Transaction,error){
if t==""{return nil,fmt.Errorf("empty")}
// Implementation...
}
```

### Naming Conventions
```go
// Exported functions: MixedCaps
func ProcessTransaction() {}

// Unexported functions: mixedCaps
func parseAmount() {}

// Constants: ALL_CAPS or MixedCaps
const (
    MaxRetryAttempts = 3
    DefaultTimeout   = 30 * time.Second
)

// Interfaces: end with -er when appropriate
type TransactionProcessor interface {
    ProcessTransaction(ctx context.Context, data []byte) (*Transaction, error)
}

// Structs: MixedCaps
type TransactionService struct {
    aiClient     AIClient
    sheetClient  SpreadsheetClient
    logger       Logger
}

// Package names: lowercase, single word when possible
package transactions // Good
package transactionprocessing // Avoid
```

### Variable and Function Guidelines
```go
// Use descriptive names
userID := "12345"           // Good
id := "12345"              // Too generic
u := "12345"               // Too short

// Avoid stuttering
type TransactionService struct{}
func (ts *TransactionService) ProcessTransaction() {} // Good
func (ts *TransactionService) TransactionProcess() {} // Stuttering

// Use receiver names consistently (2-3 characters)
func (ts *TransactionService) Process() {}  // Good
func (service *TransactionService) Process() {} // Too verbose
func (t *TransactionService) Process() {}   // Potentially confusing

// Boolean variables should be questions
isValid := true     // Good
hasError := false   // Good
valid := true       // Less clear
```

### Import Organization
```go
import (
    // Standard library imports first
    "context"
    "fmt"
    "time"

    // Third-party imports second
    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "google.golang.org/api/sheets/v4"

    // Local imports last
    "github.com/rompi/money-tracker-bot/internal/domain/transactions"
    "github.com/rompi/money-tracker-bot/internal/errors"
)
```

## 🏗️ Architecture Principles

### Hexagonal Architecture Compliance
```go
// Domain layer: No external dependencies
package transactions

type Transaction struct {
    ID          string
    Amount      float64
    Category    string
    Timestamp   time.Time
}

// Port interface: Define what domain needs
type AIAnalyzer interface {
    AnalyzeTransaction(ctx context.Context, data []byte) (*Transaction, error)
}

// Adapter: Implements port interface
package gemini

type GeminiAdapter struct {
    client *genai.Client
}

func (g *GeminiAdapter) AnalyzeTransaction(ctx context.Context, data []byte) (*Transaction, error) {
    // External service implementation
}
```

### Dependency Injection
```go
// Constructor pattern for dependency injection
func NewTransactionService(
    aiClient AIAnalyzer,
    sheetClient SpreadsheetWriter,
    logger Logger,
) *TransactionService {
    return &TransactionService{
        aiClient:    aiClient,
        sheetClient: sheetClient,
        logger:      logger,
    }
}

// Avoid global variables
var globalClient *http.Client // Avoid this

// Use dependency injection instead
type HTTPService struct {
    client *http.Client
}
```

### Interface Design
```go
// Keep interfaces small and focused
type TransactionAnalyzer interface {
    AnalyzeTransaction(ctx context.Context, data []byte) (*Transaction, error)
}

// Avoid large interfaces
type EverythingInterface interface { // Avoid this approach
    AnalyzeTransaction()
    SaveTransaction()
    SendNotification()
    LogEvent()
}

// Prefer composition of small interfaces
type TransactionProcessor interface {
    TransactionAnalyzer
    TransactionSaver
}
```

## 🚨 Error Handling

### Use Custom Error Types
```go
// Always use error constructors from internal/errors package
func ProcessTransaction(text string) (*Transaction, error) {
    if text == "" {
        return nil, errors.NewValidationError("transaction text is required", nil)
    }

    result, err := aiClient.Analyze(text)
    if err != nil {
        return nil, errors.NewGeminiError("failed to analyze transaction", err).
            WithContext("input_text", text).
            WithComponent("transaction-processor")
    }

    return result, nil
}

// Never use log.Fatal or log.Panic
func BadExample() {
    if err != nil {
        log.Fatal("This will crash the application") // Never do this
    }
}

// Always return errors
func GoodExample() error {
    if err != nil {
        return errors.NewConfigError("initialization failed", err)
    }
    return nil
}
```

### Error Context and Wrapping
```go
// Add meaningful context to errors
err := errors.NewTelegramError("message send failed", originalErr)
err.WithContext("user_id", userID)
err.WithContext("message_type", "photo")
err.WithContext("file_size", fileSize)
err.WithComponent("telegram-handler")

// Use error wrapping for external errors
if err != nil {
    return fmt.Errorf("failed to process transaction: %w", err)
}
```

### Retry Logic
```go
// Check if errors are retryable
func ProcessWithRetry(ctx context.Context, data []byte) error {
    const maxRetries = 3
    var lastErr error

    for i := 0; i < maxRetries; i++ {
        err := processTransaction(ctx, data)
        if err == nil {
            return nil
        }

        if !errors.IsRetryableError(err) {
            return err // Don't retry non-retryable errors
        }

        lastErr = err
        time.Sleep(time.Second * time.Duration(i+1)) // Exponential backoff
    }

    return fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
```

## 🧪 Testing Standards

### Test Structure and Naming
```go
// Test files: *_test.go
// Test functions: TestFunctionName

func TestTransactionService_ProcessTransaction(t *testing.T) {
    // Table-driven tests for multiple scenarios
    tests := []struct {
        name          string
        input         string
        want          *Transaction
        wantErr       bool
        errorType     string
    }{
        {
            name:  "valid receipt text",
            input: "Coffee $4.50 at Starbucks",
            want: &Transaction{
                Amount:   4.50,
                Category: "Food",
                Merchant: "Starbucks",
            },
            wantErr: false,
        },
        {
            name:      "empty input",
            input:     "",
            want:      nil,
            wantErr:   true,
            errorType: "VALIDATION_ERROR",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            service := setupTestService(t)

            got, err := service.ProcessTransaction(context.Background(), tt.input)

            if tt.wantErr {
                assert.Error(t, err)
                if tt.errorType != "" {
                    assert.Contains(t, err.Error(), tt.errorType)
                }
                return
            }

            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### Test Coverage Requirements
```go
// Aim for 85%+ coverage on new code
// Test both success and error paths
func TestTransactionService_ProcessTransaction_ErrorHandling(t *testing.T) {
    tests := []struct {
        name          string
        setupMocks    func(*testing.T) (*MockAIClient, *MockSheetClient)
        input         string
        expectedError string
    }{
        {
            name: "AI service failure",
            setupMocks: func(t *testing.T) (*MockAIClient, *MockSheetClient) {
                aiClient := &MockAIClient{}
                aiClient.On("Analyze", mock.Anything, mock.Anything).
                    Return(nil, errors.New("AI service unavailable"))

                sheetClient := &MockSheetClient{}
                return aiClient, sheetClient
            },
            input:         "test transaction",
            expectedError: "GEMINI_ERROR",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            aiClient, sheetClient := tt.setupMocks(t)
            service := NewTransactionService(aiClient, sheetClient, &MockLogger{})

            _, err := service.ProcessTransaction(context.Background(), tt.input)

            assert.Error(t, err)
            assert.Contains(t, err.Error(), tt.expectedError)
        })
    }
}
```

### Mock Usage
```go
// Use interfaces for mockable dependencies
type MockAIClient struct {
    mock.Mock
}

func (m *MockAIClient) AnalyzeTransaction(ctx context.Context, data []byte) (*Transaction, error) {
    args := m.Called(ctx, data)
    return args.Get(0).(*Transaction), args.Error(1)
}

// Create test helpers for setup
func setupTestService(t *testing.T) *TransactionService {
    aiClient := &MockAIClient{}
    sheetClient := &MockSheetClient{}
    logger := &MockLogger{}

    return NewTransactionService(aiClient, sheetClient, logger)
}
```

## 📖 Documentation Requirements

### Function Documentation
```go
// All exported functions must have comprehensive docstrings
// Include purpose, parameters, return values, and examples

// ProcessTransaction analyzes transaction text and extracts structured financial data.
// It uses AI to identify transaction components like amount, merchant, and category.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - text: Raw transaction text from receipt or user input
//
// Returns:
//   - *Transaction: Structured transaction data with amount, category, merchant
//   - error: Validation, processing, or external service errors
//
// Example Input: "Coffee $4.50 at Starbucks 2025-01-20"
// Example Output: Transaction{Amount: 4.50, Category: "Food", Merchant: "Starbucks", Date: "2025-01-20"}
//
// Errors:
//   - Returns VALIDATION_ERROR for empty or invalid input
//   - Returns GEMINI_ERROR for AI analysis failures
//   - Returns NETWORK_ERROR for connectivity issues (retryable)
func ProcessTransaction(ctx context.Context, text string) (*Transaction, error) {
    // Implementation...
}

// Internal functions should have brief comments
// parseAmount extracts monetary amount from transaction text
func parseAmount(text string) (float64, error) {
    // Implementation...
}
```

### Package Documentation
```go
// Package doc.go files for each package
// internal/service/transactions/doc.go

// Package transactions provides core business logic for processing financial transactions.
//
// This package implements the transaction processing workflow:
//   1. Input validation and sanitization
//   2. AI-powered transaction analysis
//   3. Data structuring and categorization
//   4. Integration with external storage systems
//
// The package follows hexagonal architecture principles, depending only on
// interfaces defined in the port packages for external integrations.
//
// Example usage:
//
//     service := transactions.NewService(aiClient, sheetClient, logger)
//     transaction, err := service.ProcessTransaction(ctx, "Coffee $4.50")
//     if err != nil {
//         // Handle error
//     }
//     // Use transaction data
//
// Error handling:
//
// The package uses structured errors from the internal/errors package.
// All errors include context and are categorized by type and severity.
// Retryable errors (network, timeout) can be identified using errors.IsRetryableError().
package transactions
```

### AI.md Updates
```go
// Update package AI.md when adding new features
// Example: internal/service/transactions/AI.md

## Recent Changes
- Added ProcessTransaction function with AI integration
- Implemented retry logic for external service failures
- Added comprehensive error handling with context

## Key Functions
- ProcessTransaction: Main entry point for transaction processing
- validateInput: Input sanitization and validation
- categorizeTransaction: AI-powered transaction categorization

## Integration Points
- Gemini AI: Transaction text analysis
- Google Sheets: Data persistence
- Telegram: User interaction and file handling

## Error Patterns
- Use NewValidationError for input validation failures
- Use NewGeminiError for AI service issues
- Use NewSpreadsheetError for storage failures
```

## 📁 Code Organization

### File Structure
```
internal/
├── domain/
│   └── transactions/
│       ├── doc.go           # Package documentation
│       ├── transaction.go   # Domain entities
│       └── transaction_test.go
├── service/
│   └── transactions/
│       ├── doc.go
│       ├── service.go       # Business logic
│       ├── service_test.go
│       └── AI.md           # AI development context
├── adapters/
│   ├── telegram/
│   │   ├── doc.go
│   │   ├── handler.go      # Telegram bot handlers
│   │   ├── handler_test.go
│   │   └── AI.md
│   └── gemini/
│       ├── doc.go
│       ├── client.go       # Gemini AI integration
│       ├── client_test.go
│       └── AI.md
└── port/
    └── out/
        └── ai/
            └── analyzer.go  # Interface definitions
```

### Import Paths
```go
// Use full import paths
import "github.com/rompi/money-tracker-bot/internal/domain/transactions"

// Group imports logically
import (
    // Standard library
    "context"
    "fmt"

    // External dependencies
    "github.com/stretchr/testify/assert"

    // Internal packages
    "github.com/rompi/money-tracker-bot/internal/errors"
)
```

## 🔒 Security Guidelines

### Secrets Management
```go
// Never hardcode secrets
const apiKey = "sk-1234567890" // Never do this

// Use environment variables
func NewGeminiClient() (*Client, error) {
    apiKey := os.Getenv("GEMINI_API_KEY")
    if apiKey == "" {
        return nil, errors.NewConfigError("GEMINI_API_KEY is required", nil)
    }

    return &Client{apiKey: apiKey}, nil
}

// Validate environment variables at startup
func validateConfig() error {
    required := []string{
        "TELEGRAM_BOT_TOKEN",
        "GEMINI_API_KEY",
        "SPREADSHEET_ID",
    }

    for _, env := range required {
        if os.Getenv(env) == "" {
            return errors.NewConfigError(fmt.Sprintf("%s is required", env), nil)
        }
    }

    return nil
}
```

### Input Validation
```go
// Always validate external input
func ProcessTransaction(ctx context.Context, text string) (*Transaction, error) {
    // Sanitize input
    text = strings.TrimSpace(text)
    if text == "" {
        return nil, errors.NewValidationError("transaction text cannot be empty", nil)
    }

    // Validate length
    if len(text) > 10000 {
        return nil, errors.NewValidationError("transaction text too long", nil)
    }

    // Additional validation...
    return processValidatedInput(ctx, text)
}
```

### Logging Security
```go
// Never log sensitive information
func processPayment(userID string, amount float64, cardNumber string) error {
    // Good: Log non-sensitive information
    logger.Info("Processing payment",
        "user_id", userID,
        "amount", amount,
        "timestamp", time.Now(),
    )

    // Bad: Never log sensitive data
    // logger.Info("Payment details", "card_number", cardNumber)

    // Use redaction for logs
    logger.Info("Payment method", "card_ending", cardNumber[len(cardNumber)-4:])

    return nil
}
```

## ⚡ Performance Considerations

### Context Usage
```go
// Always pass context for cancellation and timeouts
func ProcessTransaction(ctx context.Context, data []byte) (*Transaction, error) {
    // Set timeout for AI processing
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    // Pass context to all downstream calls
    result, err := aiClient.Analyze(ctx, data)
    if err != nil {
        return nil, err
    }

    return result, nil
}

// Don't use context.TODO() in production code
func BadExample() {
    // Avoid this
    result, _ := aiClient.Analyze(context.TODO(), data)
}
```

### Resource Management
```go
// Always close resources
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close() // Always defer close

    // Process file...
    return nil
}

// Use sync.Pool for frequently allocated objects
var transactionPool = sync.Pool{
    New: func() interface{} {
        return &Transaction{}
    },
}

func getTransaction() *Transaction {
    return transactionPool.Get().(*Transaction)
}

func putTransaction(t *Transaction) {
    // Reset transaction
    *t = Transaction{}
    transactionPool.Put(t)
}
```

### Concurrency
```go
// Use goroutines appropriately
func processMultipleTransactions(transactions []string) error {
    var wg sync.WaitGroup
    errCh := make(chan error, len(transactions))

    for _, tx := range transactions {
        wg.Add(1)
        go func(transaction string) {
            defer wg.Done()
            if err := processTransaction(transaction); err != nil {
                errCh <- err
            }
        }(tx)
    }

    wg.Wait()
    close(errCh)

    // Check for errors
    for err := range errCh {
        if err != nil {
            return err
        }
    }

    return nil
}
```

## 📝 Git and Versioning

### Commit Messages
```bash
# Format: <type>: <description>
#
# <body>
#
# <footer>

feat: add transaction category detection

Implement AI-based category classification for transactions.
The system now automatically categorizes expenses into predefined
categories like Food, Transport, Shopping, etc.

- Add CategoryClassifier interface
- Implement Gemini-based categorization
- Add unit tests with 92% coverage
- Update transaction service AI.md

Closes #123
```

### Branch Naming
```bash
# Feature branches
feature/transaction-categories
feature/retry-logic
feature/telegram-webhooks

# Bug fixes
bugfix/memory-leak-in-processor
bugfix/invalid-currency-parsing

# Documentation
docs/api-documentation
docs/deployment-guide

# Refactoring
refactor/extract-ai-interface
refactor/simplify-error-handling
```

### Code Review Checklist
- [ ] Code follows Go conventions and project standards
- [ ] All functions have appropriate documentation
- [ ] Error handling uses internal/errors package
- [ ] Tests cover both success and error scenarios
- [ ] No secrets or sensitive data in code
- [ ] Package AI.md updated if business logic changed
- [ ] Imports are properly organized
- [ ] No use of log.Fatal or log.Panic
- [ ] Context is properly propagated
- [ ] Resource cleanup is handled (defer statements)

## 🔄 Continuous Improvement

### Regular Reviews
- Conduct weekly code review sessions
- Update guidelines based on lessons learned
- Review and refactor legacy code gradually
- Monitor test coverage and performance metrics

### Tool Integration
```bash
# Use these tools in CI/CD
go fmt ./...              # Format code
go vet ./...              # Static analysis
golangci-lint run         # Comprehensive linting
go test -race ./...       # Race condition detection
go test -cover ./...      # Test coverage
```

### Metrics and Monitoring
- Track test coverage (target: 85%+)
- Monitor error rates and types
- Measure API response times
- Review code complexity metrics

---

**Remember**: These guidelines ensure code quality, maintainability, and team collaboration. When in doubt, prioritize clarity and correctness over cleverness.