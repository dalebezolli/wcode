package main

import (
	"strconv"
	"strings"
)

type Box struct {
	Width  int
	Height int
	Title   string
	Content string
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

	ansiString := topLeft + horizontal + " " + b.Title + " " + strings.Repeat(horizontal, b.Width-countOfMandatoryTopCharsWithoutLeftSpacer-len(b.Title)) + topRight
	ansiString += strings.Repeat(lineShift+vertical+"\x1b["+strconv.Itoa(b.Width-countOfMandatoryContainerChars)+"C"+vertical, b.Height-countOfMandatoryContainerChars)
	ansiString += lineShift + bottomLeft + strings.Repeat(horizontal, b.Width-countOfMandatoryContainerChars) + bottomRight

	ansiString += "\x1b[" + strconv.Itoa(b.Width-contentLeftPadding) + "D" + "\x1b[" + strconv.Itoa(b.Height-contentTopPadding) + "A" + b.Content

	display.Add(ansiString)
}
