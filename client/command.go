package client

import (
	"fmt"
	"strings"
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
	fmt.Printf("\nMESSAGE:%s\n%s\n%s\n%s",
		"***HELP***",
		"Inbox: Show incoming mail on LAN (inbox)",
		"Draft: Draft some message and select a destination on LAN (draft [ip])",
		"Util: Scanning, checking to see where an open receiver sits (util)")

	exit <- true
}

func (c *Command) inbox(exit chan bool) {

}

func (c *Command) draft(exit chan bool) {

}

func (c *Command) util(exit chan bool) {

}

// Main redirect method
func (c *Command) redirect(exit chan bool) func(exit chan bool) {

	cmap := map[string]func(chan bool){
		"help":  c.help,
		"inbox": c.inbox,
		"draft": c.draft,
		"util":  c.util,
	}

	methodcall := strings.TrimSpace(c.command)

	method, ok := cmap[methodcall]
	if ok {
		fmt.Printf("\nLOG: Returning appropriate method")
		return method
	} else {
		fmt.Printf("\nLOG: Could not find the method! %v \n In map of %v", c.command, cmap)
		return func(exit chan bool) { fmt.Printf("\nEMPTY METHOD") }
	}
}

// Parsing
// [command] {-optional-args...}
func Parse(input string) *Command {
	c := &Command{}

	args := strings.Split(input, " ")
	c.command = args[0]

	// Validate check
	cmap := map[string]bool{"help": true, "inbox": true, "draft": true, "util": true, "redirect": true}

	// ERROR
	_, ok := cmap[c.command]
	if ok {
		fmt.Printf("\nLOG: [Error] Command is %s", c.command)
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
