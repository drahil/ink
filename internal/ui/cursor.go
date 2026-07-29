package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type CursorPosition struct {
	Row    int
	Column int
}

var cursorStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("15")).
	Foreground(lipgloss.Color("0"))

func (c CursorPosition) RenderLineWithin(line string, width int, cursorVisible bool) string {
	if width <= 0 {
		return ""
	}

	if !cursorVisible {
		return truncateCells(line, width)
	}

	column := c.Column
	if column < 0 {
		column = 0
	}

	start := 0
	if column >= width {
		start = column - width + 1
	}

	segment := ansi.Cut(line, start, start+width)
	cursorColumn := column - start
	before := ansi.Cut(segment, 0, cursorColumn)
	cursorText := ansi.Cut(segment, cursorColumn, cursorColumn+1)
	afterStart := cursorColumn + lipgloss.Width(cursorText)
	after := ansi.Cut(segment, afterStart, width)

	if cursorText == "" {
		cursorText = " "
	}

	return truncateCells(before+cursorStyle.Render(cursorText)+after, width)
}

func (c *CursorPosition) ClampCursorColumn(lineLength int) {
	if c.Column > lineLength {
		c.Column = lineLength
	}
}
