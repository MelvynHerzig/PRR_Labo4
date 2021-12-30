package probe

import (
	"fmt"
	"server/algorithms"
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

// unknown is the character used to serialize an unknow path to a given node
const unknown = "N"

// message is the type of struct that is sent between servers during the research.
type message struct {
	mType string // Node topology
	id uint      // unique identifier, generally the id of original source emission
	src uint	 // Node Number of source
	temporaryKnownSp [][]uint  // content
}

// serialize translate a message into a string. Type, id and src are simply translated to string.
// temporaryKnownSp sp is translated specially. if the matrix looks like [ [1, 2, 3],
//																		   [2, 3],
//																		   [] ],
// it will be translated 1-2-3_2-3_N
// The resulting string is <type> <id> <src> <str of matrix>
func serialize(m *message) string {

	// temporary known shortest paths
	var strTemporaryKnownSp string
	for i := range m.temporaryKnownSp {
		if i > 0 {
			strTemporaryKnownSp += "_"
		}
		// If no shortest path is known for the node i, put N
		if len(m.temporaryKnownSp[i]) == 0 {
			strTemporaryKnownSp += unknown
		// Else put the known node.
		} else {
			for j := range m.temporaryKnownSp[i] {
				if j > 0 {
					strTemporaryKnownSp += "-"
				}
				strTemporaryKnownSp += fmt.Sprintf("%v", m.temporaryKnownSp[i][j])
			}
		}
	}

	return fmt.Sprintf("%v %v %v %v", m.mType, m.id, m.src, strTemporaryKnownSp)
}

// deserialize translate a string into a message. If the string is not well-formed, the program will crash
// temporaryKnownSp sp is translated specially. if the matrix looks like [ [1, 2, 3],
//																		   [2, 3],
//																		   [] ],
// it will be translated 1-2-3_2-3_N
// The resulting string is <type> <id> <src> <str of matrix>
func deserialize(s string) message {
	var m message

	splits := strings.Split(s, " ")

	m.mType = splits[0]

	id, _ := strconv.ParseUint(splits[1], 10, 0)
	m.id = uint(id)

	src, _ := strconv.ParseUint(splits[2], 10, 0)
	m.src = uint(src)

	// temporary known shortest paths
	var temporaryKnownSp [][]uint
	temporaryKnownSp = make([][]uint, algorithms.ServerCount)

	rows := strings.Split(splits[3], "_")
	for i := range rows {

		values := strings.Split(rows[i], "-")
		if values[0] == unknown {
			temporaryKnownSp[i] = make([]uint, 0,  algorithms.ServerCount)
		} else {
			temporaryKnownSp[i] = make([]uint, len(values), algorithms.ServerCount)
			for j := range values {
				node, _ := strconv.ParseUint(values[j], 10, 0)
				temporaryKnownSp[i][j] = uint(node)
			}
		}
	}
	m.temporaryKnownSp = temporaryKnownSp

	return m
}