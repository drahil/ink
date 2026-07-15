package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("62")).
			Padding(0, 1)

	paneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	focusedPaneStyle = paneStyle.
				BorderForeground(lipgloss.Color("39"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)
)

type CursorPosition struct {
	Row    int
	Column int
}

func RenderHeader(width int) string {
	text := fmt.Sprintf("ink == nulla dies sine linea")

	return headerStyle.
		Width(max(0, width-2)).
		Render(text)
}

func RenderStatusBar(width int, status string) string {
	text := fmt.Sprintf("status: %s | keys: tab focus | alt+1 files | q quit | ctrl+c quit", status)
	text = truncateRunes(text, max(0, width-2))

	return statusStyle.
		Width(max(0, width-2)).
		Render(text)
}

func PaneTitle(name string, active bool) string {
	if active {
		return fmt.Sprintf("[ %s * ]", name)
	}

	return fmt.Sprintf("[ %s ]", name)
}

func RenderPane(width, height int, content string, focused bool) string {
	style := paneStyle
	if focused {
		style = focusedPaneStyle
	}

	return style.
		Width(max(0, width-4)).
		Height(max(0, height-2)).
		Render(content)
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func truncateRunes(value string, maxLength int) string {
	runes := []rune(value)
	if len(runes) <= maxLength {
		return value
	}

	if maxLength <= 0 {
		return ""
	}

	if maxLength <= 3 {
		return string(runes[:maxLength])
	}

	return string(runes[:maxLength-3]) + "..."
}

func (c CursorPosition) RenderLine(line string, cursorVisible bool) string {
	if !cursorVisible {
		return line
	}

	runes := []rune(line)
	column := c.Column

	if column < 0 {
		column = 0
	}

	if column >= len(runes) {
		return string(runes) + cursorStyle.Render(" ")
	}

	before := string(runes[:column])
	cursor := cursorStyle.Render(string(runes[column]))
	after := string(runes[column+1:])

	return before + cursor + after
}

func (c *CursorPosition) ClampCursorColumn(lineLength int) {
	if c.Column > lineLength {
		c.Column = lineLength
	}
}
