package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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
