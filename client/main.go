package main

import (
	"fmt"
	"log"
	"net"
	"prr.configuration/config"
)

func main() {
	config.InitSimple("../config.json")

	for _, v := range config.GetServers() {
		srvAddr := fmt.Sprintf("%v:%v", v.Ip, v.Port)
		conn, err := net.Dial("udp", srvAddr)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintln(conn, "START")
		conn.Close()
	}
}
