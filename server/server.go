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
	reqpool  map[int]string
	pingopen bool
	reqopen  bool
	handlers map[string]func(w http.ResponseWriter, r *http.Request)
	srv      *http.Server

	// This is the UDP section
	broadcasting bool
	buffer       []byte

	// Hostname + Address
	clienthostname string

	// A connection that the user may have to a node
	connection *conn

	// Alive Channel
	maintainsignal bool
}

// Server Instance creator
var (
	serverinstance *ServerInstance
	once           sync.Once
)

func GetInstance() *ServerInstance {
	return serverinstance
}

func ServerClose() {
	logger := log.GetInstance()

	if serverinstance == nil {
		logger.Output("ERROR", "Server instance doesn't exist nothing to close")
	} else {
		serverinstance.maintainsignal = false
	}
}

// The server runner handling broadcast and normal connections
func ServerRun(alive bool) {

	logger := log.GetInstance()

	if serverinstance != nil {
		logger.Output("SERVER", "There already exists an instance!")
		return
	}

	once.Do(func() {

		// LEARNING: THIS INITIALIZES AND SETS, WE DONT NEED A LOCAL VARIABLE WE JUST NEED TO UPDATE THE GLOBAL VARIABLE
		serverinstance = &ServerInstance{
			pingpool: make(map[string]string),
			reqpool:  make(map[int]string),
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
			broadcasting:   false,
			buffer:         make([]byte, 1024),
			maintainsignal: alive,
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

	logger.Output("SERVER", "Starting server!")

	/* Learning!
	signal 0xc0000005: This is a Windows-specific error indicating an access violation (attempting to access memory that is not valid).
	addr=0x28: This is the memory address that the program tried to access, which is invalid.
	pc=0x7ff7aa2eb870: This is the program counter (instruction pointer) at the time of the crash.
	goroutine 82: The crash occurred in the ServerRun function, which was called in a goroutine.
	*/
	if serverinstance == nil {
		logger.Debug("DEBUG", "Failed to start server! Setting the maintain to be false")
		alive = false
		return
	}

	go func() {
		logger.Debug("DEBUG", fmt.Sprintf("Starting the http, instance %v", serverinstance))
		err := serverinstance.srv.ListenAndServe() // Not http listen and serve, we have our own server
		if err == nil {
			logger.Output("ERROR", "Error in starting the server")
		}
	}()

	logger.Debug("DEBUG", "Server has started!")

	// Server running loop
	go func() {
		// Context
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel() // Learning: Cancels the resources associated with the things we're canceling

		go func() {
			con := createudpcon(8080, "0.0.0.0", true)
			defer con.Close()
			for {
				if con == nil {
					logger.Debug("SERVER | ERROR", "Connection could not be created!")
					break
				}

				n, addr, err := con.ReadFromUDP(serverinstance.buffer)
				if err != nil {
					fmt.Println("ERROR: Could not read from UDP: " + err.Error())
					break
				} else {
					// Check duplicates
					_, exists := serverinstance.pingpool[addr.String()]
					senderhostname, errhostname := net.LookupHost(addr.IP.String())
					if !exists {

						if errhostname != nil {
							// Store [address] = hostname
							serverinstance.pingpool[addr.String()] = strings.Join(senderhostname, " ")
						} else {
							fmt.Printf("Could not resolve hostname!\n")
							serverinstance.pingpool[addr.String()] = ""
						}
					}

					fmt.Printf("Received response from %s: %s\n", addr.String(), string(serverinstance.buffer[:n]))
				}
			}
		}()

		for {

			// BROADCASTING AND LISTENING LOGIC
			if serverinstance.broadcasting {
				con := createudpcon(8080, "255.255.255.255", false)

				if con == nil {
					logger.Debug("SERVER | ERROR", "Connection could not be created!")
				}

				message := []byte("[QFSERVER]ALIVEPING")
				_, err := con.Write(message)

				if err != nil {
					fmt.Printf("Error: %v", err)
				}

				time.Sleep(time.Second * 2)
				fmt.Println("BROADCAST: Sending a broadcast")

				con.Close()
			}

			if !serverinstance.maintainsignal {

				// Shutdown the server
				if err := serverinstance.srv.Shutdown(ctx); err != nil {
					logger.Debug("DEBUG", fmt.Sprintf("Server Shutdown Failed:%+v", err))
				}

				logger.Debug("DEBUG", "Server has been stopped")

				return
			}

			time.Sleep(time.Second * 1)
		}
	}()
}

func createudpcon(port int, ip string, listen bool) *net.UDPConn {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}

	var conCreate *net.UDPConn = nil
	var err error = nil

	if !listen {
		conCreate, err = net.DialUDP("udp", nil, &addr)
	} else {
		conCreate, err = net.ListenUDP("udp", &addr)
	}

	if err != nil {
		fmt.Printf("%v", err)
	} else {
		con := conCreate
		con.SetDeadline(time.Now().Add(5 * time.Second))
		serverinstance.buffer = make([]byte, 1024)
	}

	return conCreate
}
