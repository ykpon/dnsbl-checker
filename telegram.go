package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type telegramBot struct {
	bot    *tgbotapi.BotAPI
	token  string
	chatID int64
}

func (t *telegramBot) init() {
	var err error
	t.bot, err = tgbotapi.NewBotAPI(t.token)
	if err != nil {
		log.Panicln(err)
	}

}

func (t telegramBot) sendMessageToChannel(text string) {
	msg := tgbotapi.NewMessage(t.chatID, text)
	t.bot.Send(msg)
}
