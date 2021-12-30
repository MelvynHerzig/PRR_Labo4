package wave

import (
	"container/list"
	"fmt"
	"server/algorithms"
	"server/debug"
	"server/network"
)


// reset creates and return the base structures to the shortest path resolution, such has:
// the topology, the array of active neighbors and the active neighbors count.
func reset() ([][]bool, []bool, uint){

	var topology [][]bool
	var activeNeighbors []bool
	var activeNeighborsCount uint

	// Size the topology according to the serverCount
	topology = make([][]bool, algorithms.ServerCount)
	for i := range topology {
		topology[i] = make([]bool, algorithms.ServerCount)
	}

	// Sets the base topology and the activeNeighbors according to the neighbors
	activeNeighbors = make([]bool, algorithms.ServerCount)
	for i := uint(0); i < algorithms.ServerCount; i++{
		topology[algorithms.LocalNumber][i] = algorithms.IsMyNeighbors[i]
		activeNeighbors[i] = algorithms.IsMyNeighbors[i]
	}

	// Set the active neighbors count
	activeNeighborsCount = algorithms.NeighborsCount

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
		sendToActiveNeighbors(& message{topology: topology, src: algorithms.LocalNumber, wave: wave, active: true}, &activeNeighbors)

		// Collecting wave response
		for i := uint(0); i < algorithms.NeighborsCount; i++ {
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
	sendToActiveNeighbors(& message{topology: topology, src: algorithms.LocalNumber, wave: wave, active: false}, &activeNeighbors)
	for i := uint(0); i < activeNeighborsCount; i++ {
		_ = receiveMessage(wave)
	}

	// Printing result
	printSP(&topology)

	// Setting to false in order to restart if necessary
	algorithms.Running = false
}

//  sendToActiveNeighbors sends the message m to all active neighbors in activeNeighbors
func sendToActiveNeighbors(m *message, activeNeighbors *[]bool) {
	for i := range algorithms.IsMyNeighbors {
		if  algorithms.IsMyNeighbors[i] && (*activeNeighbors)[i] {
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
	detailedSP := make([][]uint,  algorithms.ServerCount)
	for i := range *finalTopo  {
		detailedSP[i] = make([]uint, 0,  algorithms.ServerCount)
	}

	// The shortest path to go the current node, is the current node.
	detailedSP[algorithms.LocalNumber] = append(detailedSP[algorithms.LocalNumber], algorithms.LocalNumber)

	// Then, we start a breadth first search
	var queue = list.New()
	// Beginning with neighbors
	for i := uint(0); i < algorithms.ServerCount; i++ {
		if algorithms.IsMyNeighbors[i] {
			queue.PushBack(i)
			// Inserting local node and neighbors
			detailedSP[i] = append(detailedSP[i], algorithms.LocalNumber)
			detailedSP[i] = append(detailedSP[i], i)
		}
	}

	for queue.Len() != 0 {
		front := queue.Front()
		queue.Remove(front)
		node := front.Value.(uint)

		for i := uint(0); i < algorithms.ServerCount; i++ {
			if i != algorithms.LocalNumber && (*finalTopo )[node][i] && len(detailedSP[i]) == 0 {
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