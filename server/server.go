package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/QFServer/log"
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
	broadcasting        bool
	alreadybroadcasting bool
	buffer              []byte

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
	if !CheckServerAlive() {
		return map[string]string{"ERROR": "ERROR: Cannot change broadcast since the server isn't alive!"}
	}
	return si.pingpool
}

// Send an alive message
func (si *ServerInstance) sendbroadcast() {
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
		for si.broadcasting {
			message := []byte("[QFSERVER]ALIVEPING")
			_, err = con.Write(message)

			if err != nil {
				fmt.Printf("Error: %v", err)
			}

			time.Sleep(time.Second * 2)
			fmt.Println("BROADCAST: Sending a broadcast")
		}
	}()
}

// Listen for alive
func (si *ServerInstance) listenbroadcast() {
	go func() {
		addr := net.UDPAddr{
			Port: 12345,
			IP:   net.ParseIP("0.0.0.0"),
		}

		con, err := net.ListenUDP("udp", &addr)
		if err != nil {
			fmt.Printf("%v", err)
		}
		defer con.Close()

		con.SetDeadline(time.Now().Add(5 * time.Second))

		si.buffer = make([]byte, 1024)

		for si.broadcasting {

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
					fmt.Printf("Could not resolve hostname!\n")
					si.pingpool[addr.String()] = ""
				}
			}

			fmt.Printf("Received response from %s: %s\n", addr.String(), string(si.buffer[:n]))
		}
	}()
}

// Changing server states
func BroadcastStateChange() {
	fmt.Println("SERVER: Attempting to change the broadcast switch")

	if !CheckServerAlive() {
		fmt.Println("ERROR: Cannot change broadcast since the server isn't alive!")
		return
	}

	fmt.Println("SERVER: Instance exists!")
	serverinstance.broadcasting = !serverinstance.broadcasting

	if !serverinstance.broadcasting {
		serverinstance.alreadybroadcasting = false
	}
}

// Open the server to be pinged
func (si *ServerInstance) PingStateChange() bool {
	if !CheckServerAlive() {
		fmt.Println("ERROR: Cannot change ping status since the server isn't alive!")
		return false
	}
	si.pingopen = !si.pingopen
	return si.pingopen
}

// Open the server to be requested
func (si *ServerInstance) ReqStateChange() bool {
	if !CheckServerAlive() {
		fmt.Println("ERROR: Cannot change request status since the server isn't alive!")
		return false
	}
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
	fmt.Println("SERVER: Checking the instance!")
	if serverinstance == nil {
		return false
	} else {
		return true
	}
}

// The init of the server with once.do for singelton
func ServerInitSingleton() *ServerInstance {

	logger := log.GetInstance()

	once.Do(func() {

		// LEARNING: THIS INITIALIZES AND SETS, WE DONT NEED A LOCAL VARIABLE WE JUST NEED TO UPDATE THE GLOBAL VARIABLE
		serverinstance = &ServerInstance{
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
			broadcasting: false,
			buffer:       make([]byte, 1024),
		}

		hostget, errhost := os.Hostname()
		if errhost != nil {
			serverinstance.clienthostname = hostget
		} else {
			logger.Debug("DEBUG", "Error in getting hostname!")
		}

		// Setup the handlers
		// Learnings: The handlers here are specifically talking bout the app.routes() handler. It's sort of middle-ware.
		// The http.HandleFunc simple adds to the routes. The server would speak to that then. That's why it's not in the same struct
		// Listen and server is a blocking thing too. Since it runs in the sync.once it will block the initialization process indefinitely preventing
		// The rest of the program from running. Well technically it's just blocking the goroutine but still the same issue.
		for i, v := range serverinstance.handlers {
			http.HandleFunc(i, v)
		}

		logger.Debug("DEBUG", "Setup the handlers!")

		// Setup server configuration
		serverinstance.srv.IdleTimeout = time.Millisecond * 5
		serverinstance.srv.MaxHeaderBytes = 1024
	})

	return serverinstance
}

// The server runner handling broadcast and normal connections
func ServerRun(alive chan bool) {

	fmt.Printf("DEBUG: Starting server!")

	// Get the singleton and use it
	ServerInitSingleton()

	instance := serverinstance

	logger := log.GetInstance()

	/* Learning!
	signal 0xc0000005: This is a Windows-specific error indicating an access violation (attempting to access memory that is not valid).
	addr=0x28: This is the memory address that the program tried to access, which is invalid.
	pc=0x7ff7aa2eb870: This is the program counter (instruction pointer) at the time of the crash.
	goroutine 82: The crash occurred in the ServerRun function, which was called in a goroutine.
	*/
	if instance == nil {
		logger.Debug("DEBUG", "Failed to start server! Setting the maintain to be false")
		alive <- false
		return
	}

	go func() {
		logger.Debug("DEBUG", fmt.Sprintf("Starting the http, instance %v", instance))
		err := instance.srv.ListenAndServe() // Not http listen and serve, we have our own server
		if err == nil {
			fmt.Println("Error in starting the server", err)
		}
	}()

	logger.Debug("DEBUG", "Server has started!")

	// Server running loop
	go func() {
		// Context
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel() // Learning: Cancels the resources associated with the things we're canceling

		for {

			if instance.broadcasting && !instance.alreadybroadcasting {
				instance.listenbroadcast()
				instance.sendbroadcast()
				instance.alreadybroadcasting = true
			}

			select {

			case maintainsignal := <-alive:
				if !maintainsignal {

					// Shutdown the server
					if err := instance.srv.Shutdown(ctx); err != nil {
						logger.Debug("DEBUG", fmt.Sprintf("Server Shutdown Failed:%+v", err))
					}

					logger.Debug("DEBUG", "Server has been stopped")

					return
				}
			default:
				time.Sleep(time.Second * 2)
			}
		}
	}()
}
