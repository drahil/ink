package ui

import "strings"

type CommandPane struct {
	Text string
}

func NewCommandPane() CommandPane {
	return CommandPane{
		Text: "command/status input placeholder",
	}
}

func (c CommandPane) View(width, height int, focused bool) string {
	innerWidth, innerHeight := PaneInnerSize(width, height, focused)
	renderedLines := make([]string, 0, innerHeight)

	if innerHeight >= 1 {
		renderedLines = append(renderedLines, truncateCells(PaneTitle("command", focused), innerWidth))
	}
	if innerHeight >= 2 {
		renderedLines = append(renderedLines, "")
	}
	if innerHeight >= 3 {
		renderedLines = append(renderedLines, truncateCells(c.Text, innerWidth))
	}

	return RenderPane(width, height, strings.Join(renderedLines, "\n"), focused)
}
