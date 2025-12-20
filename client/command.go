package client

import (
	"fmt"
	"strings"

	"github.com/QFServer/log"
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
	fmt.Printf("\n%s\n%s\n%s\n%s\n%s",
		"\n***HELP***",
		"Inbox: Show incoming mail on LAN (inbox)",
		"Draft: Draft some message and select a destination on LAN (draft [ip])",
		"Util: Scanning, checking to see where an open receiver sits (util)",
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

// Main redirect method
func (c *Command) redirect(exit chan bool) func(exit chan bool) {

	logger := log.GetInstance()

	cmap := map[string]func(chan bool){
		"help":  c.help,
		"inbox": c.inbox,
		"draft": c.draft,
		"util":  c.util,
	}

	methodcall := strings.TrimSpace(c.command)

	method, ok := cmap[methodcall]
	if ok {
		logger.Store("COMMAND", "Returning appropriate method")
		return method
	} else {
		logger.Store("COMMAND", "Could not find the method! "+c.command+" In map of "+fmt.Sprint(cmap))
		return func(exit chan bool) { fmt.Printf("\nNIL") }
	}
}

// Parsing
// [command] {-optional-args...}
func Parse(input string) *Command {

	logger := log.GetInstance()

	c := &Command{}

	args := strings.Split(input, " ")
	c.command = args[0]

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
