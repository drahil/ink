package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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
	innerWidth, innerHeight := PaneInnerSize(width, height, focused)
	renderedLines := make([]string, 0, innerHeight)

	if innerHeight >= 1 {
		renderedLines = append(renderedLines, truncateCells(PaneTitle("files", focused), innerWidth))
	}
	if innerHeight >= 2 {
		renderedLines = append(renderedLines, "")
	}

	if innerHeight >= 3 {
		searchPrefix := "search: "
		searchWidth := max(0, innerWidth-lipgloss.Width(searchPrefix))
		searchLine := searchPrefix + f.Cursor.RenderLineWithin(f.SearchQuery, searchWidth, focused && cursorVisible)
		renderedLines = append(renderedLines, truncateCells(searchLine, innerWidth))
	}
	if innerHeight >= 4 {
		renderedLines = append(renderedLines, "")
	}

	itemWidth := max(0, innerWidth-2)
	items := f.visibleItems()
	visibleHeight := max(0, innerHeight-len(renderedLines))
	start, end := f.visibleItemRange(items, visibleHeight)
	for index := start; index < end; index++ {
		item := items[index]
		prefix := "  "
		if index == f.Selected {
			prefix = "> "
		}

		renderedLines = append(renderedLines, truncateCells(prefix+truncateCells(item, itemWidth), innerWidth))
	}

	return RenderPane(width, height, strings.Join(renderedLines, "\n"), focused)
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

func (f FilesPane) visibleItemRange(items []string, visibleHeight int) (int, int) {
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
