package client

// The client:
// 1. Manage the broker for workers which can be created from commands
// 2. Contain the command management
// 3. Center of control

import (
	"bufio"
	"fmt"
	"os"
)

func ClientLoop() {

	// Reader
	reader := bufio.NewReader(os.Stdin)

	// Quit channel
	exitclient := make(chan bool)

	// Give us the broker
	br := InitBroker()

	if br == nil {
		fmt.Println("\nERROR: Broker has not begun")
		exitclient <- true
	}

	select {
	case <-exitclient:
		return
	default:
		fmt.Println("\nLOG: Default")
		go func() {

			fmt.Println("\nLOG: Inside")
			for {
				fmt.Println("\nQFServer CLI! Type in ->Help<- to get started.")
				fmt.Print("> ")
				input, err := reader.ReadString('\n')

				// Parsing
				inputparse := Parse(input)
				fmt.Printf("\nLOG: We have parsed!")

				// Redirect with the command
				if !inputparse.giveerror {
					ok := br.configureworker(inputparse)

					fmt.Printf("\nLOG: [Broker] Return value %v", ok)

					if !ok {
						fmt.Printf("\nLOG: We have added a worker for this command!")
						continue
					} else {
						fmt.Printf("\nLOG: We have not made the worker!")
						fmt.Printf("\nLOG: [Broker] %s", br.message)
					}
				} else {
					fmt.Println(inputparse.message)
					exitclient <- true
				}

				if err != nil {
					fmt.Println(input)
				} else {
					exitclient <- true
				}
			}

		}()
		fmt.Println("\nLOG: Default exiting")

	}

	<-exitclient
}
