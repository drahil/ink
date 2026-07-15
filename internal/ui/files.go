package ui

import "strings"

type FilesPane struct {
	Items       []string
	SearchQuery string
	Cursor      CursorPosition
	Selected    int
}

func NewFilesPane(items []string) FilesPane {
	if len(items) == 0 {
		items = []string{"go.mod", "cmd/", "internal/"}
	}

	return FilesPane{
		Items: items,
	}
}

func (f FilesPane) View(width, height int, focused, cursorVisible bool) string {
	var b strings.Builder

	b.WriteString(PaneTitle("files", focused))
	b.WriteString("\n\n")

	b.WriteString("search: ")
	b.WriteString(f.Cursor.RenderLine(f.SearchQuery, focused && cursorVisible))
	b.WriteString("\n\n")

	itemWidth := width - 8
	items := f.visibleItems()
	start, end := f.visibleItemRange(items, height)
	for index := start; index < end; index++ {
		item := items[index]
		if index == f.Selected {
			b.WriteString("> ")
		} else {
			b.WriteString("  ")
		}

		b.WriteString(truncateRunes(item, itemWidth))
		b.WriteString("\n")
	}
	return RenderPane(width, height, b.String(), focused)
}

func (f *FilesPane) Query(r rune) {
	f.SearchQuery = f.SearchQuery + string(r)
	f.MoveCursorRight()
	f.clampSelection()
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
	f.clampSelection()
}

func (f *FilesPane) ClearSearch() {
	f.SearchQuery = ""
	f.MoveCursorToBeginningOfLine()
	f.clampSelection()
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

func (f FilesPane) SelectedItem() string {
	items := f.visibleItems()
	if len(items) == 0 {
		return ""
	}

	return items[f.Selected]
}

func (f FilesPane) visibleItemRange(items []string, height int) (int, int) {
	visibleHeight := height - 6
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	start := f.Selected - visibleHeight/2
	if start < 0 {
		start = 0
	}

	end := start + visibleHeight
	if end > len(items) {
		end = len(items)
		start = end - visibleHeight
		if start < 0 {
			start = 0
		}
	}

	return start, end
}

func (f *FilesPane) MoveSelectionUp() {
	if f.Selected > 0 {
		f.Selected--
	}
}

func (f *FilesPane) MoveSelectionDown() {
	lines := f.visibleItems()
	if f.Selected < len(lines)-1 {
		f.Selected++
	}
}

func (f *FilesPane) clampSelection() {
	items := f.visibleItems()
	if len(items) == 0 {
		f.Selected = 0
		return
	}

	if f.Selected >= len(items) {
		f.Selected = len(items) - 1
	}
}
