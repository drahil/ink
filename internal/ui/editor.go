package ui

import (
	"strings"
)

type EditorPane struct {
	Content string
}

func NewEditorPane() EditorPane {
	return EditorPane{}
}

func (e EditorPane) View(width, height int, focused bool) string {
	var b strings.Builder

	b.WriteString(PaneTitle("editor", focused))
	b.WriteString("\n\n")
	b.WriteString("editor placeholder\n")
	b.WriteString("later: open file text goes here\n\n")
	b.WriteString("tab changes focus")

	return RenderPane(width, height, b.String(), focused)
}
