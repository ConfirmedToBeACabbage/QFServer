package server

// TODO:
// This needs to all be in a structure and as a singleton
// Also have shutdown methods etc
// The shutdown would be called when the maintain channel closes
// Startup method starts from the start method

// The listener can act as both the receiver and listener. Should probablyy change the name

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello!")
}

func InitServer() {
	http.HandleFunc("/", handler)

	fmt.Println(("Listening on port 8080..."))
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error in starting the server", err)
	}
}
