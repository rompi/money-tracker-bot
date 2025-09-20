# 💸 Money Tracker Bot

A Telegram bot that intelligently tracks your financial transactions by analyzing receipts and documents using AI, then automatically organizing the data in Google Sheets.

## 📋 What This Bot Does

### Core Business Features

1. **Receipt Analysis**
   - Upload photos of receipts, invoices, or transaction screenshots
   - AI extracts key information: amount, date, merchant, category
   - Supports multiple currencies and formats

2. **Smart Transaction Processing**
   - Automatically categorizes expenses (Food, Transport, Shopping, etc.)
   - Extracts merchant names and transaction details
   - Handles both income and expense transactions

3. **Google Sheets Integration**
   - Automatically saves all transactions to your spreadsheet
   - Organized columns: Date, Amount, Category, Merchant, Notes, File ID
   - Real-time updates with transaction history

4. **Telegram Bot Interface**
   - Simple: Just send a photo to the bot
   - Instant feedback with transaction summary
   - Works with photos, documents, and even text descriptions

### Supported Transaction Types
- **Receipts**: Restaurant bills, shopping receipts, service invoices
- **Bank Statements**: Screenshot or photo of transaction history
- **Digital Receipts**: Online purchase confirmations, e-wallet transactions
- **Manual Entries**: Text descriptions like "Coffee $4.50 at Starbucks"

## 🛠️ Technology Stack

- **Backend**: Go 1.23+ with hexagonal architecture
- **AI Processing**: Google Gemini AI for document analysis
- **Storage**: Google Sheets API for data persistence
- **Interface**: Telegram Bot API for user interaction
- **Error Handling**: Robust error infrastructure with graceful degradation

## 📦 Project Structure

```
money-tracker-bot/
├── cmd/telebot/               # Application entry point
│   └── main.go               # Bootstrap and dependency injection
├── internal/
│   ├── errors/               # Centralized error handling system
│   ├── adapters/             # External service integrations
│   │   ├── telegram/         # Telegram Bot API adapter
│   │   ├── google/           # Google Sheets API adapter
│   │   └── gemini/           # Gemini AI service adapter
│   ├── service/transactions/ # Core business logic
│   ├── domain/transactions/  # Domain models and entities
│   └── port/out/ai/         # AI service interface definitions
├── scripts/                  # Development and deployment scripts
├── .env.example             # Environment variables template
└── google-service-account.json # Google API credentials (not in repo)
```

## 🚀 Getting Started

### Prerequisites

#### System Requirements
- **Go**: Version 1.23 or higher
- **Git**: For version control
- **Make**: For build commands (optional, can use `go` commands directly)

