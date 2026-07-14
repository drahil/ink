package ui

import "strings"

type FilesPane struct {
	Items       []string
	SearchQuery string
	Cursor      CursorPosition
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

func (f FilesPane) View(width, height int, focused, cursorVisible bool) string {
	var b strings.Builder

	b.WriteString(PaneTitle("files", focused))
	b.WriteString("\n\n")

	b.WriteString("search: ")
	b.WriteString(f.Cursor.RenderLine(f.SearchQuery, focused && cursorVisible))
	b.WriteString("\n\n")
	b.WriteString(strings.Join(f.visibleItems(), "\n"))

	return RenderPane(width, height, b.String(), focused)
}

func (f *FilesPane) Query(r rune) {
	f.SearchQuery = f.SearchQuery + string(r)
	f.MoveCursorRight()
}

func (f *FilesPane) MoveCursorRight() {
	f.Cursor.Column++
}

func (f *FilesPane) MoveCursorLeft() {
	if f.Cursor.Column > 0 {
		f.Cursor.Column--
	}
}

func (f *FilesPane) MoveCursorToBeginningOfLine() {
	f.Cursor.Column = 0
}

func (f *FilesPane) Backspace() {
	if f.Cursor.Column == 0 {
		return
	}

	currentLine := []rune(f.SearchQuery)
	before := currentLine[:f.Cursor.Column-1]
	after := currentLine[f.Cursor.Column:]
	nextLine := make([]rune, 0, len(currentLine)-1)
	nextLine = append(nextLine, before...)
	nextLine = append(nextLine, after...)
	f.SearchQuery = string(nextLine)
	f.MoveCursorLeft()
}

func (f *FilesPane) ClearSearch() {
	f.SearchQuery = ""
	f.MoveCursorToBeginningOfLine()
}

func (f FilesPane) visibleItems() []string {
	if f.SearchQuery == "" {
		return f.Items
	}

	visibleItems := make([]string, 0)

	for _, item := range f.Items {
		if strings.Contains(strings.ToLower(item), strings.ToLower(f.SearchQuery)) {
			visibleItems = append(visibleItems, item)
		}
	}

	return visibleItems
}
