package log

import "fmt"

type OutBuffer struct {
	outputQueue *queue
	CurrModule  string
	OutputClear bool
	shutdown    chan bool
	currInput   string
}

func (ob *OutBuffer) moduleinputtext() {
	switch ob.CurrModule {
	case "DEFAULT":
		fmt.Println("\nQFServer CLI! Type in - Help - to get started.")
	case "SERVERREQ":
		fmt.Println("\n** (C[index] to accept connection or [index] to make a request) ** ")
	default:
		return
	}
}

func (ob *OutBuffer) switchmodule(module string) {
	switch module {
	case "SERVERREQ",
		"DEFAULT":
		ob.CurrModule = module
	default:
		fmt.Printf("OUTPUT: Invalid module %s, switching to default\n", module)
		ob.CurrModule = "DEFAULT"
		return
	}
}

// Init the buffer
func (ob *OutBuffer) Init() {
	ob.OutputClear = true
	ob.CurrModule = "DEFAULT"

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
