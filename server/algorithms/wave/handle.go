// Package wave implements the logic for the waves algorithm.
// It implements serialization/deserialization logic and handlers for incoming message
package wave

import (
	"container/list"
	"prr.configuration/config"
	"server/algorithms/common"
	"server/debug"
	"server/network"
)

// messageChan is a channel to transfer message from handler to receiveMessage function.
var messageChan = make(chan *message)

// cachedMessages is a cache to store yet unneeded message but useful for the next wave.
var cachedMessages = list.New()

// Handle is a function to start in a goroutine. It receives the message from the network process.
// It either starts the shortest path resolution if is not already running or transfer message to receiveMessage function.
func Handle() {

	for {
		select {
			// On message
			case data := <- network.DataChan:

				// If start execution signal and not already running
				if data.Message == network.SignalExec && !common.Running {
					common.Running = true
					go searchSP()

				// Else if it's a message from another server, it's a wave message, so we handle it.
				} else if config.IsServerIP(data.Sender) && common.Running {

					// Sometimes the process has not yet started and already receive waves. If this is the case, we
					// store the message until the process has started.
					m := deserialize(data.Message)
					if common.Running {
						messageChan <- &m
					} else {
						cachedMessages.PushBack(&m)
					}

				}
		}
	}
}

// receiveMessage is a function used by the shortest path search function. It returns messages for a given wave number.
// Might be blocking if nothing corresponding has been yet received.
func receiveMessage(wave uint) message {

	// Checking if a cached message correspond to the current required wave
	for e := cachedMessages.Front(); e != nil; e = e.Next() {
		v := e.Value
		m, _ := v.(*message)

		if m.wave == wave {
			cachedMessages.Remove(e)
			debug.LogReceive(serialize(m))
			return *m
		}
	}

	// No cached message correspond to the wave, this kind of message has not been received yet.
	// So, we wait that the handler transmit a new message to check if it corresponds
	for {
		m := <- messageChan
		if m.wave != wave {
			cachedMessages.PushBack(m)
		} else {
			debug.LogReceive(serialize(m))
			return *m
		}
	}
}
