package log

import "fmt"

type OutBuffer struct {
	outputQueue *queue
	OutputClear bool
	shutdown    chan bool
}

// Init the buffer
func (ob *OutBuffer) Init() {
	ob.OutputClear = true

	// Init the queue
	QueueInit()
	ob.outputQueue = globalqueue

	go func() {
		for {
			select {
			default:

				queueOut, ok := ob.outputQueue.Dequeue()
				if ok {
					fmt.Printf("OUTPUT [%s]%s \n", queueOut.id, queueOut.message)
				}

				ob.checkclear()
			case <-ob.shutdown:
				return
			}
		}
	}()

}

// Add to output
func (ob *OutBuffer) addtooutput(log logMessage) {
	ob.outputQueue.Enqueue(log)
}

// // Clear the output buffer
// func (ob *OutBuffer) clearoutput() {
// 	QueueReset()
// 	ob.OutputClear = true
// }

// Check if the output is clear
func (ob *OutBuffer) checkclear() bool {
	ob.OutputClear = ob.outputQueue == nil || len(ob.outputQueue.items) == 0
	return ob.OutputClear
}
