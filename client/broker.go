package Client

import (
	"fmt"
	"sync"
)

// Workers for listening to broadcasts
type broadcastlistenerW struct {
	source chan interface{}
	quit   chan interface{}
}

func (blw *broadcastlistenerW) Start() {
	blw.source = make(chan interface{}, 10)

	go func() {
		for {
			select {
			case msg := <-blw.source: // What to do with the message
				fmt.Printf(`%s`, msg)
			case <-blw.quit: // What to do when quitting
				return
			}
		}
	}()
}

// Thread safe Worker Group
type tslice struct {
	sync.Mutex
	workers []*broadcastlistenerW
}
