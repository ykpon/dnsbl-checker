package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ykpon/dnsbl-checker/config"
	"github.com/ykpon/dnsbl-checker/lib/dnsbl"
)

var bot telegramBot

func findIP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	go func(p map[string]string) {
		dnsbls, isListed := dnsbl.IPIsListed(p["ip"])

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

func main() {
	config := config.LoadConf()

	bot = telegramBot{token: config.TelegramBotToken, chatID: config.TelegramChannelChatID}
	bot.init()

	r := mux.NewRouter()
	r.HandleFunc("/dnsbl/{ip}", findIP).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
