// Package config defines methods to get setup configurations
package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// reader that parses configuration
var reader *configReader = nil

// localServerNumber stores the server number, unused if config is initialized by a client
var localServerNumber uint

const (
	AlgoWave  = "wave"
	AlgoProbe = "probe"
)

// Server is a struct that defines the characteristics of a server
type Server struct {
	Ip        string `json:"ip"`
	Port      uint   `json:"port"`
	Neighbors []uint `json:"neighbors"`
}

// configReader defines the field to read in the config file
type configReader struct {
	Debug   bool     `json:"debug"`
	Version string   `json:"version"`
	Servers []Server `json:"servers"`
}

// Init Loads the config file as a list of Server
// path is used to find the config file
// serverNumber is used on the server side, to remember who we are
func Init(path string, serverNumber uint) {
	rand.Seed(time.Now().UnixNano())
	jsonFile, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &reader)

	localServerNumber = serverNumber
}

// InitSimple Init the reader. Used by client because the second parameter is useless
func InitSimple(path string) {
	Init(path, 0)
}

// IsDebug Returns the value of debug in config file
func IsDebug() bool {
	if reader == nil {
		log.Fatal("config not initialized")
	}

	return reader.Debug
}

// GetServerById Returns the server corresponding to the specified id
func GetServerById(id uint) *Server {
	if reader == nil {
		log.Fatal("config not initialized")
	}

	if id >= uint(len(reader.Servers)) {
		return nil
	}

	return &reader.Servers[id]
}

// GetServerRandomly Returns a server from the servers list
func GetServerRandomly() *Server {
	if reader == nil {
		log.Fatal("config not initialized")
	}

	return GetServerById(uint(rand.Intn(len(reader.Servers))))
}

// GetServers Returns all servers from the list
func GetServers() []Server {
	if reader == nil {
		log.Fatal("config not initialized")
	}

	return reader.Servers
}

// IsServerIP Checks if the ip sent correspond to one of the server in the config file
func IsServerIP(address string) bool {
	if reader == nil {
		log.Fatal("config not initialized")
	}

	var ip = strings.Split(address, ":")[0]
	for _, server := range reader.Servers {
		if server.Ip == ip || (isLocalhost(server.Ip) && isLocalhost(ip)) {
			return true
		}
	}

	return false
}

// GetNeighbors Gets all neighbors id of server id sent
func GetNeighbors(id uint) []uint {
	if reader == nil {
		log.Fatal("config not initialized")
	}

	return reader.Servers[id].Neighbors
}

// isLocalhost Checks if an address is localhost
func isLocalhost(address string) bool {
	return address == "127.0.0.1" || address == "localhost"
}

// GetLocalServerNumber Gets the id specified in the Init function
func GetLocalServerNumber() uint {
	return localServerNumber
}

// GetAlgorithm Gets the algorithm to use
func GetAlgorithm() string {
	if reader == nil {
		log.Fatal("config not initialized")
	}

	algo := reader.Version

	if algo != AlgoWave && algo != AlgoProbe {
		log.Fatal("Unknown algorithm set")
	}

	return algo
}