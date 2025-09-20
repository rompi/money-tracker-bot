# Transaction Service

## Package: `internal/service/transactions`

### Purpose
Core business logic service for handling transaction processing, including image and text input processing, and data persistence.

### Key Components

#### `svc.go`
- **Purpose**: Interface definition for transaction operations
- **Key Interface**:
  - `ITransaction`: Core contract for transaction processing

#### `handler.go`
- **Purpose**: Transaction service implementation with dependency management
- **Key Structures**:
  - `TransactionService`: Main service with AI and spreadsheet dependencies
  - `SpreadsheetServicePort`: Interface for spreadsheet operations
- **Key Functions**:
  - `SaveTransaction()`: Persists transaction data to Google Sheets
  - `HandleImageInput()`: Processes receipt images into transaction records
  - `HandleTextInput()`: Converts text messages into transactions

#### Business Logic Flow
1. **Input Processing**:
   - Accepts image paths or text messages
   - Uses AI service (with fallback support) for data extraction
   - Adds user metadata (uploader information)

2. **Data Persistence**:
   - Saves structured transaction data to Google Sheets
   - Returns category summary with budget information

3. **Dependency Injection**:
   - Supports AI service injection for testing
   - Uses default AI service when none provided
   - Abstracts spreadsheet operations through interface

#### Features
- **Multi-Input Support**: Handles both image and text inputs
- **User Tracking**: Associates transactions with uploaders
- **Flexible AI Integration**: Supports different AI implementations
- **Error Handling**: Graceful handling of processing failures

#### Dependencies
- AI port interface for transaction extraction
- Google Spreadsheet adapter for data persistence
- Transaction domain models
- Context support for operation cancellation