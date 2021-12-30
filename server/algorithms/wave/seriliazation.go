package wave

import (
	"fmt"
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
// For example, if the server src is 2, it is active, it is at the 4th round and the matrix is [0, 1, 0]
//																							   [1, 0, 1]
//																							   [0, 1, 0]
// The string will be: "0-1-0_1-0-1_0-1-0 2 4 1" (i.e <matrix> <src num> <wave num> <active true= 1 false = 0>)
func serialize(m *message) string {

	// Topology
	var strTopo string
	for i := range m.topology {
		if i > 0 {
			strTopo += "_"
		}
		for j := range m.topology[i] {
			if j > 0 {
				strTopo += "-"
			}
			if m.topology[i][j] {
				strTopo += "1"
			} else {
				strTopo += "0"
			}
		}
	}

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
// For example, if the server src is 2, it is active, it is at the 4th round and the matrix is [0, 1, 0]
//																							   [1, 0, 1]
//																							   [0, 1, 0]
// The string must be: "0-1-0_1-0-1_0-1-0 2 4 1" (i.e <matrix> <src num> <wave num> <active true= 1 false = 0>)
func deserialize(s string) message {
	var m message
	var topology [][]bool
	topology = make([][]bool, serverCount)
	for i := range topology {
		topology[i] = make([]bool, serverCount)
	}

	splits := strings.Split(s, " ")

	// Topology
	rows := strings.Split(splits[0], "_")
	for i := range rows {
		values := strings.Split(rows[i], "-")
		for j := range values {
			if values[j] == "1" {
				topology[i][j] = true
			} else {
				topology[i][j] = false
			}
		}
	}
	m.topology = topology

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