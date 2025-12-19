package client

import (
	"fmt"
	"sync"
	"time"
)

// Workers for listening to broadcasts
// Learnings: Data races
type wbroker struct {
	cmethodsig func(exit chan bool) // The method
	name       string
	status     string
	start      chan bool
	exit       chan bool  // Channel for exiting
	mu         sync.Mutex // Added for synch and protecting against data races
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
	wbrokerlist map[string]*wbroker // The list of all the workers
	error       bool
	message     string
	mu          sync.Mutex
}

// Add a worker
// Learning: We need to make sure that we're not using a copy of wbrokercontroller but the direct reference
func (w *wbrokercontroller) configureworker(c *Command) bool {

	w.mu.Lock()
	defer w.mu.Unlock()

	fmt.Printf("\n[BROKER] LOG: Configuring a new worker for command: %s\n", c.command)

	newworker := &wbroker{}

	fmt.Printf("\n[BROKER] LOG: New worker has been instantiated")

	newworker.exit = make(chan bool, 1)

	fmt.Printf("\n[BROKER] LOG: Worker has a channel created")

	// Learning: This below would cause a freeze if unbuffered channel
	newworker.exit <- false
	// Channels in go require there to be a sender and receiver. So because there is no receiver,
	// it will be an unbuffered channel.
	// What we can do is buffer it above by doing make(chan bool, 1)

	fmt.Printf("\n[BROKER] LOG: Worker channel is setup")

	// The new worker has the method sig assigned while also passing the exit channel. Which it holds in its own structure too.
	newworker.cmethodsig = c.redirect(newworker.exit)

	fmt.Printf("\n[BROKER] LOG: Worker has method assigned")

	newworker.name = c.command

	fmt.Printf("\n[BROKER] LOG: Worker has name assigned")

	newworker.status = "WORKER: Currently being configured"

	newworker.setStatus("STATUS: Configuring new worker")

	fmt.Printf("\n[BROKER] LOG: We have completed a new worker configuration")

	// Check for duplicate
	_, duplicate := w.wbrokerlist[newworker.name]
	if duplicate {
		w.error = true
		w.message = "ERROR: Broker cannot add duplicate workers"
		return w.error
	} else {

		fmt.Printf("\n[BROKER] LOG: Adding to the worker list")

		// Adding the new worker to the broker controller list
		// Learning: We have to actually initialize the map. It's nil right now, we will do that
		// In the init portion of the init broker
		w.wbrokerlist[newworker.name] = newworker

		// Signal that the worker should start
		newworker.start = make(chan bool, 1)
		newworker.start <- true

		w.error = false

		fmt.Printf("\n[BROKER] LOG: All done! Error %v", w.error)

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
func (w *wbrokercontroller) maintain(exitmaintain chan bool) {

	for {
		select {
		case <-exitmaintain:
			return
		default:
			for i := range w.wbrokerlist {
				worker := w.wbrokerlist[i]

				fmt.Printf("\nLOG: [Broker] Checking worker: %v", worker.status)

				// Learning: Select is good for not blocking and waiting for channel
				select {
				case start := <-worker.start:
					if start {
						worker.setStatus("STATUS: Beginning the goroutine")
						fmt.Printf("\nLOG: [Broker] Worker status: %v for name %v and method %v", worker.getStatus(), worker.name, worker.cmethodsig)
						go worker.cmethodsig(worker.exit)
					}
				case exitworker := <-worker.exit:
					if exitworker {
						worker.setStatus("STATUS: Exiting the worker!")
						delete(w.wbrokerlist, worker.name) // Deleting from the list
					}
				default:
					continue
				}
			}

			time.Sleep(time.Second * 5)
			fmt.Printf("\nLOG: [Broker] Checking workers...")
		}
	}

}

// The init for a broker
func InitBroker() *wbrokercontroller {
	broker := &wbrokercontroller{
		// Learning: Make does a couple things
		// 1. Allocate memory
		// 2. Initialize data structure
		// Good practice when working with slices, maps, and channels
		wbrokerlist: make(map[string]*wbroker),
	}

	// Start the listener for the broker
	exitmaintain := make(chan bool)
	go broker.maintain(exitmaintain)

	fmt.Printf("\nLOG: Done!")

	return broker
}
