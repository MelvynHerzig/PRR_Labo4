package main

import (
	"log"
	"net"
	"os"
	"prr.configuration/config"
	"server/network"
	"strconv"
)

func main() {

	// Getting program arg (server number).
	argsLen := len(os.Args)
	if argsLen < 2 {
		log.Fatal("Usage: <server number>")
	}

	noServ, errNoServ   := strconv.ParseUint(os.Args[1], 10, 0)
	if  errNoServ != nil {
		log.Fatal("Invalid parameter. Must be <no serveur>")
	}

	// Init configuration
	config.Init("../config.json", uint(noServ))

	if noServ < 0 || noServ >= uint64(len(config.GetServers())) {
		log.Fatal("Server number is an integer between [0, servers count [")
	}

	// Opening UDP Server.
	port := config.GetServerById(config.GetLocalServerNumber()).Port
	addr, _ := net.ResolveUDPAddr("udp", "localhost:" + strconv.FormatUint(uint64(port), 10))

	network.WaitNetwork(addr)
}