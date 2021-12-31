package wave

import (
	"fmt"
	"server/algorithms/common"
	"strconv"
	"strings"
)

// message is the type of struct that is sent between servers during the research.
type message struct {
	topology [][]bool // Node topology
	src uint		  // Src node
	wave uint		  // Number of the wave
	active bool       // If the src node is active
}

// serialize translate a message into a string
func serialize(m *message) string {

	// Topology
	strTopo := common.SerializeTopology(&m.topology)

	// Active
	var strActive string
	if m.active {
		strActive = "1"
	} else {
		strActive = "0"
	}

	return fmt.Sprintf("%v %v %v %v", strTopo, m.src, m.wave, strActive)
}

// deserialize translate a string into a message. If the string is not well-formed, the program will crash
func deserialize(s string) message {
	var m message
	var topology [][]bool
	topology = make([][]bool, common.ServerCount)
	for i := range topology {
		topology[i] = make([]bool, common.ServerCount)
	}

	splits := strings.Split(s, " ")

	// Topology
	m.topology = common.DeserializeTopology(splits[0])

	// Number
	src, _ := strconv.ParseUint(splits[1], 10, 0)
	m.src = uint(src)

	wave, _ := strconv.ParseUint(splits[2], 10, 0)
	m.wave = uint(wave)

	// Active
	if splits[3] == "1" {
		m.active = true
	} else {
		m.active = false
	}

	return m
}