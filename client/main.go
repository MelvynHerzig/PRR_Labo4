// Package main init config, reads file to know which algorithm is being executed and start the corresponding client.
package main

import (
	"client/algorithms/wave"
	"prr.configuration/config"
)

// main reads config and start the corresponding client
func main() {

	config.InitSimple("../config.json")

	switch config.GetAlgorithm() {

	case config.AlgoWave:
		wave.StartClient()
	case config.AlgoProbe:
	}
}
