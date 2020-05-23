package mritd

import (
	"strconv"
)

const esc = "\033["

const (
	hideCursor       = esc + "?25l"
	showCursor       = esc + "?25h"
	clearLine        = esc + "2K"
	clearDown        = esc + "J"
	clearStartOfLine = esc + "1K"
	clearScreen      = esc + "2J"
	moveUp           = esc + "1A"
	move2Up          = esc + "2A"
	moveDown         = esc + "1B"
	clearTerminal    = "\033c"
)

func upLine(n uint) string {
	return movementCode(n, 'A')
}

func movementCode(n uint, code rune) string {
	return esc + strconv.FormatUint(uint64(n), 10) + string(code)
}
