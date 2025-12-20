package client

// The client:
// 1. Manage the broker for workers which can be created from commands
// 2. Contain the command management
// 3. Center of control

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/QFServer/log"
)

func ClientLoop() {

	// Logger
	logger := log.GetInstance()

	// Reader
	reader := bufio.NewReader(os.Stdin)

	// Quit channel
	exitclient := make(chan bool, 1)

	// Ready for input channel
	readyforinput := make(chan bool, 1)
	readyforinput <- true

	// Give us the broker
	br := InitBroker(readyforinput)

	if br == nil {
		logger.Store("CLIENT", "Broker has not begun")
		exitclient <- true
	}

	go func() {

		for {
			select {
			case ready := <-readyforinput:
				if ready {
					fmt.Println("\nQFServer CLI! Type in - Help - to get started.")
					fmt.Print("> ")
					input, err := reader.ReadString('\n')

					// Default exit
					if strings.TrimSpace(input) == "quit" {
						exitclient <- true
					}

					// Parsing
					inputparse := Parse(input)
					logger.Store("CLIENT", "We have completed command parsing!")

					// Redirect with the command
					if !inputparse.giveerror {
						ok := br.configureworker(inputparse)

						if ok {
							logger.Store("CLIENT", "The worker has not been made by the broker"+br.message)
						}

					} else {
						fmt.Println(inputparse.message)
					}

					if err != nil {
						fmt.Println(input)
					}
				} else {
					time.Sleep(time.Millisecond * 100)
				}
			case exit := <-exitclient:
				if exit {
					return
				}
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	<-exitclient
}
