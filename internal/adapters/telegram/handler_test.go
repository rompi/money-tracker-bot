package telegram

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestNewTelegramHandler(t *testing.T) {
	mockBot := &MockBotAPI{}
	h := NewTelegramHandlerWithBot(mockBot, &MockTransactionService{})
	if h.Telebot == nil {
		t.Error("Telebot should be initialized")
	}
	if h.TransactionService == nil {
		t.Error("TransactionService should be set")
	}
}

func TestHandleMessage_CallsService(t *testing.T) {
	m := &MockTransactionService{}
	mockBot := &MockBotAPI{}
	h := &TelegramHandler{
		Telebot:            mockBot,
		TransactionService: m,
	}
	msg := &tgbotapi.Message{
		Text: "test",
		From: &tgbotapi.User{UserName: "user"},
		Chat: &tgbotapi.Chat{ID: 12345},
	}
	h.handleMessage(mockBot, msg)
	if !m.HandleTextInputCalled {
		t.Error("HandleTextInput should be called")
	}
	if !m.SaveTransactionCalled {
		t.Error("SaveTransaction should be called")
	}
	if len(mockBot.SentMessages) == 0 {
		t.Error("Bot should have sent a message")
	}
}
