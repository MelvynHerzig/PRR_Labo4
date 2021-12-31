package common

import "strings"

// SerializeTopology translate an adjacency matrix to a string
// For example, if the matrix is [false, true, false]
//								 [true, false, true]
//								 [false, true, false]
// The string will be: "0-1-0_1-0-1_0-1-0"
func SerializeTopology(matrix *[][]bool) string{
	// Topology
	var strTopo string
	for i := range *matrix {
		if i > 0 {
			strTopo += "_"
		}
		for j := range (*matrix)[i] {
			if j > 0 {
				strTopo += "-"
			}
			if (*matrix)[i][j] {
				strTopo += "1"
			} else {
				strTopo += "0"
			}
		}
	}

	return strTopo
}

// DeserializeTopology translate a string into adjacency matrix.
// For example, if the matrix is [false, true, false]
//								 [true, false, true]
//								 [false, true, false]
// The string must be: "0-1-0_1-0-1_0-1-0"
func DeserializeTopology(str string) [][]bool {
	var topology [][]bool
	topology = make([][]bool, ServerCount)
	for i := range topology {
		topology[i] = make([]bool, ServerCount)
	}

	rows := strings.Split(str, "_")
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

	return topology
}
