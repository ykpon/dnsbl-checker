package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ykpon/dnsbl-checker/config"
	"github.com/ykpon/dnsbl-checker/lib/dnsbl"
	"github.com/ykpon/dnsbl-checker/lib/telegram"
)

var bot telegram.TelegramBot

func findIP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	params := mux.Vars(r)
	_, notify := r.URL.Query()["notify"]
	dnsbls, isListed := dnsbl.IPIsListed(params["ip"])
	var msg string
	if isListed {
		msg = fmt.Sprintf("Info about IP: %s", params["ip"])
		msg += "\nAddress in blacklists:"
		for _, dnsbl := range dnsbls {
			msg += fmt.Sprintf("\n%s", dnsbl)
		}

		if notify && bot.Bot != nil && bot.Connected {
			bot.SendMessageToChannel(msg)
		}
	}

	fmt.Fprintf(w, "%v\n", msg)
}

func main() {
	config := config.LoadConf()

	bot = telegram.TelegramBot{Token: config.TelegramBotToken, ChatID: config.TelegramChannelChatID}
	bot.Init()

	r := mux.NewRouter()
	r.HandleFunc("/dnsbl/{ip}", findIP).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
