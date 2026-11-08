package ui

import (
	"fmt"
	"strings"
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

func RenderPrompt(width, maxHeight int, title string, items []string, message string) string {
	frameWidth := promptStyle.GetHorizontalFrameSize()
	if width <= 0 {
		return ""
	}
	if width < frameWidth {
		return blankBlock(width, 1)
	}

	textWidth := max(0, width-frameWidth)
	lines := promptLines(maxHeight, title, items, message)
	for i, line := range lines {
		lines[i] = truncateCellsWithTail(line, textWidth, "...")
	}

	return promptStyle.
		Width(textWidth).
		Render(strings.Join(lines, "\n"))
}

func promptLines(maxHeight int, title string, items []string, message string) []string {
	itemSpace := maxHeight - 2
	if len(items) == 0 || itemSpace < 1 {
		return []string{message}
	}

	shown := items
	hidden := 0
	if len(items) > itemSpace {
		shown = items[:itemSpace-1]
		hidden = len(items) - len(shown)
	}

	lines := make([]string, 0, len(shown)+3)
	lines = append(lines, title)
	for _, item := range shown {
		lines = append(lines, "  "+item)
	}
	if hidden > 0 {
		lines = append(lines, fmt.Sprintf("  ...and %d more", hidden))
	}

	return append(lines, message)
}

func RenderStatusBar(width int, numberOfChangedFiles int) string {
	text := ""
	if numberOfChangedFiles > 0 {
		text = fmt.Sprintf("%d unsaved", numberOfChangedFiles)
	}
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
