# Money Tracker Bot - Main Application

## Package: `cmd/telebot`

### Purpose
Main entry point for the Money Tracker Bot application. This package initializes all dependencies and starts the Telegram bot service.

### Key Components

#### `main.go`
- **Purpose**: Application bootstrap and dependency injection
- **Key Functions**:
  - `main()`: Entry point that loads environment variables and starts the bot
  - `startBot()`: Initializes services and starts the bot with real dependencies
  - `startBotWithDeps()`: Dependency injection wrapper for testing
- **Dependencies**:
  - Telegram bot API token (`TELEGRAM_BOT_TOKEN`)
  - Gemini API key (`GEMINI_API_KEY`)
  - Google Spreadsheet service
  - Gemini AI client

#### Error Handling
- `ErrEnvVarMissing`: Custom error type for missing environment variables

### Architecture
This package follows dependency injection patterns to enable testing and modular design. It orchestrates the initialization of:
- Google Spreadsheet service for data persistence
- Gemini AI client for transaction processing
- Transaction service for business logic
- Telegram handler for user interaction

### Environment Variables Required
- `TELEGRAM_BOT_TOKEN`: Bot token from Telegram BotFather
- `GEMINI_API_KEY`: API key for Google Gemini AI
- `GOOGLE_SPREADSHEET_ID`: ID of the Google Spreadsheet for data storage