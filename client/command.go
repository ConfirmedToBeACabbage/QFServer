package Client

import "strings"

type commandcontrol interface {
	help()
	inbox()
	draft()
	util()
	redirect()
}

type command struct {
	command string
	args    []string
	message string
}

// Command methods signed by commandcontrol
func (c command) help()  {}
func (c command) inbox() {}
func (c command) draft() {}
func (c command) util()  {}

// Main redirect method
func (c command) redirect() {

}

// Parsing
// [command] {-optional-args...}
func Parse(input string) command {
	c := &command{}

	args := strings.Split(input, " ")
	c.command = args[0]

	// Validate check
	cmap := map[string]bool{"help": true, "inbox": true, "draft": true, "util": true, "redirect": true}

	// ERROR
	if !cmap[c.command] {
		c.message = "ERROR: Command not found!"
		return *c
	}

	for i := range len(args) {
		if i >= 1 {
			c.args = append(c.args, args[i])
		}
	}

	return *c
}

// Redirect method
func Direct(cc commandcontrol) {
	cc.redirect()
}
