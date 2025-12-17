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
	exit := make(chan bool)

	// Give us the broker
	br := InitBroker()

	if br == nil {
		fmt.Println("ERROR: Broker has not begun")
		<-exit
	}

	go func() {

		for {
			fmt.Println(`QFServer CLI! Type in "Help" to get started.`)
			fmt.Print("> ")
			input, err := reader.ReadString('\n')

			inputparse := Parse(input)

			// Redirect with the command
			if !inputparse.giveerror {
				br.addworker(inputparse)
			} else {
				fmt.Println(inputparse.message)
				exit <- true
			}

			if err != nil {
				fmt.Println(input)
			} else {
				exit <- true
			}
		}

	}()

	if <-exit {
		return
	}
}
