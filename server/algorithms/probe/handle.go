// Package probe implements the logic for the probe and echo algorithm.
// It implements serialization/deserialization logic and handlers for incoming message
package probe

import (
	"fmt"
	"prr.configuration/config"
	"server/algorithms/common"
	"server/debug"
	"server/network"
)

// Handle is a function to start in a goroutine. It receives the message from the network process.
// It either starts the shortest path resolution if is not already running or transfer message to receiveMessage function.
func Handle() {

	// Adjacency matrix, updated at each echo and probe
	// Adjacency matrix i at topologies [i] is the matrix used when the server i handle a client demand
	// Reset after each end.
	var topologies = make([][][]bool, common.ServerCount)
	for i := uint(0); i < common.ServerCount; i++ {
		topologies[i] = common.ComputeBaseTopology()
	}

	// Array that holds the response count to get before sending echo back to parent for a given message (id is the key)
	// ids[i] is expected amount of answers where m.id = i
	var ids = make([]uint, common.ServerCount)

	// Array that holds the parent to send the echo when all response are received. parents[i] is parent to answer
	// where m.id = i
	var parents = make([]uint, common.ServerCount)

	for {
		select {
		// On message
		case data := <- network.DataChan:

			// If start execution signal
			if data.Message == network.SignalExec {

				if !common.Running {
					common.Running = true
					startExecution(&ids, &parents, &topologies[common.LocalNumber])
				}

			// Else if it's a message from another server, it's a probe/echo message, so we handle it.
			} else if config.IsServerIP(data.Sender) {

				// Sometimes the process has not yet started and already received waves. If this is the case, we
				// store the message until the process has started.
				m := deserialize(data.Message)
				debug.LogReceive(data.Message + " from " + fmt.Sprintf("%v", m.src))
				switch m.mType {
				case TypeEcho :
					handleEcho(&m, &ids, &topologies[m.id], &parents)
				case TypeProbe :
					handleProbe(&m, &ids, &topologies[m.id], &parents)
				}

			}
		}
	}
}
