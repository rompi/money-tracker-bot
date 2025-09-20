# Telegram Adapter

## Package: `internal/adapters/telegram`

### Purpose
Telegram bot adapter that handles user interactions, message processing, and file management for the Money Tracker Bot.

### Key Components

#### `handler.go`
- **Purpose**: Main Telegram bot handler with message processing logic
- **Key Structures**:
  - `TelegramHandler`: Main bot handler with transaction service integration
  - `StoredFile`: Represents uploaded files with metadata
- **Key Functions**:
  - `Start()`: Main bot event loop for processing updates
  - `handlePhoto()`: Processes photo uploads and extracts transaction data
  - `handleMessage()`: Processes text messages for transaction extraction
  - `handleDocument()`: Manages document uploads
  - Commands: `/list`, `/view`, `/download` for file management

#### Features
- **Transaction Processing**: Converts photos and text to transaction records
- **File Management**: Stores and manages uploaded files with metadata
- **Budget Monitoring**: Displays monthly expenses, budget, and quota information
- **Currency Formatting**: Formats amounts in Indonesian Rupiah
- **Warning System**: Shows alerts when budget or quota limits are exceeded

#### Message Flow
1. User sends photo/text → Bot processes with AI → Saves to spreadsheet → Returns formatted summary
2. File commands allow users to list, view, and download previously uploaded files

#### Dependencies
- Telegram Bot API (`github.com/go-telegram-bot-api/telegram-bot-api/v5`)
- Transaction service for business logic
- Environment variables for spreadsheet integration

### Testing Support
- `BotAPI` interface for mocking Telegram API calls
- Dependency injection pattern for transaction service
- Separate constructors for production and testing