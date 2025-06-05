package telegram

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"rompi/gobot/internal/service/transactions"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramHandler struct {
	Telebot            *tgbotapi.BotAPI
	TransactionService transactions.ITransaction
}

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

type StoredFile struct {
	FileID   string
	FileName string
	User     string
	Date     time.Time
}

var storedFiles []StoredFile

func (t *TelegramHandler) Start() {
	bot := t.Telebot
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "list":
				handleListCommand(bot, update.Message)
			case "view":
				handleViewCommand(bot, update.Message)
			case "download":
				handleDownloadCommand(bot, update.Message)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command.")
				bot.Send(msg)
			}
			continue
		}

		if update.Message.Document != nil {
			handleDocument(bot, update.Message)
		} else if update.Message.Photo != nil {
			t.handlePhoto(bot, update.Message)
		} else {
			t.handleMessage(bot, update.Message)
		}
	}
}

func handleListCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
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

func handleDocument(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
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

func (t *TelegramHandler) handlePhoto(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	photos := msg.Photo
	largest := photos[len(photos)-1]
	fileID := largest.FileID
	fileName := fmt.Sprintf("%s.jpg", fileID)
	localPath := "downloads/" + fileName

	err := downloadFile(bot, fileID, localPath)
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

	t.TransactionService.SaveTransaction(*transaction)

	spreadsheetId := os.Getenv("GOOGLE_SPREADSHEET_ID")
	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Saved photo ✅ as %s \n\ntotal amount %s. link = %s", transaction.Notes, transaction.Amount, "https://docs.google.com/spreadsheets/d/"+spreadsheetId)))
}

func (t *TelegramHandler) handleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	transaction, err := t.TransactionService.HandleTextInput(context.TODO(), msg.Text, msg.From.UserName, nil)
	if err != nil {
		log.Println("Error handling text input:", err)
		return
	}

	t.TransactionService.SaveTransaction(*transaction)

	spreadsheetId := os.Getenv("GOOGLE_SPREADSHEET_ID")
	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Saved text ✅ as %s \n\ntotal amount %s. link = %s", transaction.Notes, transaction.Amount, "https://docs.google.com/spreadsheets/d/"+spreadsheetId)))
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
