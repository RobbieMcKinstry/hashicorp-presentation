package events

import (
	"fmt"
	"log"
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
	newServiceCallback func(throughput, soft, hard uint64)
}

func NewEventLoop() (*EventLoop, func(string)) {
	var stream = make(chan string, 100)
	var loop = &EventLoop{
		eventStream:        stream,
		load:               0,
		loadCallback:       func(uint64) {},
		newMachineCallback: func(string) {},
		newServiceCallback: func(uint64, uint64, uint64) {},
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

func (loop *EventLoop) SetServiceCallback(f func(uint64, uint64, uint64)) {
	loop.newServiceCallback = f
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
	case strings.HasPrefix(event, "new service "):
		loop.HandleNewService(event)
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

func (loop *EventLoop) HandleNewService(event string) {
	// Capture the machine name by stripping away the prefix
	var serviceParams = strings.TrimPrefix(event, "new service ")
	serviceParams = strings.TrimSpace(serviceParams)
	// Now, split this into four parameters:
	// Throughput, Soft, Hard
	var params = strings.Fields(serviceParams)
	if len(params) != 3 {
		fmt.Println("Expected three arguments")
		ExitOnError(fmt.Errorf("Expected 3 parameters, found %v", len(params)))
	}
	throughput, err := strconv.ParseUint(params[0], 10, 64)
	ExitOnError(err)
	softLimit, err := strconv.ParseUint(params[1], 10, 64)
	ExitOnError(err)
	hardLimit, err := strconv.ParseUint(params[2], 10, 64)
	ExitOnError(err)
	loop.newServiceCallback(throughput, softLimit, hardLimit)
}

func (loop *EventLoop) sendLoad() {
	loop.loadCallback(loop.load)
}

func ExitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
