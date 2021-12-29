package network

import (
	"fmt"
	"log"
	"net"
)

type ReceivedData struct{
	Message string
	Sender string
}

var DataChan = make(chan ReceivedData)

func SendTo(s string, dst uint) {
	fmt.Fprintln(neighbors[dst], s)
}

func Handle(socket *net.UDPAddr) {
	buf := make([]byte, 1024)
	var n = 0
	var cliAddr *net.UDPAddr = nil
	var err error

	conn, err := net.ListenUDP("udp", socket)
	for {

		n, cliAddr, err = conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}

		DataChan <- ReceivedData{Message: string(buf[0:n-1]), Sender : cliAddr.String() }
	}
}