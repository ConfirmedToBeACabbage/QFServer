package main

import (
	"github.com/QFServer/client"
	"github.com/QFServer/log"
)

func main() {
	// Begin the logging
	logger := log.GetInstance()
	logger.Store("init", "Log has begun!")

	client.ClientLoop()
}
