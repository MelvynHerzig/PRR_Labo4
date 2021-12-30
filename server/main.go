// Package main Gets program args, launch setup network process and algorithm process
package main

import (
	"log"
	"net"
	"os"
	"prr.configuration/config"
	"server/algorithms/wave"
	"server/network"
	"strconv"
)

// main gets program args, launch setup network process and algorithm process
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

	// Resolving UDP Server.
	port := config.GetServerById(config.GetLocalServerNumber()).Port
	addr, _ := net.ResolveUDPAddr("udp", "localhost:" + strconv.FormatUint(uint64(port), 10))

	network.WaitNetwork(addr)

	// Starting network process
	go network.Handle(addr)

	// Starting algorithm process
	switch config.GetAlgorithm() {
	case config.AlgoWave:
		wave.Handle()
	case config.AlgoProbe:
	}
}