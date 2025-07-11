package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MockBotAPI struct {
	SentMessages []tgbotapi.Chattable
}

func (m *MockBotAPI) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	m.SentMessages = append(m.SentMessages, c)
	return tgbotapi.Message{}, nil
}
