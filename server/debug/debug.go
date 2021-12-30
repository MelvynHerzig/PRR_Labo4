// Package debug offers function to print message only if debug mode is set in the config.json file
package debug

import (
	"fmt"
	"prr.configuration/config"
	"time"
)

// LogReceive logs a message that is received on server
func LogReceive(message string) {
	debugLog("RECEIVED) " + message)
}

// LogSend logs a message that is sent from server
func LogSend(message string) {
	debugLog("SENDED) " + message)
}

// LogWave log the current wave being processed
func LogWave(wave uint) {
	debugLog(fmt.Sprintf(" ------------------ WAVE %v ------------------ ", wave))
}

// debugLog logs the message with DEBUG >> prefix.
func debugLog(message string) {
	if config.IsDebug() {
		fmt.Println("DEBUG >>", time.Now().Format(time.Stamp), message)
	}
}