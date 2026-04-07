package server

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/QFServer/log"
)

// The server struct
type ServerInstance struct {
	// TLS section
	pingpool map[string]string
	reqpool  map[string]conn
	pingopen bool
	reqopen  bool
	srv      *http.Server

	handlerInterface *http.ServeMux

	// This is the UDP section
	broadcasting bool
	buffer       []byte

	// Hostname + Address
	clienthostname string

	// A connection that the user may have to a node
	connection map[string]*conn
	conKeyPriv map[*conn]*rsa.PrivateKey

	// Alive Channel
	maintainsignal chan bool
}

// Server Instance creator
var (
	serverinstance *ServerInstance
)

func GetInstance() *ServerInstance {
	return serverinstance
}

func ServerClose() {
	logger := log.GetInstance()

	if serverinstance == nil {
		logger.Output("ERROR", "Server instance doesn't exist nothing to close")
	} else {
		serverinstance.maintainsignal <- false
	}
}

// The server runner handling broadcast and normal connections
func ServerRun(alive chan bool) {

	logger := log.GetInstance()

	// LEARNING: THIS INITIALIZES AND SETS, WE DONT NEED A LOCAL VARIABLE WE JUST NEED TO UPDATE THE GLOBAL VARIABLE
	tempHandle := http.NewServeMux()
	serverinstance = &ServerInstance{
		pingpool: make(map[string]string),
		reqpool:  make(map[string]conn),
		pingopen: false,
		reqopen:  false,
		// Learning: I need to assign the handler here, otherwise we will get a panic when http tries to handle the requests
		handlerInterface: tempHandle,
		srv: &http.Server{
			Addr: ":8080", // Set the address and port
			// ANOTHER LEARNING: Handler is just the interface, but ServeMux actually implements it. That's why you should assign newservemux seperately and then get
			// handlers on it.
			Handler: tempHandle, // Use a new ServeMux LEARNING, this is important to the shutdowns and everything
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
	serverinstance.handlerInterface.HandleFunc("/", serverinstance.handleping)
	serverinstance.handlerInterface.HandleFunc("/req", serverinstance.handlereq)

	logger.Debug("DEBUG", "Setup the handlers!")

	// Setup server configuration
	serverinstance.srv.IdleTimeout = time.Millisecond * 5
	serverinstance.srv.MaxHeaderBytes = 1024

	logger.Output("SERVER", "Starting server!")

	/* Learning!
	signal 0xc0000005: This is a Windows-specific error indicating an access violation (attempting to access memory that is not valid).
	addr=0x28: This is the memory address that the program tried to access, which is invalid.
	pc=0x7ff7aa2eb870: This is the program counter (instruction pointer) at the time of the crash.
	goroutine 82: The crash occurred in the ServerRun function, which was called in a goroutine.
	*/
	if serverinstance == nil {
		logger.Debug("DEBUG", "Failed to start server! Setting the maintain to be false")
		alive <- false
		return
	}

	// Main components of setting the server up
	go createserverinstance() // Create the instance for listen and serve
	go createudplistener()    // Create a udp listener
	go servershutdownflag()   // Do shutdown work when closing
	go broadcasttonodes()     // Listener to broadcast to nodes

	logger.Debug("DEBUG", "Server has started!")
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
		serverinstance.buffer = make([]byte, 1024)
	}

	return conCreate
}

func createudplistener() {

	logger := log.GetInstance()
	con := createudpcon(8080, "0.0.0.0", true)
	if con == nil { // TODO: Add handler for this
		logger.Debug("SERVER | ERROR", "Connection could not be created!")
		return
	}
	defer con.Close()

	for serverinstance != nil {
		n, addr, err := con.ReadFromUDP(serverinstance.buffer)
		if err != nil {
			fmt.Println("ERROR: Could not read from UDP: " + err.Error())
		} else {
			// Check duplicates
			_, exists := serverinstance.pingpool[strings.Split(addr.String(), ":")[0]]
			senderhostname, errhostname := net.LookupHost(addr.IP.String())
			if !exists {

				if errhostname != nil {
					// Store [address] = hostname
					serverinstance.pingpool[strings.Split(addr.String(), ":")[0]] = strings.Join(senderhostname, " ")
				} else {
					fmt.Printf("Could not resolve hostname!\n")
					serverinstance.pingpool[strings.Split(addr.String(), ":")[0]] = ""
				}
			}

			fmt.Printf("Received response from %s: %s\n", addr.String(), string(serverinstance.buffer[:n]))
		}
	}
}

func servershutdownflag() {
	logger := log.GetInstance()
	waitState := <-serverinstance.maintainsignal
	if !waitState {

		// Shutdown the server
		if err := serverinstance.srv.Shutdown(context.Background()); err != nil { //CTX in this case is not a copy
			logger.Debug("DEBUG", fmt.Sprintf("Server Shutdown Failed:%+v", err))
		}

		logger.Debug("DEBUG", "Server has been stopped")

		serverinstance = nil
		//delete(http.DefaultServeMux.Handle(), "/") Interesting implementation in this method though
	}
}

func createserverinstance() {
	logger := log.GetInstance()
	logger.Debug("DEBUG", fmt.Sprintf("Starting the http, instance %v", serverinstance))
	err := serverinstance.srv.ListenAndServe() // Not http listen and serve, we have our own server
	if err == nil {
		logger.Output("ERROR", "Error in starting the server")
	}
}

func broadcasttonodes() {
	logger := log.GetInstance()
	con := createudpcon(8080, "255.255.255.255", false)

	for serverinstance != nil {
		for serverinstance.broadcasting {
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
		}

		time.Sleep(time.Second * 2)
	}

	con.Close()
}
