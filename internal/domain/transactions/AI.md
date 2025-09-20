# Transaction Domain

## Package: `internal/domain/transactions`

### Purpose
Domain model for transaction data structure, representing the core business entity for financial transactions.

### Key Components

#### `transactions.go`
- **Purpose**: Transaction domain model definition
- **Key Structures**:
  - `Transaction`: Core business entity representing a financial transaction

### Transaction Model
The `Transaction` struct represents a complete financial transaction with the following fields:

#### Core Transaction Data
- `TransactionDate`: Date of the transaction (YYYY-MM-DD format)
- `Amount`: Transaction amount (always positive, in Indonesian Rupiah)
- `AmountCurrency`: Currency designation (typically IDR)
- `Category`: Expense category (Groceries, Utilities, Entertainment, etc.)
- `Notes`: Detailed description of the transaction

#### Payment Information
- `DestinationName`: Name of the payment recipient
- `DestinationNumber`: Account/phone number of recipient
- `SourceAccount`: Source payment method (GOPAY, BCA, OVO, etc.)

#### Metadata
- `Title`: Summary/title of the transaction
- `FileID`: Associated file identifier (for image uploads)
- `CreatedBy`: User who created the transaction
- `WarningMessage`: Optional budget/quota warning message

### Design Principles
- **JSON Serialization**: All fields support JSON marshaling for API responses
- **Positive Amounts**: Amount is always stored as positive value
- **Immutable Structure**: Represents a snapshot of transaction data
- **Rich Metadata**: Includes audit trail and file associations

### Usage Context
This domain model is used throughout the application for:
- AI-extracted transaction data
- Spreadsheet persistence
- API responses
- User interface display