package display

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"time"
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

	chartStream chan float64
}

func NewMachine(name string, x1, y1, x2, y2 int) *Machine {
	var machine = &Machine{
		Grid:        ui.NewGrid(),
		chart:       widgets.NewSparkline(),
		batch:       widgets.NewGauge(),
		cpu:         widgets.NewPieChart(),
		chartStream: make(chan float64, 100),
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
	machine.batch.Percent = 0
	machine.batch.Title = "Batch"
	/////// SPARKLINE CONFIG
	// machine.chart.Data = []float64{5.0, 4.0, 3.5, 2.5}
	machine.chart.Data = []float64{}
	machine.chart.LineColor = ui.ColorGreen
	machine.chart.TitleStyle.Fg = ui.ColorBlue

	go machine.watchSparkline()

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

func (machine *Machine) watchSparkline() {
	var ticker = time.Tick(1 * time.Second)
	for {
		select {
		case <-ticker:
			ui.Render(machine.Grid)
		case val := <-machine.chartStream:
			machine.chart.Data = append(machine.chart.Data, val)
			if len(machine.chart.Data) > 10 {
				machine.chart.Data = machine.chart.Data[1:]
			}
		}
	}
}

func (machine *Machine) GetSparklineStream() chan float64 {
	return machine.chartStream
}
