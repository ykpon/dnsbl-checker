package main

import (
	"fmt"
	"net"
	"strings"
)

var listenIps []string = []string{
	"127.0.0.2",
	"127.0.0.3",
	"127.0.0.4",
}

var dnsblList map[string]string = map[string]string{
	"zen.spamhaus.org":       "SpamHaus",
	"bl.spamcop.net":         "SpamCop",
	"dnsbl-3.uceprotect.net": "UceProtect",
	"dnsbl.spfbl.net":        "SpfBL",
	"b.barracudacentral.org": "Barracuda",
	"noptr.spamrats.com":     "SpamRats",
}

// IPIsListed Check IP is listed in spam blacklists
func IPIsListed(ip string) (ipInList []string, listen bool) {
	sliceIP := strings.Split(ip, ".")

	if len(sliceIP) < 4 {
		return
	}

	reverseIP := fmt.Sprintf("%s.%s.%s.%s", sliceIP[3], sliceIP[2], sliceIP[1], sliceIP[0])

	for dnsbl, dnsblName := range dnsblList {
		address := fmt.Sprintf("%s.%s", reverseIP, dnsbl)
		lookupResult, err := lookupIP(address)
		if err != nil {
			// here may "no such host" error, so host is wrong or address not listed in bl
			continue
		}

		if len(lookupResult) > 0 {
			for _, v := range lookupResult {
				finded := false
				if finded {
					break
				}

				for _, listenIP := range listenIps {
					if v.String() == listenIP {
						ipInList = append(ipInList, dnsblName)
						finded = true
						break
					}
				}
			}
		}
	}

	if len(ipInList) > 0 {
		listen = true
	}

	return
}

func lookupIP(address string) (ips []net.IP, err error) {
	ips, err = net.LookupIP(address)
	return
}
