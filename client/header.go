package main

import (
	"fmt"
	"image"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// A Header is a textarea for terminal input followed by a Status box for the current
// system load.
type Header struct {
	area     TextArea
	load     widgets.Paragraph
	shutdown <-chan struct{}

	loadVal uint64
}

func NewHeader(x1, y1, x2, y2 int) *Header {
	var header = &Header{
		load: *widgets.NewParagraph(),
	}
	header.area = *NewTextArea(func() {
		ui.Render(header)
	})
	header.shutdown = header.area.Shutdown()
	header.area.Title = "Terminal"
	header.load.Title = "Current Load"

	header.SetRect(x1, y1, x2, y2)
	header.SetLoad(0)

	return header
}

func (header *Header) SetLoad(load uint64) {
	header.loadVal = load
	header.load.Text = fmt.Sprintf("%d reqs/s", load)
	ui.Render(header)
}

func (header *Header) Draw(buffer *ui.Buffer) {
	header.area.Draw(buffer)
	header.load.Draw(buffer)
}

func (header *Header) GetRect() image.Rectangle {
	var r1 = header.area.GetRect()
	var r2 = header.load.GetRect()
	return image.Rectangle{
		Min: r1.Min,
		Max: r2.Max,
	}
}

func (header *Header) Lock() {
	header.area.Lock()
}

func (header *Header) Unlock() {
	header.area.Unlock()
}

func (header *Header) SetRect(x1, y1, x2, y2 int) {
	header.area.SetRect(x1, y1, x2/2, y2)
	header.load.SetRect(1+x2/2, y1, x2, y2)
}

func (header *Header) Shutdown() <-chan struct{} {
	return header.shutdown
}
