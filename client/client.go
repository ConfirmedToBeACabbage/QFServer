package Client

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

	go func() {

		for {
			fmt.Println(`QFServer CLI! Type in "Help" to get started.`)
			fmt.Print("> ")
			input, err := reader.ReadString('\n')

			if err != nil {
				fmt.Println(input)
			} else {
				exit <- true
			}
		}

	}()

	select {
	case <-exit:
		return
	}
}
