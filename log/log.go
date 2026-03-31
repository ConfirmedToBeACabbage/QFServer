package log

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

// A structure which stores the logs
// We store logs in a dictionary [LogIDName]:[Collection of logs, sequentially]
// Function to call on the logs from these different ID's
// Extra feat: Call on specific times
//
// Set an output or something if needed
type logMessage struct {
	id      string
	message string
}

type logdb struct {
	logdictstorage    map[string][]string
	outputBuffer      *OutBuffer
	inputcheckchannel chan bool
	debuggeralive     bool
	debuglogshow      bool
	mu                sync.Mutex
}

// We should only have this as a singleton
// Learning: You can group things in a var like this
// Technically they're both vars, it's just for syntax
var (
	instance *logdb
	once     sync.Once
)

// The debug logger goroutine which manages all debug messages
func (l *logdb) BeginDebugLogger() {

	if l.debuggeralive {
		fmt.Println("LOGGER: Logger is already active!")
		return
	}

	l.outputBuffer = &OutBuffer{}
	l.outputBuffer.Init()

	l.debuggeralive = true
	l.debuglogshow = true
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

// Function to see if we need to actually use the input or another module is currently using it
func (l *logdb) CheckModule() string {
	return l.outputBuffer.CurrModule
}

// Switches the module
func (l *logdb) SwitchModule(module string) {
	l.outputBuffer.switchmodule(module)
}

func (l *logdb) InputFromUser() string {

	if l.outputBuffer.checkclear() {
		reader := bufio.NewReader(os.Stdin)

		l.outputBuffer.moduleinputtext() // Get the input text depending on the module
		fmt.Print("> ")
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Print("CRITICAL: Error")
			return ""
		} else {
			return input
		}

	} else {
		return ""
	}
}

func (l *logdb) ReadyForUserInput() bool {
	return l.outputBuffer.checkclear()
}

// This would turn the debugging output off or o
func (l *logdb) DebuggingOutputOnOff() {
	l.debuglogshow = !l.debuglogshow
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

	if l.debuglogshow {
		l.outputBuffer.addtooutput(logMessage{id: id, message: log})
	}
}

// Debug Storing and Printing in sync with the input
func (l *logdb) Output(id string, log string) {
	// We use the normal storing function
	l.Store(id, log)

	l.outputBuffer.addtooutput(logMessage{id: id, message: log})
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
