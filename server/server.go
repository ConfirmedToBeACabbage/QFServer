package Server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func createserver(typeS string, doneS chan string) {

	// Setting up pointer
	server := &http.Server{
		Addr:           ":8090",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Handler
	// http.HandleFunc("/ping", ping)
	// http.HandleFunc("/session", session)

	if typeS == "receiver" {
		// Provide handlers for receiver
	} else {
		// Provide handlers for sender
	}

	// Error handling ListenAndServe
	for {
		select {
		case <-doneS:
			fmt.Println("SERVER: EXITING")
			return
		default:
			if err := server.ListenAndServe(); err != nil {
				log.Fatal(err)
			}
		}
	}
}

// Init the server (Type)
// Options:
// 1. Sender
// - Broadcast listen on LAN to find other receivers
// - Gather and collect receivers
// 2. Receiver
// - Listen in for sender requests
func InitServer(typeS string, workcloseCH chan string) {

	// Create the server
	go createserver(typeS, workcloseCH)

	// // Setting up signal for interrupt checking in a goroutine
	// exitChan := make(chan struct{})
	// go func() {
	// 	sigs := make(chan os.Signal, 1)
	// 	signal.Notify(sigs, os.Interrupt)
	// 	<-sigs

	// 	fmt.Println("GATEWAY: Received signal to close test gateway")
	// 	close(exitChan)
	// }()

	// fmt.Println("GATEWAY: Gateway is setup and awaiting exit signal")

	// // Awaiting the channel to exit
	// <-exitChan

	// // Context for exiting (This context visually is like surrounding the function it's a part of. It provides context on what to do. Defer lets you finish the context workflow.)
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := server.Shutdown(ctx); err != nil {
	// 	log.Fatal(err)
	// }

}
