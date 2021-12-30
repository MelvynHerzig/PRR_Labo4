// Package probe implements the logic for the probe and echo algorithm.
// It implements serialization/deserialization logic and handlers for incoming message
package probe

import (
	"fmt"
	"prr.configuration/config"
	"server/algorithms"
	"server/debug"
	"server/network"
)

// Handle is a function to start in a goroutine. It receives the message from the network process.
// It either starts the shortest path resolution if is not already running or transfer message to receiveMessage function.
func Handle() {

	// temporary known the shortest path, updated at each echo and probe
	var temporaryKnownSp [][]uint

	// Map that holds the response count to get before sending echo back to parent for a given message (id is the key)
	var ids map[uint]uint

	// Map that holds the parent to send the echo when all response are received
	var parents map[uint]uint

	for {
		select {
		// On message
		case data := <- network.DataChan:

			// If start execution signal and not already running
			if data.Message == network.SignalExec && !algorithms.Running  {
				algorithms.Running = true
				ids, temporaryKnownSp = startExecution()

			// Else if it's a message from another server, it's a probe/echo message, so we handle it.
			} else if config.IsServerIP(data.Sender) {

				// Sometimes the process has not yet started and already received waves. If this is the case, we
				// store the message until the process has started.
				m := deserialize(data.Message)
				debug.LogReceive(data.Message + " from " + fmt.Sprintf("%v", m.src))
				switch m.mType {
				case TypeEcho :
					handleEcho(&m, &ids, &temporaryKnownSp, &parents)
				case TypeProbe :
					handleProbe(&m, &ids, &temporaryKnownSp, &parents)
				}

			}
		}
	}
}
