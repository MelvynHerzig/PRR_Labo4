// Package algorithms define a common init methode shared between probe and wave implementation
package algorithms

import "prr.configuration/config"

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

	IsMyNeighbors   = make([]bool, ServerCount)
	NeighborsCount = uint(len(config.GetNeighbors(LocalNumber)))
	for _, v := range config.GetNeighbors(LocalNumber) {
		IsMyNeighbors[v] = true
	}
}