#### External Services
1. **Telegram Bot Token**
   - Message [@BotFather](https://t.me/BotFather) on Telegram
   - Create a new bot with `/newbot`
   - Save the bot token

2. **Google Gemini API Key**
   - Visit [Google AI Studio](https://makersuite.google.com/app)
   - Create an API key for Gemini
   - Enable the Generative AI API

3. **Google Sheets API Access**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select existing one
   - Enable Google Sheets API
   - Create service account credentials
   - Download the JSON credentials file

4. **Google Spreadsheet**
   - Create a new Google Spreadsheet
   - Share it with your service account email (found in credentials JSON)
   - Give "Editor" permissions
   - Copy the spreadsheet ID from the URL

### Installation & Setup

1. **Clone the repository**
```bash
git clone https://github.com/rompi/money-tracker-bot.git
cd money-tracker-bot
```

2. **Install Go dependencies**
```bash
go mod tidy
```

3. **Set up environment variables**
Create a `.env` file in the project root:
```env
# Telegram Bot Configuration
TELEGRAM_BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrSTUvwxYZ

# Google Gemini AI Configuration
GEMINI_API_KEY=your_gemini_api_key_here

# Google Sheets Configuration
SPREADSHEET_ID=1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms
```

4. **Add Google service account credentials**
Place your Google service account JSON file as `google-service-account.json`:
```json
{
  "type": "service_account",
  "project_id": "your-project-id",
  "private_key_id": "key-id",
  "private_key": "-----BEGIN PRIVATE KEY-----\nYOUR_PRIVATE_KEY\n-----END PRIVATE KEY-----\n",
  "client_email": "your-service-account@project.iam.gserviceaccount.com",
  "client_id": "1234567890",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/your-service-account%40project.iam.gserviceaccount.com"
}
```

5. **Run the application**
```bash
# Using Make (recommended)
make run

# Or directly with Go
go run ./cmd/telebot
```

### Verification

1. **Test the bot**: Send a message to your Telegram bot
2. **Upload a receipt**: Send a photo of a receipt to test AI processing
3. **Check Google Sheets**: Verify transactions appear in your spreadsheet

### Available Commands
```bash
make run           # Start the bot
make test          # Run all tests with coverage
make fmt           # Format Go code
make lint          # Run linting checks
make build         # Build production binary
```

## ✨ How to Use

### Basic Usage
1. **Start a chat** with your Telegram bot
2. **Send a photo** of any receipt, invoice, or transaction
3. **Receive confirmation** with extracted transaction details
4. **Check your Google Sheet** for the automatically added transaction

### Example Interaction
```
You: [Send photo of Starbucks receipt]

Bot: 💾 Transaction Processed Successfully!
     • Amount: $15.50 USD
     • Date: 2025-01-20 14:30:00
     • Category: Food & Dining
     • Merchant: Starbucks Coffee
     • Notes: Iced coffee and sandwich
     • File ID: photo_123456789

Your transaction has been saved to Google Sheets.
```

### Supported Input Types
- **Photos**: JPG, PNG receipt images
- **Documents**: PDF invoices and statements
- **Text**: Manual transaction descriptions
- **Screenshots**: Bank app or e-wallet transactions

---

## 🤝 How to Contribute

We welcome contributions to make this bot even better! Whether you're fixing bugs, adding features, or improving documentation, your help is appreciated.

### Quick Start for Contributors

1. **Fork the repository** on GitHub
2. **Clone your fork** locally
3. **Read the development context** in `AGENTS.md`
4. **Set up the development environment** (follow installation steps above)
5. **Create a feature branch** for your changes
6. **Make your changes** following our coding standards
7. **Submit a pull request** with a clear description

### Development Requirements

When contributing code, you must:

#### 📝 Documentation
- Update package `AI.md` files to reflect business/feature changes
- Add comprehensive function docstrings with examples:
```go
// ProcessTransaction analyzes transaction text and extracts structured data.
// Input: "Coffee at Starbucks $4.50"
// Output: Transaction{Amount: 4.50, Category: "Food", Merchant: "Starbucks"}
func ProcessTransaction(text string) (*Transaction, error) {
    // implementation
}
```

#### 🧪 Testing
- Write unit tests for all new functions
- Achieve **minimum 85% code coverage** for new code
- Test both success and error scenarios
- Verify with: `go test -cover ./...`

#### 🎨 Code Quality
- **Follow comprehensive standards** in `CODING-GUIDELINES.md`
- Use error constructors from `internal/errors` package
- Follow hexagonal architecture principles
- Add meaningful context to errors
- Never use `log.Fatal` - return errors instead
- Format code with `make fmt`
- Pass linting with `make lint`

### Git Workflow Options

#### Option A: Automated Workflow (Recommended)
Use our git workflow script for quality enforcement:

```bash
# One-time setup
./scripts/git-workflow.sh setup

# Start new feature
./scripts/git-workflow.sh branch your-feature-name

# After making changes
git add .
./scripts/git-workflow.sh check          # Run all quality checks
./scripts/git-workflow.sh commit "feat: your feature description"

# Push and create PR
git push origin feature/your-feature-name
```

#### Option B: Manual Workflow
For manual control:

```bash
# Create feature branch
git checkout -b feature/your-feature-name

# Make your changes, then run quality checks
make test          # Must pass with ≥85% coverage
make fmt           # Format code
make lint          # Check for issues

# Commit and push
git add .
git commit -m "feat: your feature description"
git push origin feature/your-feature-name
```

### Pre-commit Checklist
Before submitting any pull request:
- [ ] Package `AI.md` files updated
- [ ] Function docstrings with input/output examples added
- [ ] Unit tests written with ≥85% coverage
- [ ] `make test` passes completely
- [ ] `make fmt && make lint` runs clean
- [ ] Architecture diagrams added (if needed, use PlantUML format)

### What We're Looking For

#### 🐛 Bug Fixes
- Clear description of the issue
- Steps to reproduce
- Fix with appropriate tests

#### ✨ New Features
- AI transaction processing improvements
- Better error handling and user feedback
- Integration enhancements (Telegram, Sheets, Gemini)
- Performance optimizations

#### 📚 Documentation
- Code examples and usage guides
- Architecture documentation
- API documentation improvements

### Need Help?

- **Read**: `AGENTS.md` for detailed development context
- **Standards**: `CODING-GUIDELINES.md` for comprehensive coding standards
- **Issues**: Check [existing issues](https://github.com/rompi/money-tracker-bot/issues)
- **Questions**: Open a discussion or create an issue
- **Architecture**: Review the hexagonal architecture pattern

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🆘 Support & Community

- **🐛 Bug Reports**: [GitHub Issues](https://github.com/rompi/money-tracker-bot/issues)
- **💡 Feature Requests**: [GitHub Discussions](https://github.com/rompi/money-tracker-bot/discussions)
- **📖 Documentation**: Check `AGENTS.md` and package `AI.md` files
- **🤝 Contributing**: Follow the guidelines above

---

**Made with ❤️ for better financial tracking**
