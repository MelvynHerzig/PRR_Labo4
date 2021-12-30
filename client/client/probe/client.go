// Package probe define a simple client for the probe algorithm
package probe

import (
	"client/client"
	"fmt"
	"log"
	"prr.configuration/config"
)

// StartClient sends to the server of number srvId the execution signal.
// Wait a reception confirmation otherwise retry until receiving a confirmation
func StartClient(srvNum uint) {

	// For all servers
	if srvNum < 0 || srvNum >= uint(len(config.GetServers())) {
		log.Fatal("Invalid server number must be [0, M[ where M is the server count")
	}
	server := config.GetServerById(srvNum)
	srvAddr := fmt.Sprintf("%v:%v", server.Ip, server.Port)
	client.SendUntilAck(client.SignalExec, srvAddr)
}