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
	*ui.Grid
	chart *widgets.Sparkline
	batch *widgets.Gauge
	cpu   *widgets.PieChart
}

func NewMachine(name string, x1, y1, x2, y2 int) *Machine {
	var machine = &Machine{
		Grid:  ui.NewGrid(),
		chart: widgets.NewSparkline(),
		batch: widgets.NewGauge(),
		cpu:   widgets.NewPieChart(),
	}
	/////// GRID CONFIG
	machine.Grid.SetRect(x1, y1, x2, y2)
	machine.Grid.Title = name
	machine.Grid.Border = true
	machine.Grid.BorderLeft = true
	machine.Grid.BorderRight = true
	machine.Grid.BorderTop = true
	machine.Grid.BorderBottom = true
	/////// PIE CONFIG
	machine.cpu.Title = "CPU"
	machine.cpu.Data = []float64{1.0}
	machine.cpu.LabelFormatter = func(i int, v float64) string {
		var label = "Free"
		if i == 1 {
			label = "Server"
		} else if i == 2 {
			label = "Batch"
		}
		return label
	}
	/////// GAUGE CONFIG
	machine.batch.Percent = 50
	machine.batch.Title = "Batch"
	/////// SPARKLINE CONFIG
	machine.chart.Data = []float64{5.0, 4.0, 3.5, 2.5}
	machine.chart.LineColor = ui.ColorGreen
	machine.chart.TitleStyle.Fg = ui.ColorBlue

	var chartGroup = widgets.NewSparklineGroup(machine.chart)
	chartGroup.Title = "Throughput"

	var paragraph = widgets.NewParagraph()
	paragraph.Text = "Hello Goodbye"

	machine.Grid.Set(
		ui.NewRow(1.0,
			ui.NewCol(1.0/3, chartGroup),
			ui.NewCol(1.0/3, machine.batch),
			ui.NewCol(1.0/3, machine.cpu),
		),
	)
	return machine
}
