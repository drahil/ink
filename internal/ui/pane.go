package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

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
