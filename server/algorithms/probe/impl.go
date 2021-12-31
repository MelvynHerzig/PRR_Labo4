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
func startExecution() (map[uint]uint, [][]bool, map[uint]uint) {
	ids := make(map[uint]uint)
	ids[common.LocalNumber] = common.NeighborsCount

	topology := common.ComputeBaseTopology()

	parents := make(map[uint]uint)

	// Sending probe to all neighbors
	sendToNeighbors(&message{mType: TypeProbe, id: common.LocalNumber, src: common.LocalNumber, topology: topology})

	return ids, topology, parents
}

// handleProbe handles a probe (message m) updates map ids of response count left, adjacency matrix and
// first parent for a given request id
func handleProbe(m *message, ids *map[uint]uint, topology *[][]bool, parents *map[uint]uint) {

	// if ids is nil, we are not the "root" so we need to set the base struct.
	if *ids == nil {
		*ids = make(map[uint]uint)
		*parents = make(map[uint]uint)
	}

	_, ok := (*ids)[m.id]
	// if response left count is existing, we handle like an echo (circular graph)
	if ok {
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
func handleEcho(m *message, ids *map[uint]uint, topology *[][]bool, parents *map[uint]uint) {

	common.UpdateTopology(topology, &m.topology)

	cnt, _ := (*ids)[m.id]
	(*ids)[m.id] = cnt - 1

	checkAndEchoIfOver(m, ids, topology, parents)
}

// checkAndEchoIfOver if answers count hits 0, send an echo to original parent.
func checkAndEchoIfOver(m *message, ids *map[uint]uint, topology *[][]bool, parents *map[uint]uint) {
	if (*ids)[m.id] == 0 {
		str := serialize(&message{mType: TypeEcho, id:m.id, src: common.LocalNumber, topology: *topology})

		if parent, ok := (*parents)[m.id]; ok {
			debug.LogSend(str + " to " + fmt.Sprintf("%v", parent))
			network.SendToServer(str, parent)
		}

		common.PrintSP(topology, m.id)
		delete(*ids, m.id)
		delete(*parents, m.id)
		*topology = nil
		common.Running = false
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