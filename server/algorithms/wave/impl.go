package wave

import (
	"container/list"
	"fmt"
	"prr.configuration/config"
	"server/debug"
	"server/network"
)

// Those variables are initiated once when the server start in the init function. They are and must be never changed.

// localNumber is the number of the node represented by the current server
var localNumber uint

// serverCount is the number of servers in the network
var serverCount uint

// isMyNeighbors is an array of serverCount long which indicate who are the server neighbors.
// if the array[i] is true, the node i is a neighbor of the current server
var isMyNeighbors []bool

// neighborsCount is the number of neighbor for the current node/server.
var neighborsCount uint



// running indicates if the shortest path research has already started, se to true in the wave handler and reset to
// false in the function searchSP (max one at a time).
var running = false

// init is a function launched once by the handler. It instantiates the const variables such as:
// localNumber, serverCount, isMyNeighbors and neighborsCount. All these, are and must not be mutated.
func initVariables() {
	localNumber = config.GetLocalServerNumber()
	serverCount = uint(len(config.GetServers()))

	isMyNeighbors   = make([]bool, serverCount)
	neighborsCount = uint(len(config.GetNeighbors(localNumber)))
	for _, v := range config.GetNeighbors(localNumber) {
		isMyNeighbors[v] = true
	}
}

// reset creates and return the base structures to the shortest path resolution, such has:
// the topology, the array of active neighbors and the active neighbors count.
func reset() ([][]bool, []bool, uint){

	var topology [][]bool
	var activeNeighbors []bool
	var activeNeighborsCount uint

	// Size the topology according to the serverCount
	topology = make([][]bool, serverCount)
	for i := range topology {
		topology[i] = make([]bool, serverCount)
	}

	// Sets the base topology and the activeNeighbors according to the neighbors
	activeNeighbors = make([]bool, serverCount)
	for i := uint(0); i < serverCount; i++{
		topology[localNumber][i] = isMyNeighbors[i]
		activeNeighbors[i] = isMyNeighbors[i]
	}

	// Set the active neighbors count
	activeNeighborsCount = neighborsCount

	return topology, activeNeighbors, activeNeighborsCount
}

// searchSP is a function launched in a goroutine by the handler when he receives a "EXEC" message.
// There is only one goroutine at a time. It resolves the shortest path for the existing network.
func searchSP() {
	topology, activeNeighbors, activeNeighborsCount := reset()

	// Number of the wave. incremented before each sendToActiveNeighbors
	wave := uint(0)

	for {
		wave++
		sendToActiveNeighbors(& message{topology: topology, src: localNumber, wave: wave, active: true}, &activeNeighbors)

		// Collecting wave response
		for i := uint(0); i < neighborsCount; i++ {
			m := receiveMessage(wave)
			updateTopology(&topology, &m.topology)
			activeNeighbors[m.src] = m.active
			if !m.active {
				activeNeighborsCount--
			}
		}

		// Do we know a path to every node ?
		if isTopologyComplete(&topology) {
			break
		}
	}

	// Sending final wave
	wave++
	sendToActiveNeighbors(& message{topology: topology, src: localNumber, wave: wave, active: false}, &activeNeighbors)
	for i := uint(0); i < activeNeighborsCount; i++ {
		_ = receiveMessage(wave)
	}

	// Printing result
	printSP(&topology)

	// Setting to false in order to restart if necessary
	running = false
}

//  sendToActiveNeighbors sends the message m to all active neighbors in activeNeighbors
func sendToActiveNeighbors(m *message, activeNeighbors *[]bool) {
	for i := range isMyNeighbors {
		if isMyNeighbors[i] && (*activeNeighbors)[i] {
			str := serialize(m)
			debug.LogSend(str)
			network.SendToServer(str, uint(i))
		}
	}
}

// updateTopology updates the topoply currenTopo with the new topology newTopo with a logical OR
func updateTopology(currentTopo *[][]bool, newTopo *[][]bool){
	for i := range *currentTopo {
		for j := range (*currentTopo)[i] {
			if (*currentTopo)[i][j] == false {
				(*currentTopo)[i][j] = (*newTopo)[i][j]
			}
		}
	}
}

// isTopologyComplete return false if one line contains only false values else true.
func isTopologyComplete(currentTopo *[][]bool) bool {
	for i := range *currentTopo {
		isLineEmpty := true

		for j := range (*currentTopo)[i] {
			if (*currentTopo)[i][j] {
				isLineEmpty = false
				break
			}
		}

		if isLineEmpty {
			return false
		}
	}

	return true
}

// Print the shortest paths based on the final topology finalTopo.
func printSP(finalTopo *[][]bool) {

	detailedSP := computeSP(finalTopo)

	for i, SP := range detailedSP {
		fmt.Printf("Shortest path to %v, length: %v, Path:", i, len(SP) - 1)

		for i, node := range SP {
			if i != 0 {
				fmt.Printf(" ->")
			}
			fmt.Printf(" %v", node)
		}
		fmt.Printf("\n")
	}
}

// computeSP makes a breadth first search in the final topology in order to find the shortest path to every node
// Returns a 2D slice of the shortest paths. Row[i] = the shortest path to reach node i, Column[i] = ith node to visit for
// the given shortest path
func computeSP(finalTopo *[][]bool) [][]uint {

	// 2D slice of the shortest paths. Row[i] = the shortest path to reach node i, Column[i] = ith node to visit for
	// the given the shortest path
	detailedSP := make([][]uint, serverCount)
	for i := range *finalTopo  {
		detailedSP[i] = make([]uint, 0, serverCount)
	}

	// The shortest path to go the current node, is the current node.
	detailedSP[localNumber] = append(detailedSP[localNumber], localNumber)

	// Then, we start a breadth first search
	var queue = list.New()
	// Beginning with neighbors
	for i := uint(0); i < serverCount; i++ {
		if isMyNeighbors[i] {
			queue.PushBack(i)
			// Inserting local node and neighbors
			detailedSP[i] = append(detailedSP[i], localNumber)
			detailedSP[i] = append(detailedSP[i], i)
		}
	}

	for queue.Len() != 0 {
		front := queue.Front()
		queue.Remove(front)
		node := front.Value.(uint)

		for i := uint(0); i < serverCount; i++ {
			if i != localNumber && (*finalTopo )[node][i] && len(detailedSP[i]) == 0 {
				queue.PushBack(i)

				for _, v := range detailedSP[node] {
					detailedSP[i] = append(detailedSP[i], v)
				}
				detailedSP[i] = append(detailedSP[i], i)
			}
		}
	}
	return detailedSP
}