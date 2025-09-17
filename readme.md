# 💸 Telegram Money Bot + Gemini AI (Golang)

This project is a Telegram bot written in Go that:

- Receives **photos or documents** from users (e.g., receipts, transactions)
- **Downloads** the files from Telegram
- Sends them to **Gemini AI (via Google's SDK)** for intelligent parsing
- Extracts structured **transaction data as JSON**
- Replies to the user with a summary
- Stores the results in a **Google Spreadsheet** for tracking
- Cleans up the downloaded files afterward

---

## 🚀 Features

- ✅ Telegram bot handler using [`go-telegram-bot-api`](https://github.com/go-telegram-bot-api/telegram-bot-api)
- ✅ Gemini AI integration using [`generative-ai-go`](https://github.com/google/generative-ai-go)
- ✅ Google Sheets integration using [`google.golang.org/api/sheets/v4`](https://pkg.go.dev/google.golang.org/api/sheets/v4)
- ✅ Extracts fields like:
  - `transaction_datetime`
  - `amount` & `currency`
  - `notes`
  - `destination_name` & `number`
  - `source_account`
  - `category`
  - `file_id`
- ✅ Responds with Markdown messages
- ✅ Auto-removes files after processing
- ✅ Goroutine-safe and gracefully shutdown

---

## 📦 Project Structure

```
.
├── cmd/bot             # Main entry point
├── internal/           # Bot logic, Gemini integration, Sheets handler
├── downloads/          # Temporary downloaded images
├── .env                # Secrets (TELEGRAM_BOT_TOKEN, GEMINI_API_KEY, SHEETS_CREDENTIALS)
└── README.md
```

---

## 🧭 Contributor Guide

New contributors should start with `AGENTS.md` for repository conventions, workflows, and security tips.

---

## ⚙️ Requirements

- Go 1.21+
- A Telegram Bot Token from [@BotFather](https://t.me/BotFather)
- A Gemini API Key from [Google AI Studio](https://makersuite.google.com/app)
- A Google Sheets API service account & spreadsheet access

---

## ⚙️ Setup

1. Clone the repo

```bash
git clone https://github.com/rompi/money-tracker-bot.git
cd money-tracker-bot
```

2. Install dependencies

```bash
go mod tidy
```

3. Create `.env` file

```env
TELEGRAM_BOT_TOKEN=your_telegram_token
GEMINI_API_KEY=your_gemini_key
SPREAD
```

1. Create `google-service-account.json` file that you got from google api

```google-service-account.json
{
  "type": "service_account",
  "project_id": "your-project",
  "private_key_id": "1234",
  "private_key": "-----BEGIN PRIVATE KEY-----\nasdasd
  \n-----END PRIVATE KEY-----\n",
  "client_email": "service-account@gmail.com",
  "client_id": "1111111111",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509",
  "universe_domain": "googleapis.com"
}

```

5. Run the bot

```bash
go run ./cmd/bot
```

---

## ✨ Example

Send a photo of a receipt or payment transaction to your bot:
The bot will reply with:

```
💾 Transaction Info
• Amount: `150.00`
• Date: `2025-03-30T12:30:00`
• Notes: Lunch at ABC cafe
• Destination: ABC Cafe
• Category: `food`
📁 File: `photo_1743586322.jpg`
```

And the transaction will be saved into your Google Spreadsheet.

---

## 📔 License

MIT — use freely and modify as needed.

---

## 🙇‍♂️ Questions or Suggestions?

Open an issue or message me on Telegram.
