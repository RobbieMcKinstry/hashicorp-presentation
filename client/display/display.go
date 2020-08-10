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
	var width, height = ui.TerminalDimensions()

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

	var startHeight = 4
	var endHeight = 4 + 3*height/10
	for i := 0; i < 3; i++ {
		var id = fmt.Sprintf("Machine %v", i)
		var machine = NewMachine(id, 0, startHeight, width, endHeight)
		display.machines[i] = machine
		startHeight = endHeight + 1
		endHeight = startHeight + 3*height/10
		ui.Render(machine)
	}

	return display
}

func (display *Display) Shutdown() <-chan struct{} {
	return display.terminal.Shutdown()
}
