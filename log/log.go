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
//	Set an output or something if needed
type logdb struct {
	logdictstorage map[string][]string
	mu             sync.Mutex
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

// Init the logs only as a singleton
func GetInstance() *logdb {
	once.Do(func() {
		instance = &logdb{
			logdictstorage: make(map[string][]string),
		}
	})
	return instance
}
