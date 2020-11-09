package main

import (
	"encoding/json"
	"os"
)

type config struct {
	TelegramBotToken      string `json:"TELEGRAM_BOT_TOKEN"`
	TelegramChannelChatID int64  `json:"TELEGRAM_CHANNEL_CHAT_ID"`
}

func loadConf() config {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}

	return configuration
}
