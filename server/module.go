package server

import (
	"fmt"
	"time"

	"github.com/QFServer/log"
)

func (si *ServerInstance) REQmodule(alive bool) {

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
	logger.Output("", "Current Pool")
	logger.Output("", "------------")
	for _, v := range pingPool {

		logger.Output("", fmt.Sprintf("%d | %s --- %s%d | %s", counter+1, v, "C", counter+1, reqPool[counter]))

		counter += 1
	}

	for !goodInput {
		input := logger.InputFromUser()

		if input == "quit" {
			logger.SwitchModule("DEFAULT")
			alive = false
			goodInput = true
		}

		time.Sleep(time.Second * 1)
	}
}
