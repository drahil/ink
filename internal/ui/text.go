package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func truncateCells(value string, maxWidth int) string {
	return truncateCellsWithTail(value, maxWidth, "")
}

func truncateCellsWithTail(value string, maxWidth int, tail string) string {
	if maxWidth <= 0 {
		return ""
	}

	return ansi.Truncate(value, maxWidth, tail)
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
