package telegram

import (
	"context"
	"fmt"
	"io"
	"log"
	"money-tracker-bot/internal/service/transactions"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotAPI is an interface for sending messages (for testability)
type BotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type TelegramHandler struct {
	Telebot            BotAPI
	TransactionService transactions.ITransaction
}

// NewTelegramHandler creates a TelegramHandler with a real bot (for production)
func NewTelegramHandler(token string, transactionService transactions.ITransaction) *TelegramHandler {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	return &TelegramHandler{
		Telebot:            bot,
		TransactionService: transactionService,
	}
}

// NewTelegramHandlerWithBot allows injecting a bot instance (for testing)
func NewTelegramHandlerWithBot(bot BotAPI, transactionService transactions.ITransaction) *TelegramHandler {
	return &TelegramHandler{
		Telebot:            bot,
		TransactionService: transactionService,
	}
}

type StoredFile struct {
	FileID   string
	FileName string
	User     string
	Date     time.Time
}

var storedFiles []StoredFile

func (t *TelegramHandler) Start() {
	realBot, ok := t.Telebot.(*tgbotapi.BotAPI)
	if !ok {
		log.Panic("Telebot is not a *tgbotapi.BotAPI")
	}
	realBot.Debug = true
	log.Printf("Authorized on account %s", realBot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := realBot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "list":
				handleListCommand(t.Telebot, update.Message)
			case "view":
				handleViewCommand(realBot, update.Message)
			case "download":
				handleDownloadCommand(realBot, update.Message)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command.")
				t.Telebot.Send(msg)
			}
			continue
		}

		if update.Message.Document != nil {
			handleDocument(t.Telebot, update.Message)
		} else if update.Message.Photo != nil {
			t.handlePhoto(t.Telebot, update.Message)
		} else {
			t.handleMessage(t.Telebot, update.Message)
		}
	}
}

func handleListCommand(bot BotAPI, msg *tgbotapi.Message) {
	if len(storedFiles) == 0 {
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "No files received yet."))
		return
	}

	var text string
	for i, f := range storedFiles {
		text += fmt.Sprintf("%d. %s (from @%s, %s)\n", i+1, f.FileName, f.User, f.Date.Format("Jan 2 15:04"))
	}

	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, text))
}

func handleDocument(bot BotAPI, msg *tgbotapi.Message) {
	doc := msg.Document
	fileID := doc.FileID
	fileName := doc.FileName

	storedFiles = append(storedFiles, StoredFile{
		FileID:   fileID,
		FileName: fileName,
		User:     msg.From.UserName,
		Date:     time.Now(),
	})

	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Saved %s ✅", fileName)))
}

func (t *TelegramHandler) handlePhoto(bot BotAPI, msg *tgbotapi.Message) {
	photos := msg.Photo
	largest := photos[len(photos)-1]
	fileID := largest.FileID
	fileName := fmt.Sprintf("%s.jpg", fileID)
	localPath := "downloads/" + fileName

	// Cast to *tgbotapi.BotAPI for downloadFile
	realBot, ok := bot.(*tgbotapi.BotAPI)
	if !ok {
		log.Println("Bot is not *tgbotapi.BotAPI, skipping downloadFile")
		return
	}
	err := downloadFile(realBot, fileID, localPath)
	if err != nil {
		log.Println("Download error:", err)
		return
	}

	storedFiles = append(storedFiles, StoredFile{
		FileID:   fileID,
		FileName: fileName,
		User:     msg.From.UserName,
		Date:     time.Now(),
	})

	transaction, err := t.TransactionService.HandleImageInput(context.TODO(), localPath, msg.From.UserName, nil)
	if err != nil {
		log.Println("Error handling image input:", err)
		return
	}

   summary, _ := t.TransactionService.SaveTransaction(*transaction)
   spreadsheetId := os.Getenv("GOOGLE_SPREADSHEET_ID")
   spreadsheetLink := "https://docs.google.com/spreadsheets/d/" + spreadsheetId
   rupiah := formatRupiah(transaction.Amount)
   msgText := fmt.Sprintf(
	   "Saved photo ✅\nCategory: %s\nAmount: %s\nNotes: %s\nLink: %s\n"+
		   "Monthly Expenses: %s\nMonthly Budget: %s\nBudget Left: %s\n"+
		   "Monthly Quota: %s\nQuota Left: %s",
	   transaction.Category,
	   rupiah,
	   transaction.Notes,
	   spreadsheetLink,
	   summary.MonthlyExpenses,
	   summary.MonthlyBudget,
	   summary.BudgetLeft,
	   summary.Quota,
	   summary.QuotaLeft,
   )
   // Check budget and quota left, append Gemini's warning_message if needed
   budgetLeft, _ := strconv.ParseFloat(summary.BudgetLeft, 64)
   quotaLeft, _ := strconv.ParseFloat(summary.QuotaLeft, 64)
   if (budgetLeft < 0 || quotaLeft < 0) && transaction.WarningMessage != "" {
	   msgText += "\n\n⚠️ " + transaction.WarningMessage
   }
   bot.Send(tgbotapi.NewMessage(msg.Chat.ID, msgText))
}

