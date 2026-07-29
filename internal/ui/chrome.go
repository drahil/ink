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
