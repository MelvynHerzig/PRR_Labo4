// Package common define a common methods and attributes shared between probe and wave implementation
package common

import (
	"container/list"
	"fmt"
	"prr.configuration/config"
)

// Those variables are initiated once when the server start in the init function. They are and must be never changed.

// LocalNumber is the number of the node represented by the current server
var LocalNumber uint

// ServerCount is the number of servers in the network
var ServerCount uint

// IsMyNeighbors is an array of serverCount long which indicate who are the server neighbors.
// if the array[i] is true, the node i is a neighbor of the current server
var IsMyNeighbors []bool

// NeighborsCount is the number of neighbor for the current node/server.
var NeighborsCount uint

// Running indicates if the shortest path research has already started, se to true in the wave handler and reset to
// false by respective algorithms implementations. It only prevents from a local double execution, not from another node.
var Running = false

func InitVariables() {
	LocalNumber = config.GetLocalServerNumber()
	ServerCount = uint(len(config.GetServers()))

	IsMyNeighbors = make([]bool, ServerCount)
	NeighborsCount = uint(len(config.GetNeighbors(LocalNumber)))
	for _, v := range config.GetNeighbors(LocalNumber) {
		IsMyNeighbors[v] = true
	}
}

// ComputeBaseTopology compute the base adjacency matrix for the current node
func ComputeBaseTopology() [][]bool {
	// Size the topology according to the serverCount
	topology := make([][]bool, ServerCount)
	for i := range topology {
		topology[i] = make([]bool, ServerCount)
	}

	// Sets the base topology according to the neighbors
	for i := uint(0); i < ServerCount; i++{
		topology[LocalNumber][i] = IsMyNeighbors[i]
	}

	return topology
}

// UpdateTopology updates the topoply currenTopo with the new topology newTopo with a logical OR
func UpdateTopology(currentTopo *[][]bool, newTopo *[][]bool){
	if *currentTopo == nil {
		*currentTopo = ComputeBaseTopology()
	}

	for i := range *currentTopo {
		for j := range (*currentTopo)[i] {
			if (*currentTopo)[i][j] == false {
				(*currentTopo)[i][j] = (*newTopo)[i][j]
			}
		}
	}
}

// PrintSP the shortest paths based on the final topology finalTopo.
func PrintSP(finalTopo *[][]bool, start uint) {

	detailedSP := ComputeSP(finalTopo, start)

	for i, SP := range detailedSP {
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

// ComputeSP makes a breadth first search in the final topology in order to find the shortest path to every node from source.
// Returns a 2D slice of the shortest paths. Row[i] = the shortest path to reach node i, Column[i] = ith node to visit for
// the given shortest path
func ComputeSP(finalTopo *[][]bool, start uint) [][]uint {

	// 2D slice of the shortest paths. Row[i] = the shortest path to reach node i, Column[i] = ith node to visit for
	// the given the shortest path
	detailedSP := make([][]uint, ServerCount)
	for i := range *finalTopo  {
		detailedSP[i] = make([]uint, 0, ServerCount)
	}

	// The shortest path to go the current node, is the current node.
	detailedSP[start] = append(detailedSP[start], start)

	// Then, we start a breadth first search
	var queue = list.New()
	// Beginning with neighbors
	for i := uint(0); i < ServerCount; i++ {
		if (*finalTopo)[start][i] {
			queue.PushBack(i)
			// Inserting local node and neighbors
			detailedSP[i] = append(detailedSP[i], start)
			detailedSP[i] = append(detailedSP[i], i)
		}
	}

	for queue.Len() != 0 {
		front := queue.Front()
		queue.Remove(front)
		node := front.Value.(uint)

		for i := uint(0); i < ServerCount; i++ {
			if i != start && (*finalTopo )[node][i] && len(detailedSP[i]) == 0 {
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
