package main

import (
	ui "github.com/gizak/termui/v3"
)

type Simulation interface {
	AddNeighbor()
}

func main() {
	// Let's start the UI and add a textbox.
	if err := ui.Init(); err != nil {
		ExitOnError(err)
	}
	defer ui.Close()

	var shutdown = addTextbox()

	// TODO
	// º Add a textbox to track the current load.
	// º Add a paragraph element to display the current load.

	go addMachines()

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

func addLoad() {

}

func addTextbox() <-chan struct{} {
	// Add a textbox.
	var width, _ = ui.TerminalDimensions()
	var header = NewHeader(0, 0, width, 3)
	ui.Render(header) // Render again, in case a key has yet to be pressed.
	return header.Shutdown()
}

func addMachines() {
	/*
		var machine = NewMachine("Machine 1")
		machine.Title = "Machine 1"
		machine.SetRect(0, 4, 50, 9)
		ui.Render(machine)
	*/
}

var nodeTmpl = `Addr: %v
ID: %v
Name: %v
Status: %v`
