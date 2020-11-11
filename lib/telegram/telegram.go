package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// TelegramBot struct ...
type TelegramBot struct {
	Bot       *tgbotapi.BotAPI
	Token     string
	ChatID    int64
	Connected bool
}

// Init telegram bot
func (t *TelegramBot) Init() {
	var err error
	t.Bot, err = tgbotapi.NewBotAPI(t.Token)
	if err != nil {
		log.Panicln(err)
	}
	t.Connected = true
}

// SendMessageToChannel - send message to channel of telegram
func (t TelegramBot) SendMessageToChannel(text string) {
	msg := tgbotapi.NewMessage(t.ChatID, text)
	t.Bot.Send(msg)
}
