package main

import (
	"fmt"
	"strconv"
	"strings"
)

// What kind of events does the EventLoop handle?
// 1. It handles changes of load.
// 2. It handles requests to make a new machine.
// 3. It handles requests to add a service.
// 4. It handles requests to add a neighbor.
type EventLoop struct {
	eventStream  <-chan string
	load         uint64
	loadCallback func(uint64)
}

func NewEventLoop() (*EventLoop, func(string)) {
	var stream = make(chan string, 100)
	var loop = &EventLoop{
		eventStream:  stream,
		load:         0,
		loadCallback: func(uint64) {},
	}

	go loop.watchEvents()

	var callback = func(contents string) {
		stream <- contents
	}

	return loop, callback
}

func (loop *EventLoop) SetLoadCallback(f func(uint64)) {
	loop.loadCallback = f
}

func (loop *EventLoop) watchEvents() {
	for e := range loop.eventStream {
		// Check what kind of event this is.
		e = strings.ToLower(e)
		switch {
		case strings.HasPrefix(e, "new machine"):
			// make a new machine
		case strings.HasPrefix(e, "set load "):
			// Strip out the "set load" substring.
			var edited = strings.TrimPrefix(e, "set load ")
			// Extract the entered load.
			var load, err = strconv.ParseUint(edited, 10, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			loop.load = load
			loop.sendLoad()
		}
	}
}

func (loop *EventLoop) sendLoad() {
	loop.loadCallback(loop.load)
}
