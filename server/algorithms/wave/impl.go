package wave

import (
	"server/algorithms/common"
	"server/debug"
	"server/network"
)


// reset creates and return the base structures to the shortest path resolution, such has:
// the topology, the array of active neighbors and the active neighbors count.
func reset() ([][]bool, []bool, uint){

	var topology [][]bool
	var activeNeighbors []bool
	var activeNeighborsCount uint

	topology = common.ComputeBaseTopology()

	// Set active neighbors
	activeNeighbors = make([]bool, common.ServerCount)
	for i := uint(0); i < common.ServerCount; i++{
		activeNeighbors[i] = common.IsMyNeighbors[i]
	}

	// Set the active neighbors count
	activeNeighborsCount = common.NeighborsCount

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
		sendToActiveNeighbors(& message{topology: topology, src: common.LocalNumber, wave: wave, active: true}, &activeNeighbors)

		// Collecting wave response
		for i := uint(0); i < activeNeighborsCount; i++ {
			m := receiveMessage(wave)
			common.UpdateTopology(&topology, &m.topology)
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
	sendToActiveNeighbors(& message{topology: topology, src: common.LocalNumber, wave: wave, active: false}, &activeNeighbors)
	for i := uint(0); i < activeNeighborsCount; i++ {
		_ = receiveMessage(wave)
	}

	// Printing result
	common.PrintSP(&topology, common.LocalNumber)

	// Setting to false in order to restart if necessary
	common.Running = false
}

//  sendToActiveNeighbors sends the message m to all active neighbors in activeNeighbors
func sendToActiveNeighbors(m *message, activeNeighbors *[]bool) {
	for i := range common.IsMyNeighbors {
		if  common.IsMyNeighbors[i] && (*activeNeighbors)[i] {
			str := serialize(m)
			debug.LogSend(str)
			network.SendToServer(str, uint(i))
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