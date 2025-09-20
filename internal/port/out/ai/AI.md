# AI Port Interface

## Package: `internal/port/out/ai`

### Purpose
Output port interface for AI services, defining the contract for artificial intelligence operations in the money tracker application.

### Key Components

#### `aiport.go`
- **Purpose**: AI service interface definition following hexagonal architecture
- **Key Interface**:
  - `AiPort`: Contract for AI-powered transaction processing

### Interface Definition
The `AiPort` interface provides three core AI operations:

#### Methods
1. **`GenerateContent(ctx context.Context, prompt string)`**
   - General-purpose content generation
   - Context-aware for cancellation support
   - Used for custom AI prompts

2. **`ReadImageToTransaction(ctx context.Context, imgPath string) (*Transaction, error)`**
   - Processes receipt/transaction images
   - Extracts structured transaction data from images
   - Returns domain transaction model or error

3. **`TextToTransaction(ctx context.Context, message string) (*Transaction, error)`**
   - Processes natural language transaction descriptions
   - Converts text to structured transaction data
   - Returns domain transaction model or error

### Architecture Benefits
- **Dependency Inversion**: Business logic depends on interface, not implementation
- **Testability**: Easy mocking for unit tests
- **Flexibility**: Supports multiple AI service implementations
- **Context Support**: Enables timeout and cancellation handling

### Implementation Notes
- Currently implemented by the Gemini adapter
- Supports dependency injection in transaction service
- Enables testing with mock AI services
- Part of hexagonal architecture pattern (ports and adapters)

### Usage Context
This interface is injected into the transaction service to provide AI capabilities while maintaining clean architecture boundaries.