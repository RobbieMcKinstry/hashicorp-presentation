package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// A Machine represents a node within a cluster.
// Each machine renders on screen as a group of UI elements.
// Machines can have at most one server.
// And two batch jobs.
type Machine struct {
	ui.Grid
	batch widgets.Gauge
}

func NewMachine(name string) *Machine {
	var machine = &Machine{
		Grid:  *ui.NewGrid(),
		batch: *widgets.NewGauge(),
	}
	// machine.Grid.Border = true
	//machine.Grid.BorderLeft = true
	//machine.Grid.BorderRight = true
	//machine.Grid.BorderTop = true
	//machine.Grid.BorderBottom = true

	machine.batch.Percent = 50
	machine.batch.Title = name

	machine.Grid.Set(
		ui.NewRow(
			1.0,
			ui.NewCol(1.0, &machine.batch),
			// ui.NewCol(),
			// ui.NewCol(),
		),
	)
	return machine
}
