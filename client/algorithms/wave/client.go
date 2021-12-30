// Package wave define a simple client for the wave algorithm
package wave

import (
	"fmt"
	"log"
	"net"
	"prr.configuration/config"
	"time"
)

// signalExec is the signal to notify servers that they must start the shortest path research
const signalExec = "EXEC"

// signalReceived is the message that a server send when it handles the signalExec
const signalReceived = "ACK"

// StartClient sends to all servers the execution signal.
// Wait a reception confirmation otherwise retry until receiving a confirmation
func StartClient() {

	// For all servers
	for _, v := range config.GetServers() {
		srvAddr := fmt.Sprintf("%v:%v", v.Ip, v.Port)
		conn, err := net.Dial("udp", srvAddr)
		if err != nil {
			log.Fatal(err)
		}

		// Signal to exec
		fmt.Fprintln(conn, signalExec)

		buf := make([]byte, 64)
		var n = 0

		for  n == 0 || string(buf[0:n]) != signalReceived{

			// set read timeout
			err := conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			if err != nil {
				log.Fatal(err)
			}

			// Try to read ack
			n, err = conn.Read(buf)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() { // Timeout, re-send message
					fmt.Fprintln(conn, signalExec)
				} else {
					log.Fatal(err)
				}
			}
		}

		conn.Close()
	}
}
