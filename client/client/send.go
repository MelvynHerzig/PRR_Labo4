// Package client implements a methode that start a start execution message to a given server.
package client

import (
	"fmt"
	"log"
	"net"
	"time"
)

// SignalExec is the signal to notify servers that they must start the shortest path research
const SignalExec = "EXEC"

// signalReceived is the message that a server send when it handles the signalExec
const signalReceived = "ACK"

// SendUntilAck sends the message msg with udp at address addr. It wait a "ACK" otherwise retry after 1 second.
func SendUntilAck(msg string, addr string) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	// Signal to exec
	fmt.Fprintln(conn, msg)

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
				fmt.Fprintln(conn, msg)
			} else {
				log.Fatal(err)
			}
		}
	}

	conn.Close()
}
