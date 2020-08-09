package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"fmt"
	"log"
	"os"
	"unicode/utf8"
)

// TextArea allows editable text to be rendered to the screen.
type TextArea struct {
	buffer []rune
	widgets.Paragraph
	events   <-chan ui.Event
	shutdown chan struct{}
}

func NewTextArea() *TextArea {
	var area = &TextArea{
		buffer:    []rune{},
		Paragraph: *widgets.NewParagraph(),
		events:    ui.PollEvents(),
		shutdown:  make(chan struct{}),
	}
	ui.Render(area)
	go area.watchEvents()

	return area
}

func (area *TextArea) Shutdown() <-chan struct{} {
	return area.shutdown
}

func (area *TextArea) watchEvents() {
	for e := range area.events {
		if e.Type != ui.KeyboardEvent {
			continue
		}

		var event = e.ID
		if isRenderable(event) {
			// Append it to the buffer.
			// Unless its a space...
			var character = event
			if isSpace(event) {
				character = " "
			}
			area.buffer = append(area.buffer, []rune(character)...)

		} else if isBackspace(event) {
			// Remove the last item from the buffer.
			if len(area.buffer) < 1 {
				continue
			}
			area.buffer = area.buffer[:len(area.buffer)-1]
		} else if isEscape(event) {
			fmt.Println("Exit signal received.")
			os.Exit(0)
		}

		area.Paragraph.Text = string(area.buffer)
		ui.Render(area)
	}
}

var _ ui.Drawable = &TextArea{}

func ExitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func isAlphabetic(event string) bool {
	var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var first, _ = utf8.DecodeRuneInString(event)
	for _, letter := range alphabet {
		if first == letter {
			return true
		}
	}
	return false
}

func isNumeric(event string) bool {
	var numbers = []rune("0123456789")
	var first, _ = utf8.DecodeRuneInString(event)
	for _, num := range numbers {
		if first == num {
			return true
		}
	}
	return false
}

func isAlphanumeric(event string) bool {
	return isAlphabetic(event) || isNumeric(event)
}

func isSpace(event string) bool {
	return event == "<Space>"
}

func isBackspace(event string) bool {
	return event == "<Backspace>"
}

func isEscape(event string) bool {
	return event == "<Escape>"
}

func isEnter(event string) bool {
	return event == "<Enter>"
}

func isRenderable(event string) bool {
	return isAlphabetic(event) || isSpace(event)
}
