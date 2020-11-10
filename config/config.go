package config

import (
	"encoding/json"
	"os"
)

// Config struct for file config.json
type Config struct {
	TelegramBotToken      string `json:"TELEGRAM_BOT_TOKEN"`
	TelegramChannelChatID int64  `json:"TELEGRAM_CHANNEL_CHAT_ID"`
}

// LoadConf func ...
func LoadConf() Config {
	file, _ := os.Open("resources/config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}

	return configuration
}
