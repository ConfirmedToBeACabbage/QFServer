package client

// The client:
// 1. Manage the broker for workers which can be created from commands
// 2. Contain the command management
// 3. Center of control

import (
	"strings"
	"time"

	"github.com/QFServer/log"
)

func ClientLoop() {

	// Logger (We're passing the input channel here)
	logger := log.GetInstance()
	logger.BeginDebugLogger()

	// Quit channel
	exitclient := make(chan bool, 1)

	// Give us the broker
	br := InitBroker()

	if br == nil {
		logger.Debug("CLIENT", "Broker has not begun")
		exitclient <- true
	}

	go func() {

		for {
			select {
			default:
				if logger.CheckModule() == "DEFAULT" {
					time.Sleep(time.Second * 1)
					input := logger.InputFromUser()

					// Default exit TODO: (Should be moved to a command)
					if strings.TrimSpace(input) == "quit" {
						shutdownbroker := make(chan bool)
						go br.gracefulshutdown(shutdownbroker)
						<-shutdownbroker
						exitclient <- true
					}

					// Parsing
					inputparse := Parse(input)
					logger.Debug("CLIENT", "We have completed command parsing!")

					// Redirect with the command
					if !inputparse.giveerror {
						logger.Debug("DEBUG", "Lets begin configuration!")
						errorreceive := br.configureworker(inputparse)
						logger.Debug("DEBUG", "We have configured!")

						if !errorreceive {
							logger.Debug("CLIENT", "The worker has not been made by the broker"+br.message)
						} else {
							logger.Debug("DEBUG", "Error in creating the worker!")
						}

					}

				}
			case exit := <-exitclient:
				if exit {
					return
				}
			}
		}
	}()

	<-exitclient
	close(exitclient)
}
