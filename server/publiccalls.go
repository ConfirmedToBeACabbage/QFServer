package server

import "fmt"

// Request pool
func (si *ServerInstance) GetRequestPool() map[int]string {
	if !CheckServerAlive() {
		return map[int]string{0: "ERROR: Cannot change broadcast since the server isn't alive!"}
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
	if !CheckServerAlive() {
		fmt.Println("ERROR: Cannot change ping status since the server isn't alive!")
		return false
	}
	si.pingopen = !si.pingopen
	return si.pingopen
}

// Open the server to be requested
func (si *ServerInstance) ReqStateChange() bool {
	if !CheckServerAlive() {
		fmt.Println("ERROR: Cannot change request status since the server isn't alive!")
		return false
	}
	si.reqopen = !si.reqopen
	return si.reqopen
}

// Changing server states
func BroadcastStateChange() {
	fmt.Println("SERVER: Attempting to change the broadcast switch")

	if !CheckServerAlive() {
		fmt.Println("ERROR: Cannot change broadcast since the server isn't alive!")
		return
	}

	fmt.Println("SERVER: Instance exists!")
	serverinstance.broadcasting = !serverinstance.broadcasting
}

// Simple check alive for the server instance
func CheckServerAlive() bool {
	fmt.Println("SERVER: Checking the instance!")
	if serverinstance == nil {
		return false
	} else {
		return true
	}
}
