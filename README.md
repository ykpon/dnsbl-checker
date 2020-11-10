# TODO
- [x] Write tests
- [x] Dnsbl servers from text file
- [ ] Queue support, adding IPs from API and handle by queue
- [ ] Info about dnsbl servers, link to details for ip

## Example №1. Check all IP-Addresses in subnets
<details>
  <summary>main.go</summary>

```golang
package main

import (
	"log"
	"net"
	"sync"

	"github.com/korovkin/limiter"
	"github.com/ykpon/dnsbl-checker/lib/dnsbl"
)

func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	lenIPs := len(ips)
	switch {
	case lenIPs < 2:
		return ips, nil

	default:
		return ips[1 : len(ips)-1], nil
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func main() {
	subnets := []string{
		"50.19.0.0/16",
		"54.239.98.0/24",
	}

	limit := limiter.NewConcurrencyLimiter(10)
	for _, v := range subnets {
		tmp, err := Hosts(v)
		if err != nil {
			log.Fatal(err)
      continue
		}
		for _, ip := range tmp {
			limit.Execute(func() {
				dnsbl.IPIsListed(ip)
			})
		}
	}
	limit.Wait()
}
```

</details>

## Example №2. HTTP Server for checking IP-Address with notify to Telegram
<details>
  <summary>resources/config.json</summary>

```json
{
    "TELEGRAM_BOT_TOKEN": "BOT_TOKEN_FROM_@BOTFATHER",
    "TELEGRAM_CHANNEL_CHAT_ID": "CHANNEL_ID (prefix -100 required)"
}
```

</details>
<details>
  <summary>main.go</summary>

```golang
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

```
</details>

Usage:    
> curl --request GET "http://127.0.0.1:8000/dnsbl/1.2.3.4"
