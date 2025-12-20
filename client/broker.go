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
	start           chan bool  // Channel for starting
	maintain        chan bool  // Channel for maintain
	exit            chan bool  // Channel for exiting
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

		worker.exit <- true
		worker.start <- false

		close(worker.exit)
		close(worker.start)
	}

	shutdown <- true
}

// Add a worker
// Learning: We need to make sure that we're not using a copy of wbrokercontroller but the direct reference
func (w *wbrokercontroller) configureworker(c *Command) bool {

	w.mu.Lock()
	defer w.mu.Unlock()

	logger := log.GetInstance()
	logger.Store("BROKER", "Configuring a new worker for command:"+c.command)
	newworker := &wbroker{}
	logger.Store("BROKER", "New Worker has been instantiated")
	newworker.exit = make(chan bool, 1)
	newworker.start = make(chan bool, 1)
	newworker.maintain = make(chan bool, 1)
	logger.Store("BROKER", "Worker has a channel created")

	// Learning: This below would cause a freeze if unbuffered channel
	newworker.exit <- false
	// Channels in go require there to be a sender and receiver. So because there is no receiver,
	// it will be an unbuffered channel.
	// What we can do is buffer it above by doing make(chan bool, 1)

	logger.Store("BROKER", "Worker channel is setup")
	// The new worker has the method sig assigned while also passing the exit channel. Which it holds in its own structure too.
	start, maintain := c.redirect(newworker.exit, newworker.maintain)
	newworker.cmethodsig = start
	newworker.cmethodmaintain = maintain

	// Should I have a shutdown method also? AKA listener needs to shutdown
	//newworker.cmethodexit = c.redirectexit(newworker.shutdown)

	logger.Store("BROKER", "Worker has method assigned")
	newworker.name = c.command

	logger.Store("BROKER", "Worker has name assigned")
	newworker.status = "WORKER: Currently being configured"
	newworker.setStatus("STATUS: Configuring new worker")

	logger.Store("BROKER", "We have completed a new worker configuration")

	// Check for duplicate
	_, duplicate := w.wbrokerlist[newworker.name]
	if duplicate {
		w.error = true
		w.message = "ERROR: Broker cannot add duplicate workers"
		return w.error
	} else {

		logger.Store("BROKER", "Adding to the worker list")

		// Adding the new worker to the broker controller list
		// Learning: We have to actually initialize the map. It's nil right now, we will do that
		// In the init portion of the init broker
		w.wbrokerlist[newworker.name] = newworker

		// Signal that the worker should start
		newworker.start = make(chan bool, 1)
		newworker.start <- true

		w.error = false

		logger.Store("BROKER", "All done! Error "+fmt.Sprint(w.error))

		newworker.setStatus("[STATUS] Ready to init!")
	}

	return w.error
}

// // Check status
// func (w *wbrokercontroller) status(name string) []string {
// 	worker := w.wbrokerlist[name]

// 	return []string{worker.name, fmt.Sprintf("%v", worker.start), worker.status}
// }

// The main broker routine
func (w *wbrokercontroller) maintain(readyforinput chan bool) {

	logger := log.GetInstance()

	go func() {

		for {
			for i := range w.wbrokerlist {
				worker := w.wbrokerlist[i]

				logger.Store("BROKER", "Checking worker: "+worker.status)

				// Learning: Select is good for not blocking and waiting for channel
				select {
				case start := <-worker.start:
					if start {
						worker.setStatus("STATUS: Beginning the goroutine")
						logger.Store("BROKER", "Worker status "+worker.status+" for name "+worker.name)
						go worker.cmethodsig(worker.exit)
						worker.start <- false // To make sure we don't restart it
					}
				case exitworker := <-worker.exit:
					if exitworker {
						worker.setStatus("STATUS: Exiting the worker!")
						delete(w.wbrokerlist, worker.name) // Deleting from the list

						// Closing the channels used in the worker
						close(worker.exit)
						close(worker.start)

						// Ready for next input (It has been completed)
						readyforinput <- true
					}
				default:
					time.Sleep(time.Millisecond * 100) // Add a small delay
				}
			}

			logger.Store("BROKER", "Checking workers...")
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
