package client

import (
	"fmt"
	"strings"

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
}

// Command methods signed by commandcontrol
func (c *Command) help(exit chan bool) {
	fmt.Printf("\n%s\n%s\n%s\n%s\n%s\n%s\n%s",
		"\n***HELP***",
		"Inbox: Show incoming mail on LAN (inbox)",
		"Draft: Draft some message and select a destination on LAN (draft [ip])",
		"Util: Scanning, checking to see where an open receiver sits (util)",
		"      - server open: This would start the server and get it ready for scanning",
		"      - server broadcast: This would start broadcasting your server. Other node pools can pick it up and add it on LAN",
		"Quit: This will quit the program\n")

	exit <- true
}

func (c *Command) inbox(exit chan bool) {
	// We would have to just check for connections pooled?
	// Then when the connections are pooled we can either open them with a token
	// Or choose to receive them. We can also see the contents before we download
}

func (c *Command) draft(exit chan bool) {
	// This is where we would have a pool of known nodes on the network.
}

func (c *Command) util(exit chan bool) {
	// Util should have a couple functions; We're starting with scanning and locating
	// possible receivers
}

// SERVER ARGS

// SERVER: Listener; This would start the broadcast listener
func (c *Command) srvbroadcast(maintain chan bool) {

	// Check server
	serveractive := server.CheckServerAlive()
	if !serveractive {
		fmt.Printf("FAIL: You must run 'util server open' to first open the server")
		return
	} else {
		server.ServerInitSingleton().BroadcastStateChange()
	}

	fmt.Printf("SERVER: We have begun broadcasting")

}

// SERVER: Open; This should open the server
func (c *Command) srvopen(maintain chan bool) {

	// Init the server
	serveractive := server.CheckServerAlive()
	if !serveractive {
		server.ServerRun(maintain)
	}

	fmt.Printf("SERVER: The server has now begun running and is open")

}

// SERVER: pool; This should show us the pool of users which we have on lan that we can send to
func (c *Command) srvpool(maintain chan bool) {

	// Check server
	serveractive := server.CheckServerAlive()
	if !serveractive {
		fmt.Printf("FAIL: You must run 'util server open' to first open the server")
		return
	} else {
		poollist := server.ServerInitSingleton().GetPingPool()

		fmt.Println("Address | Host")
		for i, v := range poollist {
			fmt.Println(i + " | " + v)
		}
	}

}

// Main redirect method
func (c *Command) redirect(exit chan bool, maintain chan bool) (func(chan bool), func(chan bool)) {

	// logger := log.GetInstance()

	cmapstart := map[string]func(chan bool){
		"help":  c.help,
		"inbox": c.inbox,
		"draft": c.draft,
		"util":  c.util,
	}

	cmapserver := map[string]func(chan bool){
		"broadcast": c.srvbroadcast, // Broadcast our client
		"open":      c.srvopen,
		"pool":      c.srvpool,
	}

	cmaprouteutil := map[string]map[string]func(chan bool){
		"server": cmapserver,
	}

	// Method call
	methodcall := strings.TrimSpace(c.command)

	// Setting up the return list
	returnlist := make([]func(controller chan bool), 2)
	returnlist[0] = func(exit chan bool) { fmt.Print("\nEMPTY\n") }
	returnlist[1] = func(maintain chan bool) { fmt.Print("\nEMPTY\n") }

	// TODO: this is just a tree traverse, I should automate this away. Commands might grow in length.
	// AKA make a command tree
	// Lower priority for now
	if len(c.args) > 1 {
		argcall := strings.TrimSpace(c.args[1])
		route, okroute := cmaprouteutil[argcall]

		if okroute {
			furtherargcall := strings.TrimSpace(c.args[2])
			furthercommand, okcommand := route[furtherargcall]

			if okcommand {
				returnlist[1] = func(maintain chan bool) { furthercommand(maintain) }
			}
		}
	}

	startmethod, okstart := cmapstart[methodcall]

	if okstart {
		returnlist[0] = func(exit chan bool) { startmethod(exit) }
	}

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
	_, ok := cmap[c.command]
	if ok {
		logger.Store("COMMAND", "Could not find command "+c.command)
		c.message = "ERROR: Command not found!"
		c.giveerror = true
		return c
	}

	for i := range len(args) {
		if i >= 1 {
			c.args = append(c.args, args[i])
		}
	}

	return c
}
