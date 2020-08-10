package events

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
	eventStream        <-chan string
	load               uint64
	loadCallback       func(uint64)
	newMachineCallback func(string)
}

func NewEventLoop() (*EventLoop, func(string)) {
	var stream = make(chan string, 100)
	var loop = &EventLoop{
		eventStream:        stream,
		load:               0,
		loadCallback:       func(uint64) {},
		newMachineCallback: func(string) {},
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

func (loop *EventLoop) SetMachineCallback(f func(string)) {
	loop.newMachineCallback = f
}

func (loop *EventLoop) watchEvents() {
	for e := range loop.eventStream {
		// Check what kind of event this is.
		e = strings.ToLower(e)
		loop.DispatchEvent(e)
	}
}

func (loop *EventLoop) DispatchEvent(event string) {
	switch {
	case strings.HasPrefix(event, "new machine"):
		loop.HandleNewMachine(event)
	case strings.HasPrefix(event, "set load "):
		loop.HandleSetLoad(event)
	}
}

func (loop *EventLoop) HandleSetLoad(event string) {
	// Strip out the "set load" substring.
	var edited = strings.TrimPrefix(event, "set load ")
	// Extract the entered load.
	var load, err = strconv.ParseUint(edited, 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	loop.load = load
	loop.sendLoad()
}

func (loop *EventLoop) HandleNewMachine(event string) {
	// Capture the machine name by stripping away the prefix
	var machineName = strings.TrimPrefix(event, "set load ")
	machineName = strings.TrimSpace(machineName)
	loop.newMachineCallback(machineName)
}

func (loop *EventLoop) sendLoad() {
	loop.loadCallback(loop.load)
}
