package wave

import (
	"fmt"
	"strconv"
	"strings"
)

type message struct {
	topology [][]bool
	src    uint
	active bool
}

func serialize(m message) string {

	// Topology
	var strTopo string
	for i := range topology {
		if i > 0 {
			strTopo += "_"
		}
		for j := range topology[i] {
			if j > 0 {
				strTopo += "-"
			}
			if topology[i][j] {
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

	return fmt.Sprintf("%v %v %v", strTopo, m.src, strActive)
}

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
	num, _ := strconv.ParseUint(splits[1], 10, 0)
	m.src = uint(num)

	// Active
	if splits[2] == "1" {
		m.active = true
	} else {
		m.active = false
	}

	return m
}