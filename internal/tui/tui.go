package tui

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"golang.org/x/term"
)

type Event any

type EventResize struct {
	Width  int
	Height int
}

type EventKeyPress struct {
	ReadBuffer []byte
}

type Model interface {
	View(tui *TUI)
	Update(e Event) bool
}

type TUI struct {
	tty    *os.File
	Width  int
	Height int

	modifier       string
	buffer         string
	event          chan Event
	eventListeners map[reflect.Type]map[string]func(any)

	oldTermState *term.State
	readBuffer   []byte
	value        []byte

	model Model
}

func NewTUI(model Model) *TUI {
	f, err := os.OpenFile("/dev/tty", os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Error:", err)
	}

	oldState, _ := term.MakeRaw(int(f.Fd()))

	tui := &TUI{
		tty:   f,
		event: make(chan Event, 1),

		oldTermState: oldState,
		readBuffer:   make([]byte, 3),
		value:        make([]byte, 0, 80),

		model: model,
	}
	tui.SyncSize()

	go tui.listenForResize()
	go tui.listenForInput()

	return tui
}

func (tui *TUI) listenForResize() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGWINCH)

	for {
		<-c
		tui.SyncSize()
		tui.event <- EventResize{Width: tui.Width, Height: tui.Height}
	}
}

func (tui *TUI) listenForInput() {
	for {
		tui.readBuffer[1] = 0
		tui.readBuffer[2] = 0
		tui.tty.Read(tui.readBuffer)

		copyBuffer := make([]byte, 3)
		copy(copyBuffer, tui.readBuffer)
		tui.event <- EventKeyPress{copyBuffer}
	}
}

func (tui *TUI) Clear() {
	tui.buffer += ANSI_CLEAR_SCREEN
}

func (tui *TUI) Flush() {
	tui.tty.WriteString(tui.buffer)
	tui.buffer = ""
}

func (tui *TUI) MoveAt(x, y int) {
	tui.buffer += fmt.Sprintf(ANSI_MOVE_TO, y, x)
}

func (tui *TUI) Add(data string) {
	tui.buffer += data
}

func (tui *TUI) SyncSize() {
	w, h, _ := term.GetSize(int(tui.tty.Fd()))
	tui.Width = w
	tui.Height = h
}

type Status int

const (
	Status_Ok = iota
	Status_Finished
	Status_Terminated
)

func (tui *TUI) GetValue() string {
	return string(tui.value)
}

func (tui *TUI) Close() {
	close(tui.event)
	term.Restore(int(tui.tty.Fd()), tui.oldTermState)
	tui.tty.Close()
}

func (tui *TUI) Run() {
	for {
		tui.model.View(tui)

		event := <-tui.event
		keepGoing := tui.model.Update(event)

		if !keepGoing {
			return
		}
	}
}
