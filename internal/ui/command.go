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
	var b strings.Builder

	b.WriteString(PaneTitle("command", focused))
	b.WriteString("\n\n")
	b.WriteString(c.Text)

	return RenderPane(width, height, b.String(), focused)
}
