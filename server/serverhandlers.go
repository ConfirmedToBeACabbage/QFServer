package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

func (si *ServerInstance) handleconn(w http.ResponseWriter, r *http.Request) {
	// Get the string which correlates to this item you want to handle in this
	specHandle, ok := si.reqpool[strings.Split(r.RemoteAddr, ":")[0]]

	if ok {

		// Get a key I can use for this connection
		token := make([]byte, 4)
		rand.Read(token)

		encryptedMessage, error := rsa.EncryptOAEP(sha256.New(), rand.Reader, &specHandle.masterPublic, specHandle.data, token)
		if error == nil {
			w.Write(encryptedMessage)
		}

	}

}

// Functions to pool everything
func (si *ServerInstance) handlereq(w http.ResponseWriter, r *http.Request) {
	address := strings.Split(r.RemoteAddr, ":")[0]

	// Check duplicates
	_, exists := si.reqpool[address]
	if !exists {
		newConn := new(conn)
		requestBody := r.Body

		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(requestBody)

		// Interpret the request data
		if err == nil {
			content := strings.Split(buf.String(), "||")

			// Get the bigint conversion
			bigIntN := new(big.Int)
			bigIntN, ok := bigIntN.SetString(content[1], 10)

			if !ok { // TODO: CHANGE TO BETTER ERROR
				fmt.Println("ERROR")
			}

			// Get the E conversion
			eInt, err := strconv.Atoi(content[2])

			if err != nil { // TODO: CHANGE TO BETTER ERROR
				fmt.Println("ERROR")
			}

			// Setup the public key object
			pubKey := new(rsa.PublicKey)
			pubKey.N = bigIntN
			pubKey.E = eInt

			// Setup the connection object
			newConn.endpointCON = content[0]
			newConn.sourceCON = address
			newConn.masterPublic = *pubKey

			// Store the request
			si.reqpool[address] = *newConn

			fmt.Println("Secured the connection object")
		}
	}
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
