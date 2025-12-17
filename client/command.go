package client

import "strings"

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
	for {

		if <-exit {
			return
		}
	}
}

func (c *Command) inbox(exit chan bool) {
	for {

		if <-exit {
			return
		}
	}
}

func (c *Command) draft(exit chan bool) {
	for {

		if <-exit {
			return
		}
	}
}

func (c *Command) util(exit chan bool) {
	for {

		if <-exit {
			return
		}
	}
}

// Main redirect method
func (c *Command) redirect(exit chan bool) func(exit chan bool) {

	funcmapping := map[string]func(chan bool){
		"help":  c.help,
		"inbox": c.inbox,
		"draft": c.draft,
		"util":  c.util,
	}

	if method, exists := funcmapping[c.command]; exists {
		return method
	}

	return func(exit chan bool) {}
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
	if !cmap[c.command] {
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
