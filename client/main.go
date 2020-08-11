package main

import (
	_ "fmt"
	"github.com/RobbieMcKinstry/hashicorp-presentation/client/display"
	"github.com/RobbieMcKinstry/hashicorp-presentation/client/events"
	ui "github.com/gizak/termui/v3"
	_ "github.com/gizak/termui/v3/widgets"
	"log"
)

func main() {
	// Let's start the UI and add a textbox.
	var display = display.NewDisplay()
	defer ui.Close()
	var shutdown = display.Shutdown()

	// var cluster SimVisualization = NewSimulatedCluster()
	// var newMachine = cluster.NewMachineCallback()

	// TODO Set the NewMachineCallback coming from cluster to
	// intercept values from the event loop.

	// var loadTextCallback = addLoadText()
	var eventLoop, onEnter = events.NewEventLoop() // onEnter is how events stream into the loop
	display.SetEventCallback(onEnter)              // What to call when the UI receives an event.
	eventLoop.SetLoadCallback(display.SetLoad)
	eventLoop.SetMachineCallback(display.NewMachine)
	// eventLoop.SetMachineCallback(newMachine)

	// var shutdown = addTextbox(eventWriter)

	// go addMachines()

	// First, we create a list of machines.
	// Each machine has at most one service.

	// Map[machine name] -> Machine

	// Next, create a list of servers.

	// Next, start a timer. Every second, we're going to ping each
	// server and get it's result.
	// We write that result to the widget responsible for this machine.

	// Create a variable to track the current load.

	<-shutdown
}

/*
func addTextbox(callback func(string)) <-chan struct{} {
	// Add a textbox.

	area := NewTextArea()
	area.OnEnter(callback)
	area.Title = "Terminal"
	var width, _ = ui.TerminalDimensions()
	area.SetRect(0, 0, width/2, 3)
	ui.Render(area) // Render again, in case a key has yet to be pressed.
	return area.Shutdown()
}

func addLoadText() func(uint64) {
	var textbox = widgets.NewParagraph()
	textbox.Title = "Current Load"
	textbox.Text = "0 reqs/s"
	var width, _ = ui.TerminalDimensions()
	textbox.SetRect(1+width/2, 0, width, 3)
	ui.Render(textbox)
	return func(load uint64) {
		textbox.Text = fmt.Sprintf("%d reqs/s", load)
		ui.Render(textbox)
	}
}

func addMachines() {
	var width, height = ui.TerminalDimensions()
	var startHeight = 4
	var endHeight = 4 + 3*height/10

	var startHeight2 = endHeight + 1
	var endHeight2 = startHeight2 + 3*height/10

	var startHeight3 = endHeight2 + 1
	var endHeight3 = startHeight3 + 3*height/10

	var machine1 = NewMachine("Machine 1", 0, startHeight, width, endHeight)
	var machine2 = NewMachine("Machine 1", 0, startHeight2, width, endHeight2)
	var machine3 = NewMachine("Machine 1", 0, startHeight3, width, endHeight3)
	ui.Render(machine1)
	ui.Render(machine2)
	ui.Render(machine3)
}
*/
func ExitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
