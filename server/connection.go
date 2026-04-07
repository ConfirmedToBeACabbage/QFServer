package server

import "crypto/rsa"

// Connection structure
type conn struct {
	// Basic check
	sourceCON   string
	endpointCON string

	// Information
	data []byte

	// Keys
	masterPublic rsa.PublicKey
}

// Send information commands

// Handshake commands
