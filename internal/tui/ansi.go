package tui

import "fmt"

const ANSI_CLEAR_MODIFIER = "\x1b[m"
const ANSI_CLEAR_SCREEN = "\x1b[H\x1b[J"

const ANSI_MOVE_TO = "\x1b[%d;%dH"
const ANSI_MOVE_UP = "\x1b[%dA"
const ANSI_MOVE_DOWN = "\x1b[%dB"
const ANSI_MOVE_RIGHT = "\x1b[%dC"
const ANSI_MOVE_LEFT = "\x1b[%dD"

func AnsiMoveTo(x, y int) string {
	return fmt.Sprintf(ANSI_MOVE_TO, y, x)
}

func AnsiMoveUp(i int) string {
	return fmt.Sprintf(ANSI_MOVE_UP, i)
}

func AnsiMoveDown(i int) string {
	return fmt.Sprintf(ANSI_MOVE_DOWN, i)
}

func AnsiMoveLeft(i int) string {
	return fmt.Sprintf(ANSI_MOVE_LEFT, i)
}

func AnsiMoveRight(i int) string {
	return fmt.Sprintf(ANSI_MOVE_RIGHT, i)
}

const ANSI_BOLD = "\x1b[1m"
