package wave

import (
	"container/list"
	"fmt"
	"prr.configuration/config"
	"server/network"
)

var localNumber uint
var serverCount uint
var isMyNeighbors []bool
var activeNeighbors []bool
var activeNeighborsCount uint
var neighborsCount uint

var topology [][]bool

var running = false

func Init() {
	localNumber = config.GetLocalServerNumber()
	serverCount = uint(len(config.GetServers()))

	isMyNeighbors   = make([]bool, serverCount)
	activeNeighbors = make([]bool, serverCount)
	neighborsCount = uint(len(config.GetNeighbors(localNumber)))
	for _, v := range config.GetNeighbors(localNumber) {
		isMyNeighbors[v] = true
	}
}

func reset() {
	topology = make([][]bool, serverCount)
	for i := range topology {
		topology[i] = make([]bool, serverCount)
	}

	for i := uint(0); i < serverCount; i++{
		topology[localNumber][i] = isMyNeighbors[i]
		activeNeighbors[i] = isMyNeighbors[i]
	}

	activeNeighborsCount = neighborsCount
}

func searchSP() {
	reset()

	for {
		sendToNeighbors(message{topology: topology, src: localNumber, active: true})
		for i := uint(0); i < neighborsCount; i++ {
			m := <- messageChan
			updateTopology(m.topology)
			activeNeighbors[m.src] = m.active
			if !m.active {
				activeNeighborsCount--
			}
		}
		if isTopologyComplete() {
			break
		}
	}

	sendToActiveNeighbors(message{topology: topology, src: localNumber, active: false})
	for i := uint(0); i < activeNeighborsCount; i++ {
		_ = <- messageChan
	}

	printSP()

	running = false
}

func sendToNeighbors(m message) {
	for i := range isMyNeighbors {
		if isMyNeighbors[i] {
			network.SendTo(serialize(m), uint(i))
		}
	}
}

func sendToActiveNeighbors(m message) {
	for i := range isMyNeighbors {
		if isMyNeighbors[i] && activeNeighbors[i] {
			network.SendTo(serialize(m), uint(i))
		}
	}
}

func updateTopology(newTopo [][]bool){
	for i := range topology {
		for j := range topology[i] {
			if topology[i][j] == false {
				topology[i][j] = newTopo[i][j]
			}
		}
	}
}

func isTopologyComplete() bool {
	for i := range topology {
		isLineEmpty := true

		for j := range topology[i] {
			if topology[i][j] {
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

func printSP() {

	// 2D slice of the shortest paths
	detailedSP := make([][]uint, serverCount)
	for i := range topology {
		detailedSP[i] = make([]uint, 0, serverCount)
	}

	// Each shortest path begin from the current node
	for i := uint(0); i < serverCount; i++ {
		detailedSP[i] = append(detailedSP[i], localNumber + 1)
	}

	// Then, we start a breadth first search
	var queue = list.New()
	// Beginning with neighbors
	for i := uint(0); i < serverCount; i++ {
		if isMyNeighbors[i] {
			queue.PushBack(i)
			detailedSP[i] = append(detailedSP[i], i + 1)
		}
	}

	for queue.Len() != 0 {
		front := queue.Front()
		queue.Remove(front)
		node := front.Value.(uint)

		for i := uint(0); i < serverCount; i++ {
			if topology[node][i] && len(detailedSP[i]) != 1 {
				queue.PushBack(i)

				for _, v := range detailedSP[node] {
					detailedSP[i] = append(detailedSP[i], v)
				}
			}
		}
	}


	for _, SP := range detailedSP {
		fmt.Printf("Length: %v, Path:", len(SP) - 1)

		for i, node := range SP {
			if i != 0 {
				fmt.Printf(" ->")
			}
			fmt.Printf(" %v", node)
		}

		fmt.Printf("\n")
	}

}