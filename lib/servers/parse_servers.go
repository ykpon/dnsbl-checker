package servers

import (
	"encoding/csv"
	"os"
)

// Server ...
type Server struct {
	Host string
	Name string
}

// GetServers ...
func GetServers() []Server {
	file, err := os.Open("resources/servers.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	var servers []Server

	for {
		record, e := reader.Read()
		if e != nil {
			break
		}

		servers = append(servers, Server{Name: record[0], Host: record[1]})
	}

	return servers
}
