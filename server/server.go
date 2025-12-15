package Server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func InitGateway() {

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

	// Error handling ListenAndServe
	go func() {

		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}

	}()

	fmt.Println("GATEWAY: Up and running!")

	// Setting up signal for interrupt checking in a goroutine
	exitChan := make(chan struct{})
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt)
		<-sigs

		fmt.Println("GATEWAY: Received signal to close test gateway")
		close(exitChan)
	}()

	fmt.Println("GATEWAY: Gateway is setup and awaiting exit signal")

	// Awaiting the channel to exit
	<-exitChan

	// Context for exiting (This context visually is like surrounding the function it's a part of. It provides context on what to do. Defer lets you finish the context workflow.)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

}
