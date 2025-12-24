package server

// TODO:
// This needs to all be in a structure and as a singleton
// Also have shutdown methods etc
// The shutdown would be called when the maintain channel closes
// Startup method starts from the start method

// The listener can act as both the receiver and listener. Should probablyy change the name

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// The server struct
// Pool of incoming requests
// Pool of lan avaliable users
type ServerInstance struct {
	// TLS section
	pingpool map[string]string
	reqpool  map[string]*http.Request
	pingopen bool
	reqopen  bool
	handlers map[string]func(w http.ResponseWriter, r *http.Request)
	srv      *http.Server

	// This is the UDP section
	broadcast chan bool
	buffer    []byte

	// Hostname + Address
	clienthostname string
}

// Functions to pool everything
func (si *ServerInstance) handlereq(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello!")
}

// Ping response and receive
func (si *ServerInstance) handleping(w http.ResponseWriter, r *http.Request) {
	// Store ping in the pool
	// - Check duplicates
	// Respond with a yes or no to the ping
	address := r.RemoteAddr
	hostname := r.Host

	// Check duplicates
	_, exists := si.pingpool[address]
	if !exists {
		// Store [address] = hostname
		si.pingpool[address] = hostname
	}

	// Write a response
	currhostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "[QFServer]\nHostname: %s\n", currhostname)
	}
}

// Return all of the ping pools (This is everyone we can contact)
func (si *ServerInstance) GetPingPool() map[string]string {
	return si.pingpool
}

// Send an alive message
func (si *ServerInstance) SendAlive() {
	go func() {
		addr := net.UDPAddr{
			Port: 12345,
			IP:   net.ParseIP("255.255.255.255"),
		}

		con, err := net.DialUDP("udp", nil, &addr)
		if err != nil {
			fmt.Printf("%v", err)
		}
		defer con.Close()

		select {
		case <-si.broadcast:
			return
		default:
			for {
				// Write our message
				message := []byte("[QFSERVER]ALIVEPING")
				_, err = con.Write(message)

				if err == nil {
					fmt.Printf("Error: %v", err)
				}

				time.Sleep(time.Second * 5)
				fmt.Println("Sending a broadcast")
			}
		}
	}()
}

// Listen for alive
func (si *ServerInstance) ListenAlive() {
	go func() {
		addr := net.UDPAddr{
			Port: 12345,
			IP:   net.ParseIP("0.0.0.0"),
		}

		con, err := net.DialUDP("udp", nil, &addr)
		if err != nil {
			fmt.Printf("%v", err)
		}
		defer con.Close()

		con.SetDeadline(time.Now().Add(5 * time.Second))

		si.buffer = make([]byte, 1024)

		for {
			n, addr, err := con.ReadFromUDP(si.buffer)
			if err != nil {
				break
			}

			// Check duplicates
			_, exists := si.pingpool[addr.String()]
			senderhostname, errhostname := net.LookupHost(addr.IP.String())
			if !exists {

				if errhostname != nil {
					// Store [address] = hostname
					si.pingpool[addr.String()] = strings.Join(senderhostname, " ")
				} else {
					fmt.Printf("Could not resolve hostname!")
				}
			}

			fmt.Printf("Received response from %s: %s\n", addr.String(), string(si.buffer[:n]))
		}

		<-si.broadcast
	}()
}

// Open channel for listening
func (si *ServerInstance) AliveChange() {
	si.broadcast <- !<-si.broadcast

	select {
	case current := <-si.broadcast:
		if !current {
			si.broadcast <- true
			si.ListenAlive()
			si.SendAlive()
		} else {
			si.broadcast <- false
			si.buffer = make([]byte, 0)
		}
	default:
		return
	}
}

// Open port command
var (
	instance *ServerInstance
	once     sync.Once
)

// Server Instance
func ServerInitSingleton() *ServerInstance {
	once.Do(func() {
		instance := &ServerInstance{
			pingpool: make(map[string]string),
			reqpool:  make(map[string]*http.Request),
			pingopen: false,
			reqopen:  false,
			handlers: map[string]func(w http.ResponseWriter, r *http.Request){
				"/":    instance.handleping, // Handle pings
				"/req": instance.handlereq,  // Handle pools
			},
			srv:       &http.Server{},
			broadcast: make(chan bool),
			buffer:    make([]byte, 1024),
		}

		// Setup
		instance.broadcast <- false

		hostget, errhost := os.Hostname()
		if errhost != nil {
			instance.clienthostname = hostget
		}

		// Setup the handlers
		for i, v := range instance.handlers {
			http.HandleFunc(i, v)
		}

		// Setup server configuration
		instance.srv.IdleTimeout = time.Millisecond * 5
		instance.srv.MaxHeaderBytes = 1024

		// Learnings: The handlers here are specifically talking bout the app.routes() handler. It's sort of middle-ware.
		// The http.HandleFunc simple adds to the routes. The server would speak to that then. That's why it's not in the same struct
		err := http.ListenAndServe(":8080", nil)

		if err != nil {
			fmt.Println("Error in starting the server", err)
		}
	})

	return instance
}

// Singleton of the server
func ServerRun(exit chan bool) {

	// Get the singleton and use it
	instance := ServerInitSingleton()

	// TODO: select statemetn for both udp and tls (this would be the exit channel and the broadcast channel)
	// 	if <-si.broadcast { // If its true
	// Lets listen here
	// Go routine here which just listens
	//}

	// Signal we wait on to exit
	<-exit
	close(exit)

	// Learning: context is just created so we can shutdown the server in this case, in 5 seconds in the background
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel() // Learning: Cancels the resources associated with the things we're canceling

	// Shutdown the server
	if err := instance.srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
}
