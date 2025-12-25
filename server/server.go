package server

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
func (si *ServerInstance) SendBroadcast() {
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

		// Learning: If i'm just looping over one channel I can do this
		for current := range si.broadcast {
			if !current {
				return
			} else {
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
func (si *ServerInstance) ListenBroadcast() {
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

// Changing server states
func (si *ServerInstance) BroadcastStateChange() {
	select {
	case broadcastsignal := <-si.broadcast:
		if !broadcastsignal {
			si.broadcast <- true
		} else {
			si.broadcast <- false
		}
	default:
		return
	}
}

// Open the server to be pinged
func (si *ServerInstance) PingStateChange() bool {
	si.pingopen = !si.pingopen
	return si.pingopen
}

// Open the server to be requested
func (si *ServerInstance) ReqStateChange() bool {
	si.reqopen = !si.reqopen
	return si.reqopen
}

// Server Instance creator
var (
	serverinstance *ServerInstance
	once           sync.Once
)

// Simple check alive for the server instance
func CheckServerAlive() bool {
	if serverinstance == nil {
		return false
	} else {
		return true
	}
}

// The init of the server with once.do for singelton
func ServerInitSingleton() *ServerInstance {
	once.Do(func() {
		instance := &ServerInstance{
			pingpool: make(map[string]string),
			reqpool:  make(map[string]*http.Request),
			pingopen: false,
			reqopen:  false,
			handlers: map[string]func(w http.ResponseWriter, r *http.Request){
				"/":    serverinstance.handleping, // Handle pings
				"/req": serverinstance.handlereq,  // Handle pools
			},
			// Learning: I need to assign the handler here, otherwise we will get a panic when http tries to handle the requests
			srv: &http.Server{
				Addr:    ":8080",              // Set the address and port
				Handler: http.DefaultServeMux, // Use the default ServeMux
			},
			broadcast: make(chan bool, 1),
			buffer:    make([]byte, 1024),
		}

		// Setup
		instance.broadcast <- false

		hostget, errhost := os.Hostname()
		if errhost != nil {
			instance.clienthostname = hostget
		} else {
			fmt.Println("DEBUG: Error in getting hostname!")
		}

		// Setup the handlers
		for i, v := range instance.handlers {
			http.HandleFunc(i, v)
		}

		fmt.Println("DEBUG: Setup the handlers!")

		// Setup server configuration
		instance.srv.IdleTimeout = time.Millisecond * 5
		instance.srv.MaxHeaderBytes = 1024

		// Learnings: The handlers here are specifically talking bout the app.routes() handler. It's sort of middle-ware.
		// The http.HandleFunc simple adds to the routes. The server would speak to that then. That's why it's not in the same struct
		// Listen and server is a blocking thing too. Since it runs in the sync.once it will block the initialization process indefinitely preventing
		// The rest of the program from running. Well technically it's just blocking the goroutine but still the same issue.

		go func() {
			fmt.Println("DEBUG: Starting the http ")
			err := http.ListenAndServe(":8080", nil)
			if err != nil {
				fmt.Println("Error in starting the server", err)
			}
		}()

		fmt.Println("DEBUG: Server has started! You can try letting it broadcast now")
	})

	return serverinstance
}

// The server runner handling broadcast and normal connections
func ServerRun(maintain chan bool) {

	fmt.Println("DEBUG: Starting server!")

	// Get the singleton and use it
	instance := ServerInitSingleton()

	fmt.Println("DEBUG: Server has started!")

	// Context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel() // Learning: Cancels the resources associated with the things we're canceling

	// Server running loop
	for {
		select {
		case broadcastsignal := <-instance.broadcast:
			if broadcastsignal {
				instance.ListenBroadcast()
				instance.SendBroadcast()
			}
		case maintainsignal := <-maintain:
			if !maintainsignal {
				close(maintain)
				close(instance.broadcast)

				// Shutdown the server
				if err := instance.srv.Shutdown(ctx); err != nil {
					log.Fatalf("Server Shutdown Failed:%+v", err)
				}

				return
			}
		default:
			time.Sleep(time.Second * 2)
		}
	}
}