func (t *TelegramHandler) handleMessage(bot BotAPI, msg *tgbotapi.Message) {
	transaction, err := t.TransactionService.HandleTextInput(context.TODO(), msg.Text, msg.From.UserName, nil)
	if err != nil {
		log.Println("Error handling text input:", err)
		return
	}

   summary, _ := t.TransactionService.SaveTransaction(*transaction)
   spreadsheetId := os.Getenv("GOOGLE_SPREADSHEET_ID")
   spreadsheetLink := "https://docs.google.com/spreadsheets/d/" + spreadsheetId
   rupiah := formatRupiah(transaction.Amount)
   msgText := fmt.Sprintf(
	   "Saved text ✅\nCategory: %s\nAmount: %s\nNotes: %s\nLink: %s\n"+
		   "Monthly Expenses: %s\nMonthly Budget: %s\nBudget Left: %s\n"+
		   "Monthly Quota: %s\nQuota Left: %s",
	   transaction.Category,
	   rupiah,
	   transaction.Notes,
	   spreadsheetLink,
	   summary.MonthlyExpenses,
	   summary.MonthlyBudget,
	   summary.BudgetLeft,
	   summary.Quota,
	   summary.QuotaLeft,
   )
   // Check budget and quota left, append Gemini's warning_message if needed
   budgetLeft, _ := strconv.ParseFloat(summary.BudgetLeft, 64)
   quotaLeft, _ := strconv.ParseFloat(summary.QuotaLeft, 64)
   if (budgetLeft < 0 || quotaLeft < 0) && transaction.WarningMessage != "" {
	   msgText += "\n\n⚠️ " + transaction.WarningMessage
   }
   bot.Send(tgbotapi.NewMessage(msg.Chat.ID, msgText))
}

// formatRupiah formats a string amount to Indonesian Rupiah currency
func formatRupiah(amount string) string {
	// Try to parse as float, fallback to original string
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return "Rp " + amount
	}
	// Format with thousands separator
	return fmt.Sprintf("Rp %s", formatThousands(int64(f)))
}

// formatThousands formats an integer with thousands separator
func formatThousands(n int64) string {
	s := fmt.Sprintf("%d", n)
	var out []byte
	for i, c := range s {
		if i != 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	return string(out)
}

func downloadFile(bot *tgbotapi.BotAPI, fileID, localPath string) error {
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return err
	}

	url := file.Link(bot.Token)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func parseIndexArg(text string) (int, error) {
	parts := strings.Split(text, " ")
	if len(parts) < 2 {
		return -1, fmt.Errorf("missing index")
	}
	i, err := strconv.Atoi(parts[1])
	if err != nil || i < 1 || i > len(storedFiles) {
		return -1, fmt.Errorf("invalid index")
	}
	return i - 1, nil
}

func handleViewCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	index, err := parseIndexArg(msg.Text)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Usage: /view <number>"))
		return
	}

	file := storedFiles[index]
	photo := tgbotapi.NewPhoto(msg.Chat.ID, tgbotapi.FilePath("downloads/"+file.FileName))
	photo.Caption = fmt.Sprintf("Viewing: %s", file.FileName)
	bot.Send(photo)
}

func handleDownloadCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	index, err := parseIndexArg(msg.Text)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Usage: /download <number>"))
		return
	}

	file := storedFiles[index]
	doc := tgbotapi.NewDocument(msg.Chat.ID, tgbotapi.FilePath("downloads/"+file.FileName))
	doc.Caption = fmt.Sprintf("Download: %s", file.FileName)
	bot.Send(doc)
}
