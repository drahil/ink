package ui

import (
	"fmt"
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

func RenderPrompt(width int, message string) string {
	frameWidth := promptStyle.GetHorizontalFrameSize()
	if width <= 0 {
		return ""
	}
	if width < frameWidth {
		return blankBlock(width, 1)
	}

	return promptStyle.
		Width(max(0, width-frameWidth)).
		Render(truncateCells(message, max(0, width-frameWidth)))
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
