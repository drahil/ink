package ui

import "strings"

func (f *FilesPane) InsertSearchRune(r rune) {
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

func (f FilesPane) filteredItems() []FileRow {
	rows := f.visibleRows()
	if f.SearchQuery == "" {
		return rows
	}

	visibleItems := make([]FileRow, 0)
	query := strings.ToLower(f.SearchQuery)
	for _, item := range rows {
		if strings.Contains(strings.ToLower(item.Path), query) {
			visibleItems = append(visibleItems, item)
		}
	}

	return visibleItems
}

func (f FilesPane) visibleRows() []FileRow {
	rows := make([]FileRow, 0, len(f.Items))
	for _, row := range f.Items {
		if f.rowVisible(row) {
			rows = append(rows, row)
		}
	}

	return rows
}

func (f FilesPane) rowVisible(row FileRow) bool {
	parent := parentPath(row.Path)
	for parent != "" {
		if !f.dirExpanded(parent) {
			return false
		}

		parent = parentPath(parent)
	}

	return true
}

func parentPath(filePath string) string {
	index := strings.LastIndex(filePath, "/")
	if index < 0 {
		return ""
	}

	return filePath[:index]
}

func (f FilesPane) dirExpanded(path string) bool {
	if f.Expanded == nil {
		return true
	}

	expanded, ok := f.Expanded[path]
	if !ok {
		return true
	}

	return expanded
}

func (f *FilesPane) ToggleExpanded(row FileRow) {
	if !row.IsDir {
		return
	}
	if f.Expanded == nil {
		f.Expanded = make(map[string]bool)
	}

	f.Expanded[row.Path] = !f.dirExpanded(row.Path)
	f.clampSelection()
}

func (f FilesPane) SelectedRow() (FileRow, bool) {
	items := f.filteredItems()
	if len(items) == 0 {
		return FileRow{}, false
	}

	return items[f.Selected], true
}

func (f FilesPane) visibleItemRange(items []FileRow, visibleHeight int) (int, int) {
	if visibleHeight <= 0 {
		return 0, 0
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
	lines := f.filteredItems()
	if f.Selected < len(lines)-1 {
		f.Selected++
	}
}

func (f *FilesPane) clampSelection() {
	items := f.filteredItems()
	if len(items) == 0 {
		f.Selected = 0
		return
	}

	if f.Selected >= len(items) {
		f.Selected = len(items) - 1
	}
}
