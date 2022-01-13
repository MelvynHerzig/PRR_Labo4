package probe

import (
	"fmt"
	"server/algorithms/common"
	"server/debug"
	"server/network"
)

// maxUint is used when no filter is needed when using sendToNeighborsButException (in sendToNeighbors)
const maxUint = ^uint(0)

// startExecution is launched when a new shortest path search is asked. One at a time.
// Sends a probe to every neighbor. Set expected answer.
func startExecution(ids *[]uint, parents *[]uint, topology *[][]bool) {
	// Sending probe to all neighbors
	sendToNeighbors(&message{mType: TypeProbe, id: common.LocalNumber, src: common.LocalNumber, topology: *topology})

	(*ids)[common.LocalNumber] = common.NeighborsCount
	(*parents)[common.LocalNumber] = common.LocalNumber
}

// handleProbe handles a probe (message m) updates array ids of response count left, adjacency matrix and
// first parent for a given request id
func handleProbe(m *message, ids *[]uint, topology *[][]bool, parents *[]uint) {

	// if response left count is existing, we handle like an echo (circular graph)
	if (*ids)[m.id] > 0 {
		handleEcho(m, ids, topology, parents)

	// else first time we receive prob with this id
	} else {
		common.UpdateTopology(topology, &m.topology)
		sendToNeighborsButException(&message{mType: TypeProbe, id:m.id, src: common.LocalNumber, topology: *topology}, m.src)
		(*ids)[m.id] = common.NeighborsCount - 1
		(*parents)[m.id] = m.src

		// If we have no children we send directly an echo back
		checkAndEchoIfOver(m, ids, topology, parents)
	}
}

// handleEcho decrease needed answers count left for the given message id and checks if over
func handleEcho(m *message, ids *[]uint, topology *[][]bool, parents *[]uint) {
	common.UpdateTopology(topology, &m.topology)

	(*ids)[m.id] -= 1

	checkAndEchoIfOver(m, ids, topology, parents)
}

// checkAndEchoIfOver if answers count hits 0, send an echo to original parent.
func checkAndEchoIfOver(m *message, ids *[]uint, topology *[][]bool, parents *[]uint) {
	if (*ids)[m.id] == 0 {

		str := serialize(&message{mType: TypeEcho, id:m.id, src: common.LocalNumber, topology: *topology})

		if (*parents)[m.id] != common.LocalNumber {
			debug.LogSend(str + " to " + fmt.Sprintf("%v", (*parents)[m.id]))
			network.SendToServer(str, (*parents)[m.id])
		}

		common.PrintSP(topology, m.id)
		common.Running = false
		*topology = common.ComputeBaseTopology() // Reset for a new iteration
	}
}

// sendToNeighbors sends msg m to all neighbors. Used in startExecution.
func sendToNeighbors(m *message) {
	// We are using maxUint as exception because no node will ever get this id normally.
	sendToNeighborsButException(m, maxUint)
}

// sendToNeighborsButException sends msg m to all neighbors but the exception. Used in startExecution.
func sendToNeighborsButException(m *message, exception uint){
	for i := range common.IsMyNeighbors {
		if  common.IsMyNeighbors[i] && uint(i) != exception{
			str := serialize(m)
			debug.LogSend(str + " to " + fmt.Sprintf("%v", i))
			network.SendToServer(str, uint(i))
		}
	}
}