// Package wave define a simple client for the wave algorithm
package wave

import (
	"client/client"
	"fmt"
	"prr.configuration/config"
)

// StartClient sends to all servers the execution signal.
// Wait a reception confirmation otherwise retry until receiving a confirmation
func StartClient() {

	// For all servers
	for _, v := range config.GetServers() {
		srvAddr := fmt.Sprintf("%v:%v", v.Ip, v.Port)
		client.SendUntilAck(client.SignalExec, srvAddr)

	}
}
