package main

import (
	"fmt"

	"github.com/QFServer/client"
)

func main() {
	// Begin client loop
	fmt.Printf("\nLOG: Starting client!")
	client.ClientLoop()
}
