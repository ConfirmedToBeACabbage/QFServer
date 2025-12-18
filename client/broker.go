package client

import "fmt"

// Workers for listening to broadcasts
type wbroker struct {
	cmethodsig func(exit chan bool) // The method
	name       string
	status     string
	start      chan bool
	exit       chan bool // Channel for exiting
}

type wbrokercontroller struct {
	wbrokerlist map[string]wbroker // The list of all the workers
	error       bool
	message     string
}

// Add a worker
// Learning: We need to make sure that we're not using a copy of wbrokercontroller but the direct reference
func (w *wbrokercontroller) addworker(c *Command) bool {

	newworker := &wbroker{}
	newworker.exit = make(chan bool)
	newworker.exit <- false
	// The new worker has the method sig assigned while also passing the exit channel. Which it holds in its own structure too.
	newworker.cmethodsig = c.redirect(newworker.exit)
	newworker.name = c.command
	newworker.status = "WORKER: Currently being configured"

	w.message = "STATUS: Configuring new worker"

	// Check for duplicate
	for i := range w.wbrokerlist {
		worker := w.wbrokerlist[i]
		if worker.name == newworker.name {
			w.error = true
			w.message = "ERROR: Broker cannot add duplicate workers"
		}
	}

	// No error then we should just add it to the list and then start the method
	if !w.error {

		// Adding the new worker to the broker controller list
		w.wbrokerlist[newworker.name] = *newworker

		// Signal that the worker should start
		newworker.start = make(chan bool)
		newworker.start <- true

	} else {
		return w.error
	}

	return w.error
}

// Check status
func (w *wbrokercontroller) status(name string) []string {
	worker := w.wbrokerlist[name]

	return []string{worker.name, fmt.Sprintf("%v", worker.start), worker.status}
}

// The main broker routine
func (w *wbrokercontroller) maintain(exitmaintain chan bool) {

	for {
		for i := range w.wbrokerlist {
			worker := w.wbrokerlist[i]

			select {
			case chkstart := <-worker.start:
				if chkstart {
					go worker.cmethodsig(worker.exit)
				}
			case chkend := <-worker.exit:
				if chkend {
					worker.exit <- true
					worker.start <- false
				}
			case <-exitmaintain:
				return
			}
		}
	}

}

// The init for a broker
func InitBroker() *wbrokercontroller {
	broker := &wbrokercontroller{}

	// Start the listener for the broker
	exitmaintain := make(chan bool)
	go broker.maintain(exitmaintain)

	return broker
}
