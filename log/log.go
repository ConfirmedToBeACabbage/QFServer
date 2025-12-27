package log

import (
	"fmt"
	"sync"
)

// A structure which stores the logs
// We store logs in a dictionary [LogIDName]:[Collection of logs, sequentially]
// Function to call on the logs from these different ID's
// Extra feat: Call on specific times
//
// Set an output or something if needed
type logMessage struct {
	id            string
	message       string
	userinputting bool
}

type logdb struct {
	logdictstorage    map[string][]string
	inputchannel      chan logMessage
	inputcheckchannel chan bool
	debuggeralive     bool
	mu                sync.Mutex
}

// Store something in logs according to the id
func (l *logdb) Store(id string, log string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	_, exist := l.logdictstorage[id]
	if !exist {
		l.logdictstorage[id] = make([]string, 0)
	}

	l.logdictstorage[id] = append(l.logdictstorage[id], log)

	return true
}

// Debug Storing and Printing in sync with the input
func (l *logdb) Debug(id string, log string) {
	// We use the normal storing function
	l.Store(id, log)

	localinputcheck := false
	select {
	case checkinput := <-l.inputcheckchannel:
		if checkinput {
			localinputcheck = true
		}
	default:
		localinputcheck = false
	}

	l.inputchannel <- logMessage{id: id, message: log, userinputting: localinputcheck}
}

// The debug logger goroutine which manages all debug messages
func (l *logdb) BeginDebugLogger(inputcheckchannel chan bool) {

	if inputcheckchannel == nil {
		fmt.Printf("ERROR: Cannot start the debug logger! Please provide an appropriate input channel")
		return
	}

	l.SetInputChannel(inputcheckchannel)
	l.inputchannel = make(chan logMessage, 1)

	if !l.debuggeralive {
		go func() {
			for logMessage := range l.inputchannel {
				fmt.Printf("%s LOG: %s\n", logMessage.id, logMessage.message)
				if logMessage.userinputting {
					fmt.Printf("> ")
				}
			}
		}()
	} else {
		fmt.Println("LOGGER: Logger is already active!")
	}

	l.debuggeralive = true
}

// Retreive something in logs according to the id
func (l *logdb) Read(id string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i := range l.logdictstorage[id] {
		log := l.logdictstorage[id][i]
		fmt.Printf(`\n%s LOG: %s`, id, log)
	}
}

// We should only have this as a singleton
// Learning: You can group things in a var like this
// Technically they're both vars, it's just for syntax
var (
	instance *logdb
	once     sync.Once
)

// Setting the input channel
func (l *logdb) SetInputChannel(inputcheckchannel chan bool) {
	l.inputcheckchannel = inputcheckchannel
}

// Init the logs only as a singleton (You can pass an input channel here)
// Although this is just a one time setting operationg. To change the input channel you can set it in SetInputChannel.
func GetInstance() *logdb {

	once.Do(func() {
		instance = &logdb{
			logdictstorage: make(map[string][]string),
		}
	})

	return instance
}
