package network

import (
	"fmt"
	"log"
	"net"
)

// SignalExec is that signal to start the algorithm execution. Every other message is considered to be server messages.
const SignalExec = "EXEC"

// ReceivedData is the struct defines the data received on the UDP socket. It is transmitted to who reads DataChan.
type ReceivedData struct{
	Message string
	Sender string
}

// DataChan is the chan that transfer the received messages as ReceivedData struct
var DataChan = make(chan *ReceivedData)

// SendToServer send the s string to the neighbor server with number dst.
func SendToServer(s string, dst uint) {
	fmt.Fprintln(neighbors[dst], s)
}

// Handle is a function to be started in a goroutine in order to receive the message on the UDP socket.
func Handle(socket *net.UDPAddr) {

	buf := make([]byte, 1024)
	var n = 0
	var cliAddr *net.UDPAddr = nil
	var err error

	conn, err := net.ListenUDP("udp", socket)
	fmt.Println("STARTING) Ready to listen incoming messages")

	for {
		n, cliAddr, err = conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}

		// If is a start execution -> confirm
		if string(buf[0:n-1]) == SignalExec {
			conn.WriteTo([]byte(signalReceived), cliAddr)
		}

		DataChan <- &ReceivedData{Message: string(buf[0:n-1]), Sender : cliAddr.String() }
	}
}