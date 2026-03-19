package tui

import (
	"strconv"
	"strings"
)

type Box struct {
	Width  int
	Height int
	Title  string
}

const horizontal = "─"
const vertical = "│"
const topLeft = "┌"
const topRight = "┐"
const bottomLeft = "└"
const bottomRight = "┘"

const countOfMandatoryTopChars int = 6
const countOfMandatoryTopCharsWithoutLeftSpacer int = countOfMandatoryTopChars - 1
const countOfMandatoryContainerChars int = 2

const contentLeftPadding int = 1
const contentTopPadding int = 2

func (b *Box) Render(display *TUI) {
	if b.Width < 0 || b.Height < 0 {
		return
	}

	lineShift := "\x1b[1B\x1b[" + strconv.Itoa(b.Width) + "D"

	if b.Width-countOfMandatoryTopChars-len(b.Title) < 0 {
		b.Title = b.Title[:b.Width-countOfMandatoryTopChars]
	}

	display.Add(topLeft)
	display.Add(horizontal)
	display.Add(" ")
	display.Add(b.Title)
	display.Add(" ")
	display.Add(strings.Repeat(
		horizontal,
		b.Width-countOfMandatoryTopCharsWithoutLeftSpacer-len(b.Title),
	))
	display.Add(topRight)

	display.Add(strings.Repeat(
		lineShift+vertical+AnsiMoveRight(b.Width-countOfMandatoryContainerChars)+vertical,
		b.Height-countOfMandatoryContainerChars,
	))

	display.Add(lineShift)
	display.Add(bottomLeft)
	display.Add(strings.Repeat(
		horizontal,
		b.Width-countOfMandatoryContainerChars,
	))
	display.Add(bottomRight)
}
