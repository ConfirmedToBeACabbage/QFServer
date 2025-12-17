package client

// Workers for listening to broadcasts
type wbroker struct {
	cmethodsig func(exit chan bool) // The method
	name       string
	status     string
	start      chan bool
	exit       chan bool // Channel for exiting
}

type wbrokercontroller struct {
	wbrokerlist []wbroker // The list of all the workers
	error       bool
	message     string
}

// Add a worker
// Learning: We need to make sure that we're not using a copy of wbrokercontroller but the direct reference
func (w *wbrokercontroller) addworker(c *Command) bool {

	newworker := &wbroker{}
	newworker.exit = make(chan bool)
	newworker.cmethodsig = c.redirect(newworker.exit) // Running the command redirect to give us the function we're using
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
		w.wbrokerlist = append(w.wbrokerlist, *newworker)

		// Signal that the worker should start
		newworker.start = make(chan bool)
		newworker.start <- true

	} else {
		return w.error
	}

	return w.error
}

// Remove a worker
func (w *wbrokercontroller) removeworker() {

}

// Check status
func (w *wbrokercontroller) status() {

}

func (w *wbrokercontroller) maintain(exitmaintain chan bool) {

	for {
		for i := range w.wbrokerlist {
			worker := w.wbrokerlist[i]

			select {
			case <-worker.start:
				go worker.cmethodsig(worker.exit)
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
