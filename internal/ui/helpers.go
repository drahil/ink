package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
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
	const text = "ink == nulla dies sine linea"
	frameWidth := headerStyle.GetHorizontalFrameSize()
	if width <= 0 {
		return ""
	}
	if width < frameWidth {
		return blankBlock(width, 1)
	}

	return headerStyle.
		Width(max(0, width-frameWidth)).
		Render(truncateCells(text, max(0, width-frameWidth)))
}

func RenderStatusBar(width int, status string) string {
	text := fmt.Sprintf("status: %s | keys: tab focus | alt+1 files | q quit | ctrl+c quit", status)
	frameWidth := statusStyle.GetHorizontalFrameSize()
	if width <= 0 {
		return ""
	}
	if width < frameWidth {
		return blankBlock(width, 1)
	}

	text = truncateCells(text, max(0, width-frameWidth))

	return statusStyle.
		Width(max(0, width-frameWidth)).
		Render(text)
}

func PaneTitle(name string, active bool) string {
	if active {
		return fmt.Sprintf("[ %s * ]", name)
	}

	return fmt.Sprintf("[ %s ]", name)
}

func RenderPane(width, height int, content string, focused bool) string {
	style := paneRenderStyle(focused)
	frameWidth, frameHeight := style.GetFrameSize()
	if width <= 0 || height <= 0 {
		return ""
	}
	if width < frameWidth || height < frameHeight {
		return blankBlock(width, height)
	}

	innerWidth, innerHeight := max(0, width-frameWidth), max(0, height-frameHeight)
	content = fitBlock(content, innerWidth, innerHeight)

	return style.Render(content)
}

func PaneInnerSize(width, height int, focused bool) (int, int) {
	return paneInnerSize(paneRenderStyle(focused), width, height)
}

func paneRenderStyle(focused bool) lipgloss.Style {
	if focused {
		return focusedPaneStyle
	}

	return paneStyle
}

func paneInnerSize(style lipgloss.Style, width, height int) (int, int) {
	frameWidth, frameHeight := style.GetFrameSize()
	return max(0, width-frameWidth), max(0, height-frameHeight)
}

func blankBlock(width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	lines := make([]string, height)
	for i := range lines {
		lines[i] = strings.Repeat(" ", width)
	}

	return strings.Join(lines, "\n")
}

func fitBlock(content string, width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	sourceLines := strings.Split(content, "\n")
	lines := make([]string, height)
	for i := range lines {
		if i >= len(sourceLines) {
			lines[i] = strings.Repeat(" ", width)
			continue
		}

		line := truncateCells(sourceLines[i], width)
		lines[i] = line + strings.Repeat(" ", max(0, width-lipgloss.Width(line)))
	}

	return strings.Join(lines, "\n")
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func truncateCells(value string, maxWidth int) string {
	return truncateCellsWithTail(value, maxWidth, "")
}

func truncateCellsWithTail(value string, maxWidth int, tail string) string {
	if maxWidth <= 0 {
		return ""
	}

	return ansi.Truncate(value, maxWidth, tail)
}

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
