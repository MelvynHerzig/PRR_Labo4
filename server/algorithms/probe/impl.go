package probe

import (
	"fmt"
	"server/algorithms"
	"server/debug"
	"server/network"
)

// maxUint is used when no filter is needed when using sendToNeighborsButException (in sendToNeighbors)
const maxUint = ^uint(0)

// startExecution is launched when a new shortest path search is asked. One at a time.
func startExecution() (map[uint]uint, [][]uint) {
	ids := make(map[uint]uint)
	ids[algorithms.LocalNumber] = algorithms.NeighborsCount

	// Size the temporary known shortest paths
	temporaryKnownSp := make([][]uint, algorithms.ServerCount)
	for i := range temporaryKnownSp {
		temporaryKnownSp[i] = make([]uint, 0, algorithms.ServerCount)

		// Trivial, shortest path to our neighbors or ourselves
		if algorithms.IsMyNeighbors[i] {
			temporaryKnownSp[i] = append(temporaryKnownSp[i], algorithms.LocalNumber, uint(i))
		} else if uint(i) == algorithms.LocalNumber {
			temporaryKnownSp[i] = append(temporaryKnownSp[i], algorithms.LocalNumber)
		}
	}

	// Sending probe to all neighbors
	sendToNeighbors(&message{mType: TypeProbe, id:algorithms.LocalNumber, src: algorithms.LocalNumber, temporaryKnownSp: temporaryKnownSp})

	return ids, temporaryKnownSp
}

// handleProbe handles a probe (message m) updates map ids of response count left, 2D slices of know shortest path and
// first parent for a given request id
func handleProbe(m *message, ids *map[uint]uint, temporaryKnownSp *[][]uint, parents *map[uint]uint) {

	// if ids is nil, we are not the "root" so we need to set the base struct.
	if *ids == nil {
		*ids = make(map[uint]uint)
		*parents = make(map[uint]uint)
	}

	_, ok := (*ids)[m.id]
	// if response left count is existing, we handle like an echo (circular graph)
	if ok {
		handleEcho(m, ids, temporaryKnownSp, parents)

	// else first time we receive prob with this id
	} else {
		updateTemporaryShortestPaths(temporaryKnownSp, &m.temporaryKnownSp)
		sendToNeighborsButException(&message{mType: TypeProbe, id:m.id, src: algorithms.LocalNumber, temporaryKnownSp: *temporaryKnownSp}, m.src)
		(*ids)[m.id] = algorithms.NeighborsCount - 1
		(*parents)[m.id] = m.src

		// If we have no children we send directly an echo back
		checkAndEchoIfOver(m, ids, temporaryKnownSp, parents)
	}
}

// handleEcho decrease needed answers count left for the given message id and checks if over
func handleEcho(m *message, ids *map[uint]uint, temporaryKnownSp *[][]uint, parents *map[uint]uint) {

	updateTemporaryShortestPaths(temporaryKnownSp, &m.temporaryKnownSp)

	cnt, _ := (*ids)[m.id]
	(*ids)[m.id] = cnt - 1

	checkAndEchoIfOver(m, ids, temporaryKnownSp, parents)
}

// checkAndEchoIfOver if answers count hits 0, send an echo to original parent.
func checkAndEchoIfOver(m *message, ids *map[uint]uint, temporaryKnownSp *[][]uint, parents *map[uint]uint) {
	if (*ids)[m.id] == 0 {
		str := serialize(&message{mType: TypeEcho, id:m.id, src: algorithms.LocalNumber, temporaryKnownSp: *temporaryKnownSp})

		if parent, ok := (*parents)[m.id]; ok {
			debug.LogSend(str + " to " + fmt.Sprintf("%v", parent))
			network.SendToServer(str, parent)
		}

		printSP(temporaryKnownSp)
		delete(*ids, m.id)
		delete(*parents, m.id)
		*temporaryKnownSp = nil
		algorithms.Running = false
	}
}

// sendToNeighbors sends msg m to all neighbors. Used in startExecution.
func sendToNeighbors(m *message) {
	// We are using maxUint as exception because no node will ever get this id normally.
	sendToNeighborsButException(m, maxUint)
}

// sendToNeighborsButException sends msg m to all neighbors but the exception. Used in startExecution.
func sendToNeighborsButException(m *message, exception uint){
	for i := range algorithms.IsMyNeighbors {
		if  algorithms.IsMyNeighbors[i] && uint(i) != exception{
			str := serialize(m)
			debug.LogSend(str + " to " + fmt.Sprintf("%v", i))
			network.SendToServer(str, uint(i))
		}
	}
}

// updateTemporaryShortestPaths updates the local list of shortest paths localSP if shortest paths are found in the
// received list receivedSP
func updateTemporaryShortestPaths(localSP *[][]uint, receivedSP *[][]uint) {

	// if we have no SP we copy the received list.
	if *localSP == nil {
		*localSP = *receivedSP

		// Trying to add our neighbors
		for i := uint(0); i < algorithms.ServerCount ; i++ {
			// If no path exist or if our path is shortest
			if  algorithms.IsMyNeighbors[i] &&
				(len((*localSP)[i]) == 0 || len((*localSP)[algorithms.LocalNumber]) + 1 <  len((*localSP)[i])) {
				(*localSP)[i] = make([]uint, len((*localSP)[algorithms.LocalNumber]))
				copy((*localSP)[i], (*localSP)[algorithms.LocalNumber])
				(*localSP)[i] = append((*localSP)[i], i)
			}
		}
	} else {
		// For each path, if our path is longer that the received one (if he exists) we remplace our path.
		for i := range *receivedSP {
			if len( (*localSP)[i]) == 0 ||  ( len((*localSP)[i]) > len((*receivedSP)[i]) && len((*receivedSP)[i]) != 0 ) {
				(*localSP)[i] = (*receivedSP)[i]
			}
		}
	}
}

// Print the shortest paths
func printSP(finalSP *[][]uint) {

	for i, SP := range *finalSP {
		if len(SP) > 0 {
			fmt.Printf("Shortest path to %v, length: %v, Path:", i, len(SP) - 1)
			for i, node := range SP {
				if i != 0 {
					fmt.Printf(" ->")
				}
				fmt.Printf(" %v", node)
			}
		} else {
			fmt.Printf("unknown")
		}
		fmt.Printf("\n")
	}
}