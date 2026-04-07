package server

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	FR "github.com/QFServer/fr"
	"github.com/QFServer/log"
)

func (si *ServerInstance) REQmodule(alive chan bool) {

	// Show the list of addresses in the request pool indexed with number starting at 1
	// Beside it show the list of addresses which are requesting to you with different indexes
	// AKA 1 for address 1 in the request
	// AKA C1 to address the connection attempt to you from 1

	// Either we make a request or we accept a request
	// If we make a request then we add the request as a connection object to the pool of connections
	// If we make a connection acceptance then we have a job which takes the connection object that is coming from that and works on it

	logger := log.GetInstance()
	goodInput := false

	pingPool := si.GetPingPool()
	reqPool := si.GetRequestPool()

	// Switch modules
	logger.SwitchModule("SERVERREQ")

	counter := 0
	logger.Output("SERVERREQ", "Current Pool")
	pingablePool := make(map[int]string)
	requestablePool := make(map[int]string)
	for i := range pingPool {

		logger.Output("NODE", fmt.Sprintf("%d | %s", counter+1, i))
		pingablePool[counter] = i

		counter += 1
	}

	counter = 0
	for i := range reqPool {

		logger.Output("REQ", fmt.Sprintf("C%d | %s", counter+1, i))
		requestablePool[counter] = i

		counter += 1
	}

	for !goodInput {
		input := logger.InputFromUser()

		if input == "quit" {
			logger.SwitchModule("DEFAULT")
			alive <- false
			goodInput = true
		}

		if input == "1" { // TODO: THis is hardcoded for now
			nodeToPing, exist := pingablePool[0]

			if exist == false {
				logger.Debug("ERROR", "That entry doesnt exist!")
			} else {
				// Prepare the file that we want to send over
				getFile := FR.ReadFromFile(filepath.Join(os.TempDir(), "example"))

				// Generate a private key
				masterPriv, _ := rsa.GenerateKey(rand.Reader, 2048)

				// We store this information inside a connection for this node that we're requesting to
				connObject := &conn{
					endpointCON:  nodeToPing,
					data:         getFile,
					masterPublic: masterPriv.PublicKey,
				}

				si.connection[nodeToPing] = connObject

				// TODO: This should be in a key manager internally
				// Storing the private key (Since we're the requesting node we're the "master" or "server")
				si.conKeyPriv[si.connection[nodeToPing]] = masterPriv

				// Building the reader for the connection | We're sending only vital information to establish a secure connection
				r := strings.NewReader(
					fmt.Sprintf("%s||%s||%s",
						connObject.endpointCON,
						connObject.masterPublic.N.String(),
						strconv.Itoa(connObject.masterPublic.E)))

				// Send over the connection object
				http.Post("http://"+nodeToPing+":8080"+"/req", "text/plain", r)
			}
		}

		// if input == "C1" {
		// 	nodeToPing, exist := requestablePool[0]

		// 	if exist == false {
		// 		logger.Debug("ERROR", "That entry doesn't exist!")
		// 	} else {
		// 		// This means that we
		// 	}
		// }

		time.Sleep(time.Second * 1)
	}
}
