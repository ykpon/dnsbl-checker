package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gorilla/mux"
)

type telegramBot struct {
	bot    *tgbotapi.BotAPI
	err    error
	token  string
	chatID int64
}

type config struct {
	TelegramBotToken      string `json:"TELEGRAM_BOT_TOKEN"`
	TelegramChannelChatID int64  `json:"TELEGRAM_CHANNEL_CHAT_ID"`
}

var bot telegramBot

func (t *telegramBot) init() {
	t.bot, t.err = tgbotapi.NewBotAPI(t.token)

}
func (t telegramBot) sendMessageToChannel(text string) {
	msg := tgbotapi.NewMessage(t.chatID, text)
	t.bot.Send(msg)
}

func findIP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	go func(p map[string]string) {
		dnsbls, isListed := IPIsListed(p["ip"])

		if isListed {
			msg := fmt.Sprintf("Info about IP: %s", p["ip"])
			msg += "\nAddress in blacklists:"
			for _, dnsbl := range dnsbls {
				msg += fmt.Sprintf("\n%s", dnsbl)
			}
			bot.sendMessageToChannel(msg)
		}
	}(params)

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

func main() {
	config := loadConf()
	fmt.Println(config)

	bot = telegramBot{token: config.TelegramBotToken, chatID: config.TelegramChannelChatID}
	bot.init()

	r := mux.NewRouter()
	r.HandleFunc("/dnsbl/{ip}", findIP).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
