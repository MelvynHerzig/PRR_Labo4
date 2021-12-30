// Package main init config, reads file to know which algorithm is being executed and start the corresponding client.
package main

import (
	"client/client/probe"
	"client/client/wave"
	"log"
	"os"
	"prr.configuration/config"
	"strconv"
)

// main reads config and start the corresponding client
func main() {

	config.InitSimple("../config.json")

	switch config.GetAlgorithm() {

	case config.AlgoWave:
		wave.StartClient()

	case config.AlgoProbe:
		// Getting arg and simple check
		if len(os.Args) < 2 {
			log.Fatal("Usage: <server number> from [0, M[ where M is the server count")
		}

		srvNum, err := strconv.ParseUint(os.Args[1], 10, 0)
		if err != nil {
			log.Fatal("Argument is not an unsigned")
		}

		probe.StartClient(uint(srvNum))
	}
}
