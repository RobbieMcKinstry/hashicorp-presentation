package display

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Display struct {
	// Terminal
	terminal *TextArea
	// Load Textarea
	loadText *widgets.Paragraph
	// Machines
	machines [3]*Machine
}

func NewDisplay() *Display {
	var err = ui.Init()
	ExitOnError(err)
	var width, _ = ui.TerminalDimensions()

	var display = &Display{
		terminal: NewTextArea(),
		loadText: widgets.NewParagraph(),
		machines: [3]*Machine{nil, nil, nil},
	}

	defer ui.Render(display.terminal)
	defer ui.Render(display.loadText)

	display.terminal.Title = "Terminal"
	display.terminal.SetRect(0, 0, width/2, 3)

	display.loadText.Title = "Current Load"
	display.loadText.Text = "0 reqs/s"
	display.loadText.SetRect(1+width/2, 0, width, 3)

	return display
}

func (display *Display) countMachines() int {
	for i := 0; i < len(display.machines); i++ {
		if display.machines[i] == nil {
			return i
		}
	}
	return 3
}

func (display *Display) addMachine() {
	// First, let's see how many machines we currently have.
	var machineCount = display.countMachines()
	if machineCount >= 3 {
		panic("Cannot add another machine.")
	}
	display.addMachineAtIndex(machineCount)
}

func (display *Display) addMachineAtIndex(index int) {
	var width, height = ui.TerminalDimensions()
	var startHeight = 4 + 3*index*height/10
	var endHeight = 4 + (index+1)*3*height/10
	var id = fmt.Sprintf("Machine %v", index)
	var machine = NewMachine(id, 0, startHeight, width, endHeight)
	display.machines[index] = machine
	ui.Render(machine)
}

func (display *Display) Shutdown() <-chan struct{} {
	return display.terminal.Shutdown()
}