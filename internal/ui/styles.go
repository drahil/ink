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

	cursorStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("15")).
			Foreground(lipgloss.Color("0"))
)

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
	text := fmt.Sprintf("status: %s | keys: tab focus | alt+1 files | alt+f search | q quit | ctrl+c quit", status)
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
