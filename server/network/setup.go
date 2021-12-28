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
// confirmUp is the message that a server send when a receive a message from another server during the network build
const confirmUp = "OK"

// startSignal is the message sent from parent to children in order to signal that all servers are online, and
// they can start to handle clients.
const startSignal = "GO"

// Array of all neighbors connections.
var neighborsConn [] net.Conn

// WaitNetwork is a blocking function that wait for network initialisation.
// In a first approach, we image a fake tree graph where a node m is connected to m-1.
// With this approach, the biggest node will be the root and there will always be only one branch.
// For example, if our network contains 5 node, the tree would be 5 -> 4 -> 3 -> 2 -> 1.
// Each node will wait until his child is ready. When the root get his child answer, it will send a "GO" message
// that means that the servers can start.
func WaitNetwork(socket *net.UDPAddr) {

	fmt.Println("STARTING) Connecting to neighbors...")

	localNumber := config.GetLocalServerNumber()
	neighborsId := config.GetNeighbors(localNumber)
	neighborsConn = make([]net.Conn, len(neighborsId))

	// Sorting neighbors ids
	sort.Slice(neighborsId, func(i, j int) bool { return neighborsId[i] < neighborsId[j] })

	// Waiting for all servers to come online.
	waitForChild(localNumber, socket)
	waitForStartSignal(localNumber)


	// Step 2
	//servers := config.GetServers()

}


func waitForChild(localNumber uint, socket *net.UDPAddr) {

	if localNumber > 0 {

		conn, err := net.ListenUDP("udp", socket)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		fakeChildInfo := config.GetServerById(localNumber - 1)

		buf := make([]byte, 64)
		var n = 0
		var cliAddr *net.UDPAddr = nil
		var receivedFromChild = false

		// Read with timeout
		for n == 0 || string(buf[0:n-1]) != confirmUp || !receivedFromChild {

			n, cliAddr, err = conn.ReadFromUDP(buf)
			if err != nil {
				log.Fatal(err)
			}

			receivedFromChild =  strings.Split(cliAddr.String(), ":")[0] == fakeChildInfo.Ip
		}
	}
}

func waitForStartSignal(localNumber uint) {
	notifyParentAndWaitStart(localNumber)
	sendStartToChild(localNumber)
}

func notifyParentAndWaitStart(localNumber uint) {
	if localNumber < uint(len(config.GetServers()) - 1) {

		var conn net.Conn
		var err error

		fakeParentInfo := config.GetServerById(localNumber + 1)
		fakeParentAddr := fmt.Sprintf("%v:%v", fakeParentInfo.Ip, fakeParentInfo.Port)

		conn, err = net.Dial("udp", fakeParentAddr)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		fmt.Fprintln(conn, confirmUp)

		buf := make([]byte, 64)
		var n = 0

		for  n == 0 || string(buf[0:n-1]) != startSignal{
			// set read timeout
			err := conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			if err != nil {
				log.Fatal(err)
			}

			n, err = conn.Read(buf)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					fmt.Fprintln(conn, confirmUp)
				} else {
					log.Fatal(err)
				}
			}
		}
	}
}

func sendStartToChild(localNumber uint) {
	if localNumber > 0 {

		var conn net.Conn
		var err error

		fakeChildInfo := config.GetServerById(localNumber - 1)
		fakeChildAddr := fmt.Sprintf("%v:%v", fakeChildInfo.Ip, fakeChildInfo.Port)

		conn, err = net.Dial("udp", fakeChildAddr)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		fmt.Fprintln(conn, startSignal)
	}
}

