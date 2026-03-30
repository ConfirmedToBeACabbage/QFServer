package server

// Connection structure
type conn struct {
	// Basic check
	sourceCON   bool
	endpointCON bool

	// Basic
	sourceIP   string
	endpointIP string

	// Asymettric keys
	privateKey string // CHange types as we go
	publicKey  string

	// Symmetric
	encryptKey string

	// Connection flags
	flagPUBKEY         bool // Public key is exchanged
	flagPRIVKEY        bool // We have a private key
	flagFORWARDSECRECY bool // First part of handshake done

	flagENCRYPTKEY       bool // Symmetric key is established
	flagTRANSFERCOMPLETE bool // Transfer of information is complete
}

// Send information commands

// Handshake commands
