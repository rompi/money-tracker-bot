# Money Tracker Bot - Technical Documentation

## Architecture Overview

The Money Tracker Bot is a Telegram bot built in Go that processes financial documents and receipts using AI to track expenses. The application follows a clean architecture pattern with clear separation of concerns.

### Core Components

1. **Telegram Bot Handler** (`internal/adapters/telegram`)
   - Handles incoming Telegram messages and file uploads
   - Manages user interactions and bot responses
   - Implements message routing and command handling

2. **Transaction Service** (`internal/service/transactions`)
   - Core business logic layer
   - Orchestrates the flow between different components
   - Processes transactions and manages data flow

3. **Gemini AI Integration** (`internal/adapters/gemini`)
   - Handles communication with Google's Gemini AI
   - Processes images and extracts structured data
   - Converts AI responses into transaction data

4. **Google Sheets Integration** (`internal/adapters/google/spreadsheet`)
   - Manages data persistence in Google Sheets
   - Handles spreadsheet operations and updates

## Technical Stack

- **Language**: Go 1.21+
- **Bot Framework**: `go-telegram-bot-api`
- **AI Integration**: Google Generative AI Go SDK (`generative-ai-go`)
- **Data Storage**: Google Sheets API (`google.golang.org/api/sheets/v4`)

## Configuration

The application requires the following environment variables:

- `TELEGRAM_BOT_TOKEN`: Your Telegram bot token from BotFather
- `GEMINI_API_KEY`: API key for Google's Gemini AI service
- Google Service Account credentials (stored in `google-service-account.json`)

## Service Flow

1. **Initialization**
   ```go
   func startBot() error {
       // Load environment variables
       // Initialize services
       // Start the Telegram bot
   }
   ```

2. **Message Processing**
   - Bot receives photos or documents from users
   - Files are temporarily downloaded to the `downloads/` directory
   - Images are processed through Gemini AI
   - Structured data is extracted and stored
   - User receives a confirmation message
   - Downloaded files are cleaned up

3. **Data Extraction**
   The service extracts the following fields from transactions:
   - Transaction datetime
   - Amount and currency
   - Notes
   - Destination name and number
   - Source account
   - Category
   - File ID

## Testing

The application is designed with testability in mind:

- Uses dependency injection for easier unit testing
- Includes mock implementations for external services
- Separate test files for each major component

## Error Handling

The application implements robust error handling:
- Environment variable validation
- Service initialization checks
- Graceful error responses to users
- Clean shutdown procedures

## Dependencies

All dependencies are managed through Go modules. Main external dependencies include:
- `github.com/joho/godotenv` for environment variable management
- Telegram Bot API client
- Google API clients for Gemini AI and Sheets

## Security Considerations

1. **API Keys and Tokens**
   - Stored in environment variables
   - Not committed to version control
   - Loaded at runtime

2. **File Handling**
   - Temporary files are automatically cleaned up
   - Secure file download handling
   - Validated file types and sizes

3. **Service Account**
   - Google service account with minimal required permissions
   - Credentials stored securely in `google-service-account.json`

## Deployment

The application can be deployed using the following methods:

1. **Direct Execution**
   ```bash
   go run cmd/telebot/main.go
   ```

2. **Build and Run**
   ```bash
   go build -o bot cmd/telebot/main.go
   ./bot
   ```

Make sure all environment variables and service account credentials are properly configured before deployment.
