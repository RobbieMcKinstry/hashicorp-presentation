package main

import (
	"github.com/RobbieMcKinstry/hashicorp-presentation/client/cluster"
	"github.com/RobbieMcKinstry/hashicorp-presentation/client/display"
	"github.com/RobbieMcKinstry/hashicorp-presentation/client/events"
	ui "github.com/gizak/termui/v3"
	"log"
)

func IgnoreResponse(cluster cluster.Cluster) func(uint64, uint64, uint64) {
	return func(throughput, soft, hard uint64) {
		cluster.NewService(throughput, soft, hard)
	}
}

func main() {
	// Let's start the UI and add a textbox.
	var display = display.NewDisplay()
	var lb = NewLoadBalancer(display)
	defer ui.Close()
	var shutdown = display.Shutdown()
	var eventLoop, onEnter = events.NewEventLoop() // onEnter is how events stream into the loop
	display.SetEventCallback(onEnter)              // What to call when the UI receives an event.
	eventLoop.SetLoadCallback(lb.SetLoad)
	eventLoop.SetMachineCallback(display.NewMachine)
	eventLoop.SetServiceCallback(lb.OnNewService)

	<-shutdown
}

func ExitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
