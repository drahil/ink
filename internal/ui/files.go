package ui

import "strings"

type FilesPane struct {
	Items []string
}

func NewFilesPane() FilesPane {
	return FilesPane{
		Items: []string{
			"project root",
			"cmd/",
			"internal/",
			"go.mod",
		},
	}
}

func (f FilesPane) View(width, height int, focused bool) string {
	var b strings.Builder

	b.WriteString(PaneTitle("files", focused))
	b.WriteString("\n\n")
	b.WriteString(strings.Join(f.Items, "\n"))

	return RenderPane(width, height, b.String(), focused)
}
