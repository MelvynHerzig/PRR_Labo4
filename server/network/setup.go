package network

import (
	"fmt"
	"log"
	"net"
	"prr.configuration/config"
	"sort"
	"strings"
	"time"
)

const (
	// signalUp is the message that a server send to his parent to say "Hi im up and ready, as well as my children"
	signalUp = "OK"

	// signalStart is the message that his emitted (initially by the root) from a parent to his child
	// to say "You can start handling clients son."
	signalStart = "GO"

	// signalReceived is the message to ack a received message. Its is mainly used during this phase because,
	// at this point, when a server sends a message, he is not sur that the receiver is up and running.
	signalReceived = "ACK"
)

// Array of all neighbors connections.
var neighborsConn [] net.Conn

// WaitNetwork is a blocking function that wait for network initialisation.
// In a first approach, we simulate a fake tree graph where a node m is connected to m-1.
// With this approach, the biggest node will be the root and there will always be only one branch.
// For example, if our network contains 5 node, the tree would be 4 -> 3 -> 2 -> 1 -> 0.
// Each node will wait until his child is ready. When the root get his child answer, it will send a "GO" message
// that means that the servers can start.
func WaitNetwork(socket *net.UDPAddr) {

	fmt.Println("STARTING) Connecting to neighbors...")

	localNumber := config.GetLocalServerNumber()

	// Starting listener UDP to get child's OK and parent GO
	conn, err := net.ListenUDP("udp", socket)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Step 1 waiting for every server to become up and ready
	// Waiting for child to come up and ready
	if localNumber > 0 {
		waitMessageFrom(conn, config.GetServerById(localNumber - 1), signalUp)
	}
	// Notify parent and wait for start signal
	if localNumber < uint(len(config.GetServers()) - 1) {
		sendMessageAndWaitConfirmation(config.GetServerById(localNumber + 1), signalUp)
		waitMessageFrom(conn, config.GetServerById(localNumber + 1), signalStart)
	}
	// Propagate to child start signal
	if localNumber > 0 {
		sendMessageAndWaitConfirmation(config.GetServerById(localNumber - 1), signalStart)
	}

	// Step 2 now we forget the fake tree. We get our real neighbors and make a connection
	neighborsId := config.GetNeighbors(localNumber)
	neighborsConn = make([]net.Conn, len(neighborsId))

	// Sorting neighbors ids
	sort.Slice(neighborsId, func(i, j int) bool { return neighborsId[i] < neighborsId[j] })

	fmt.Println("STARTING) All neighbors connected...")
}

// sendMessageAndWaitConfirmation sends the message msg to the distant server dtsSrv.
// The message is sent with UDP so, when the message is sent, it waits for a confirmation with a timeout of 1s.
func sendMessageAndWaitConfirmation(dstSrv *config.Server, msg string) {
	var conn net.Conn
	var err error

	dstAddr := fmt.Sprintf("%v:%v", dstSrv.Ip, dstSrv.Port)

	conn, err = net.Dial("udp", dstAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Fprintln(conn, msg)

	buf := make([]byte, 64)
	var n = 0

	for  n == 0 || string(buf[0:n]) != signalReceived{
		// set read timeout
		err := conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		if err != nil {
			log.Fatal(err)
		}

		n, err = conn.Read(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() { // Timeout, re-send message
				fmt.Fprintln(conn, msg)
			} else {
				log.Fatal(err)
			}
		}
	}
}

// waitMessageFrom used the connection conn to receive a message expectedMsg from a given sender expectedDstSrv.
// If the sender or the message is not the expected one, the message is ignored until the good one is received.
// Finally, a confirmation is sent to sender.
func waitMessageFrom(conn *net.UDPConn, expectedDstSrv *config.Server, expectedMsg string){
	buf := make([]byte, 64)
	var n = 0
	var cliAddr *net.UDPAddr = nil
	var err error
	var isExpectedSender = false

	// Reading incoming messages
	for n == 0 || string(buf[0:n-1]) != expectedMsg || !isExpectedSender {

		n, cliAddr, err = conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}

		isExpectedSender =  strings.Split(cliAddr.String(), ":")[0] == expectedDstSrv.Ip
	}

	// Confirm
	if _, err := conn.WriteTo([]byte(signalReceived), cliAddr); err != nil {
		log.Fatal(err)
	}
}

