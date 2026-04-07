package server

import (
	"fmt"
	"net/http"
	"strings"
)

// Functions to pool everything
func (si *ServerInstance) handlereq(w http.ResponseWriter, r *http.Request) {
	address := strings.Split(r.RemoteAddr, ":")[0]
	hostname := r.Host

	// Check duplicates
	_, exists := si.reqpool[address]
	if !exists {
		// Store [address] = hostname
		si.reqpool[address] = hostname
	}

	fmt.Println("GOT IT!")
}

// Ping response and receive
func (si *ServerInstance) handleping(w http.ResponseWriter, r *http.Request) {
	// Store ping in the pool
	address := strings.Split(r.RemoteAddr, ":")[0]
	hostname := r.Host

	// Check duplicates
	_, exists := si.pingpool[address]
	if !exists {
		// Store [address] = hostname
		si.pingpool[address] = hostname
	}
}
