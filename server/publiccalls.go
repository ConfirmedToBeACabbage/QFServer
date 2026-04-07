package server

import (
	"fmt"

	"github.com/QFServer/log"
)

// Request pool
func (si *ServerInstance) GetRequestPool() map[string]conn {
	if !CheckServerAlive() {
		return map[string]conn{"ERROR": {}}
	}
	return si.reqpool
}

// Return all of the ping pools (This is everyone we can contact)
func (si *ServerInstance) GetPingPool() map[string]string {
	if !CheckServerAlive() {
		return map[string]string{"ERROR": "ERROR: Cannot change broadcast since the server isn't alive!"}
	}
	return si.pingpool
}

// Open the server to be pinged
func (si *ServerInstance) PingStateChange() bool {
	logger := log.GetInstance()
	if !CheckServerAlive() {
		logger.Output("ERROR", "Cannot change ping status since the server isn't alive!")
		return false
	}
	si.pingopen = !si.pingopen
	return si.pingopen
}

// Open the server to be requested
func (si *ServerInstance) ReqStateChange() bool {
	logger := log.GetInstance()
	if !CheckServerAlive() {
		logger.Output("ERROR", "Cannot change request status since the server isn't alive!")
		return false
	}
	si.reqopen = !si.reqopen
	return si.reqopen
}

// Changing server states
func BroadcastStateChange() {
	logger := log.GetInstance()
	logger.Output("SERVER", "Attempting to change the broadcast switch")

	if !CheckServerAlive() {
		fmt.Println("ERROR: Cannot change broadcast since the server isn't alive!")
		return
	}

	logger.Output("SERVER", "Instance exists! Server broadcast will be changed now")
	serverinstance.broadcasting = !serverinstance.broadcasting
	logger.Output("SERVER", fmt.Sprintf("Broadcasting: %b", serverinstance.broadcasting))
}

// Simple check alive for the server instance
func CheckServerAlive() bool {
	logger := log.GetInstance()
	logger.Output("SERVER", "Checking the instance!")
	if serverinstance == nil {
		return false
	} else {
		return true
	}
}
