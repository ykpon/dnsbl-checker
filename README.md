# TODO
- [x] Write tests
- [x] Dnsbl servers from text file
- [ ] Queue support, adding IPs from API and handle by queue
- [ ] Info about dnsbl servers, link to details for ip

## Example №1 Check all IP-Addresses in subnets
```golang
package main

import (
	"log"
	"net"
	"sync"

	"github.com/korovkin/limiter"
	"github.com/ykpon/dnsbl-checker/lib/dnsbl"
)

// Hosts ...
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

func worker(ip string, wg *sync.WaitGroup) {
	defer wg.Done()
}
```
