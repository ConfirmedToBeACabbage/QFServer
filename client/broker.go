package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/QFServer/log"
)

// Workers for listening to broadcasts
// Learnings: Data races
type wbroker struct {
	cmethodsig      func(exit chan bool)     // The start method
	cmethodmaintain func(maintain chan bool) // The maintain method
	name            string
	status          string
	alive           chan bool  // Channel to make sure it's alive
	start           bool       // Channel for starting
	maintain        bool       // Channel for maintain
	exit            bool       // Channel for exiting
	mu              sync.Mutex // Added for synch and protecting against data races
}

func (w *wbroker) setStatus(status string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.status = status
}

func (w *wbroker) getStatus() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.status
}

type wbrokercontroller struct {
	wbrokerlist  map[string]*wbroker // The list of all the workers
	error        bool
	message      string
	mu           sync.Mutex
	exitmaintain chan bool
	mur          sync.RWMutex
}

// This is simply a graceful shutdown. Since all the channels are basically associated with the
// relationship between the broker and the workers, we can in a centralized fashion shut it all
// down
func (w *wbrokercontroller) gracefulshutdown(shutdown chan bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.exitmaintain <- true

	for i := range w.wbrokerlist {
		worker := w.wbrokerlist[i]

		close(worker.alive)
	}

	shutdown <- true
}

// This is fun. It's like working with lego bricks but your computer does stuff. Falling back in love with a hobby is great.

// Add a worker
// Learning: We need to make sure that we're not using a copy of wbrokercontroller but the direct reference
func (w *wbrokercontroller) configureworker(c *Command) bool {

	w.mu.Lock()
	defer w.mu.Unlock()

	logger := log.GetInstance()
	logger.Store("BROKER", "Configuring a new worker for command:"+c.command)
	newworker := &wbroker{}
	logger.Store("BROKER", "New Worker has been instantiated")
	newworker.alive = make(chan bool)
	newworker.exit = false
	newworker.start = true
	newworker.maintain = true
	logger.Debug("DEBUG", "Setup the channels!")
	logger.Store("BROKER", "Worker has a channel created")
	logger.Store("BROKER", "Worker channel is setup")
	// The new worker has the method sig assigned while also passing the exit channel. Which it holds in its own structure too.
	logger.Debug("DEBUG", "Starting the redirect")

	start, maintain := c.redirect(newworker.alive)
	newworker.cmethodsig = start
	newworker.cmethodmaintain = maintain

	logger.Debug("DEBUG", "We have gotten past redirect")

	logger.Store("BROKER", "Worker has method assigned")
	namebuilder := c.command // A new name builder which just uses the whole name
	for i := range c.args {
		namebuilder += c.args[i]
		logger.Debug("DEBUG", "The name we're building! "+namebuilder)
	}
	newworker.name = namebuilder

	logger.Store("BROKER", "Worker has name assigned")
	newworker.status = "WORKER: Currently being configured"
	newworker.setStatus("STATUS: Configuring new worker")

	logger.Store("BROKER", "We have completed a new worker configuration")
	logger.Debug("DEBUG", "We have done the configuration!")

	// Check for duplicate
	_, duplicate := w.wbrokerlist[newworker.name]
	if duplicate {
		w.error = true
		w.message = "ERROR: Broker cannot add duplicate workers"
		logger.Debug("DEBUG", "We have a duplicate worker name"+newworker.name)
		return w.error
	} else {

		logger.Store("BROKER", "Adding to the worker list")
		w.error = false
		logger.Store("BROKER", "All done! Error "+fmt.Sprint(w.error))
		newworker.setStatus("[STATUS] Ready to init!")

		// Learning: This below would cause a freeze if unbuffered channel
		// Channels in go require there to be a sender and receiver. So because there is no receiver,
		// it will be an unbuffered channel.
		// What we can do is buffer it above by doing make(chan bool, 1)

		w.wbrokerlist[newworker.name] = newworker
	}

	return w.error
}

// // Check status
// func (w *wbrokercontroller) status(name string) []string {
// 	worker := w.wbrokerlist[name]

// 	return []string{worker.name, fmt.Sprintf("%v", worker.start), worker.status}
// }

// Return latest boolean from a boolean channel
// Must be buffered!
// Learning: Channels are queues. When you drain from a channel it becomes empty! AKa return false!
func (w *wbroker) controlworkerstate() {
	select {
	case workeralive := <-w.alive:
		if !workeralive {
			w.exit = true
		}
	default:
	}
}

// The main broker routine
func (w *wbrokercontroller) maintain(readyforinput chan bool) {

	logger := log.GetInstance()

	go func() {

		for {
			for i := range w.wbrokerlist {

				time.Sleep(time.Millisecond * 200) // Add a small delay

				worker := w.wbrokerlist[i]

				logger.Store("BROKER", "Checking worker: "+worker.status)

				// Learning: Select is good for not blocking and waiting for channel
				// Learning: Select statements are Randomly selecting cases, we might need for loops here to enforce order
				// Learning: Use sync.RWMutex: If the wbrokerlist is read-heavy, you can use a sync.RWMutex to allow multiple readers while still protecting writes.

				// Now we can do sequentially since select statements select at random
				if worker.start {
					worker.setStatus("STATUS: Running a start method for " + worker.name)
					logger.Debug("BROKER", worker.status)
					go worker.cmethodsig(worker.alive)
					worker.start = false // To make sure we don't restart it
				}

				if worker.maintain {
					worker.setStatus("STATUS: Running a maintain method for " + worker.name)
					logger.Debug("BROKER", worker.status)
					go worker.cmethodmaintain(worker.alive)
					worker.alive <- true
					go worker.controlworkerstate()
					worker.maintain = false
				}

				readyforinput <- true

				if worker.exit {
					worker.setStatus("STATUS: Exiting the worker " + worker.name)
					logger.Debug("BROKER", worker.status)
					// Closing the channels used in the worker
					close(worker.alive)
					delete(w.wbrokerlist, worker.name) // Deleting from the list
				}
			}
		}

	}()

	<-w.exitmaintain
	close(w.exitmaintain)

}

// The init for a broker
func InitBroker(readyforinput chan bool) *wbrokercontroller {
	broker := &wbrokercontroller{
		// Learning: Make does a couple things
		// 1. Allocate memory
		// 2. Initialize data structure
		// Good practice when working with slices, maps, and channels
		wbrokerlist:  make(map[string]*wbroker),
		exitmaintain: make(chan bool),
	}

	// Start the listener for the broker
	go broker.maintain(readyforinput)

	fmt.Printf("\nLOG: Done!")

	return broker
}
