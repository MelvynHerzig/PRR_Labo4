package probe

import (
	"fmt"
	"server/algorithms/common"
	"strconv"
	"strings"
)

// Types of message that can be sent or received
const (
	// TypeProbe type used for probe
	TypeProbe = "probe"

	// TypeEcho type used for echo
	TypeEcho = "echo"
)

// message is the type of struct that is sent between servers during the research.
type message struct {
	mType string // Node topology
	id uint      // unique identifier, generally the id of original source emission
	src uint	 // Node Number of source
	topology [][]bool  // content
}

// serialize translate a message into a string. Type, id and src are simply translated to string.
func serialize(m *message) string {

	// topology
	strTemporaryKnownSp := common.SerializeTopology(&m.topology)

	return fmt.Sprintf("%v %v %v %v", m.mType, m.id, m.src, strTemporaryKnownSp)
}

// deserialize translate a string into a message. If the string is not well-formed, the program will crash
func deserialize(s string) message {
	var m message

	splits := strings.Split(s, " ")

	m.mType = splits[0]

	id, _ := strconv.ParseUint(splits[1], 10, 0)
	m.id = uint(id)

	src, _ := strconv.ParseUint(splits[2], 10, 0)
	m.src = uint(src)

	// temporary known shortest paths
	m.topology = common.DeserializeTopology(splits[3])

	return m
}