package client

import (
	"fmt"
	"strings"
	"sync"

	"github.com/QFServer/log"
	"github.com/QFServer/server"
)

type CMethodSig interface {
	help(chan bool)
	inbox(chan bool)
	draft(chan bool)
	util(chan bool)
	redirect(chan bool)
}

type Command struct {
	command   string
	args      []string
	message   string
	giveerror bool

	mu sync.Mutex
}

// Command methods signed by commandcontrol
func (c *Command) help() {
	fmt.Printf("\n%s\n%s\n%s\n%s\n%s\n%s\n%s",
		"\n***HELP***",
		"Inbox: Show incoming mail on LAN (inbox)",
		"Draft: Draft some message and select a destination on LAN (draft [ip])",
		"Util: Scanning, checking to see where an open receiver sits (util)",
		"      - server open: This would start the server and get it ready for scanning",
		"      - server broadcast: This would start broadcasting your server. Other node pools can pick it up and add it on LAN",
		"Quit: This will quit the program\n")
}

func (c *Command) inbox() {
	// We would have to just check for connections pooled?
	// Then when the connections are pooled we can either open them with a token
	// Or choose to receive them. We can also see the contents before we download
	fmt.Println("Inbox!")
}

func (c *Command) draft() {
	// This is where we would have a pool of known nodes on the network.
	fmt.Println("Exit!")
}

func (c *Command) util() {
	// Util should have a couple functions; We're starting with scanning and locating
	// possible receivers
	fmt.Println("Utility!")
}

// SERVER ARGS

// SERVER: Listener; This would start the broadcast listener
func (c *Command) srvbroadcast(alive chan bool) {
	<-alive // Wait for alive

	fmt.Println("Changing the server state for broadcasting!")
	server.BroadcastStateChange()
	alive <- false
}

// SERVER: Open; This should open the server
func (c *Command) srvopen(alive chan bool) {
	<-alive // Wait for alive

	fmt.Println("SERVER: In the command method to begin server!")
	server.ServerRun(alive)
	fmt.Println("SERVER: Server goroutine has begun. It should be created soon!")
}

// SERVER: pool; This should show us the pool of users which we have on lan that we can send to
func (c *Command) srvpool(alive chan bool) {
	<-alive // Wait for alive

	fmt.Println("Showing the pool!")
	poollist := server.ServerInitSingleton().GetPingPool()

	fmt.Println("Address | Host")
	for i, v := range poollist {
		fmt.Println(i + " | " + v)
	}

	alive <- false
}

// Check if the server is alive
func (c *Command) srvcheckalive(alive chan bool) {
	<-alive // Wait for alive

	serveractive := server.CheckServerAlive()
	fmt.Printf("\nSERVER ACTIVE STATUS: %v\n", serveractive)

	alive <- false
}

// Main redirect method
func (c *Command) redirect(alive chan bool) (func(chan bool), func(chan bool)) {
	c.mu.Lock()
	defer c.mu.Unlock()

	logger := log.GetInstance()

	cmapstart := map[string]func(){
		"help":  c.help,
		"inbox": c.inbox,
		"draft": c.draft,
		"util":  c.util,
	}

	cmapserver := map[string]func(chan bool){
		"broadcast": c.srvbroadcast, // Broadcast our client
		"open":      c.srvopen,
		"pool":      c.srvpool,
		"alive":     c.srvcheckalive,
	}

	cmaprouteutil := map[string]map[string]func(chan bool){
		"server": cmapserver,
	}

	// Method call
	methodcall := strings.TrimSpace(c.command)

	// Setting up the return list
	returnlist := make([]func(controller chan bool), 2)
	returnlist[0] = func(alive chan bool) { fmt.Print("\nEMPTY\n") }
	returnlist[1] = func(alive chan bool) { fmt.Print("\nEMPTY\n") }

	// TODO: this is just a tree traverse, I should automate this away. Commands might grow in length.
	// AKA make a command tree
	// Lower priority for now
	if len(c.args) > 1 {
		argcall := strings.TrimSpace(c.args[0])
		fmt.Printf("DEBUG: This is the first argcall %v\n", argcall)
		route, okroute := cmaprouteutil[argcall]
		logger.Debug("DEBUG", "We have gotten past the route!")
		logger.Debug("DEBUG", fmt.Sprintf("%v\n", route))

		if okroute {
			furtherargcall := strings.TrimSpace(c.args[1])
			logger.Debug("DEBUG", ""+furtherargcall)
			furthercommand, okcommand := route[furtherargcall]

			if okcommand {
				returnlist[1] = func(alive chan bool) { furthercommand(alive) }
			}
		}
	}

	startmethod, okstart := cmapstart[methodcall]

	if okstart {
		returnlist[0] = func(alive chan bool) { startmethod() }
	}

	fmt.Println("Got all the commands and stuff!")

	return returnlist[0], returnlist[1]
}

// Parsing
// [command] {-optional-args...}
func Parse(input string) *Command {

	logger := log.GetInstance()

	c := &Command{}

	args := strings.Split(input, " ")
	c.command = strings.TrimSpace(args[0])

	// Validate check
	cmap := map[string]bool{"help": true, "inbox": true, "draft": true, "util": true, "redirect": true}

	// ERROR
	_, ok := cmap[strings.TrimSpace(c.command)]
	if !ok {
		logger.Store("COMMAND", "Could not find command "+c.command)
		c.message = "ERROR: Command not found! " + c.command
		c.giveerror = true
		return c
	}

	for i := range len(args) {
		fmt.Println(i)
		if i >= 1 {
			c.args = append(c.args, args[i])
		}
	}

	return c
}
