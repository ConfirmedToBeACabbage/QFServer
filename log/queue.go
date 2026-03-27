package log

import (
	"sync"
)

type queue struct {
	items []logMessage
	mu    sync.Mutex
}

// Enqueue adds an item to the end of the queue
func (q *queue) Enqueue(item logMessage) {
	q.mu.Lock()
	defer q.mu.Unlock()             // Defering so it only is unlocked after method end
	q.items = append(q.items, item) // Appending the item to the end of the queue
}

// Dequeue removes and returns the item at the front of the queue
func (q *queue) Dequeue() (logMessage, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return logMessage{}, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Generate a queue and then return it's pointer
func GenerateQueue(itype logMessage) *queue {
	rqueue := &queue{}
	return rqueue
}

var globalqueue *queue
var onceQ sync.Once

func QueueReset() {
	if globalqueue != nil {
		globalqueue = GenerateQueue(logMessage{})
	}
}

func QueueInit() {
	onceQ.Do(func() {
		globalqueue = GenerateQueue(logMessage{})
	})
}
