# Common Utilities

## Package: `internal/common`

### Purpose
Shared utilities and constants used across the Money Tracker Bot application, particularly for AI prompt generation and transaction categorization.

### Key Components

#### `prompt.go`
- **Purpose**: AI prompt building utilities for consistent transaction processing
- **Key Functions**:
  - `BuildPrompt()`: Constructs structured prompts for Gemini AI
- **Key Constants**:
  - `TransactionCategoryList`: Predefined expense categories
  - `SourceAccountList`: Supported payment methods

### Transaction Categories
Predefined categories for expense classification:
- Groceries
- Utilities
- Entertainment
- Gifting
- Household
- Eating Out
- Health
- Transportation
- Savings
- Emergency
- Rent House

### Source Accounts
Supported payment methods:
- GOPAY
- BCA
- OVO
- DANA
- ISAKU
- MANDIRI
- BNI
- BRI
- CASH

### Prompt Building System
The `BuildPrompt()` function creates structured prompts for AI processing:

#### Features
- **Context-Aware**: Handles both image and text inputs differently
- **Structured Output**: Ensures consistent JSON response format
- **Field Validation**: Includes predefined categories and accounts
- **Date Handling**: Manages current date for text inputs
- **File Association**: Links file IDs for image inputs

#### Prompt Structure
- Field definitions with validation rules
- Category and account constraints
- Example JSON output format
- Currency formatting guidelines (Indonesian Rupiah)
- Warning message generation for budget alerts

### Design Principles
- **Consistency**: Standardized prompts across all AI interactions
- **Validation**: Ensures data conforms to business rules
- **Flexibility**: Supports both image and text processing modes
- **Localization**: Indonesian currency and business context