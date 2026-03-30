package server

import (
	"fmt"
	"net/http"
	"os"
)

// Functions to pool everything
func (si *ServerInstance) handlereq(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello!")
}

// Ping response and receive
func (si *ServerInstance) handleping(w http.ResponseWriter, r *http.Request) {
	// Store ping in the pool
	address := r.RemoteAddr
	hostname := r.Host

	// Check duplicates
	_, exists := si.pingpool[address]
	if !exists {
		// Store [address] = hostname
		si.pingpool[address] = hostname
	}

	// Write a response
	currhostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "[QFServer]\nHostname: %s\n", currhostname)
	}
}
