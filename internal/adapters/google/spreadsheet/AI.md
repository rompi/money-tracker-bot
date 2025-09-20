# Google Spreadsheet Adapter

## Package: `internal/adapters/google/spreadsheet`

### Purpose
Adapter for Google Sheets API that provides data persistence for transaction records and budget tracking.

### Key Components

#### `client.go`
- **Purpose**: Google Sheets API client for transaction data management
- **Key Structures**:
  - `SpreadsheetService`: Main service for Google Sheets operations
  - `CategorySummary`: Budget and quota summary for categories
- **Key Functions**:
  - `AppendRow()`: Adds new transaction records to the detailed sheet
  - `GetCellValue()`: Reads data from specific cells (utility function)

#### Data Management
- **Transaction Storage**: Stores transactions in "detailed" sheet with columns:
  - Transaction Date
  - Category
  - Notes
  - Amount
  - Created By
  - File ID
  - Created At (UTC+7 timezone)

- **Budget Tracking**: Reads from "summary" sheet for:
  - Monthly expenses by category
  - Monthly budget limits
  - Budget remaining
  - Shopping quota limits
  - Quota remaining

#### Features
- **Timezone Support**: Uses Asia/Bangkok timezone (UTC+7) for timestamps
- **Category Matching**: Links transactions to budget categories
- **Real-time Updates**: Immediately reflects budget changes after transactions
- **Defensive Programming**: Handles missing quota columns gracefully

#### Dependencies
- Google Sheets API v4 (`google.golang.org/api/sheets/v4`)
- Google service account credentials (`google-service-account.json`)
- Transaction domain models

### Sheet Structure
- **detailed**: Transaction records with full details
- **summary**: Category-based budget and quota tracking (A2:F12 range)