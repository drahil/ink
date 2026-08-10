package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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
	items := f.filteredItems()
	visibleHeight := max(0, innerHeight-len(renderedLines))
	start, end := f.visibleItemRange(items, visibleHeight)
	for index := start; index < end; index++ {
		item := items[index]
		prefix := "  "
		if index == f.Selected {
			prefix = "> "
		}

		marker := "  "
		name := item.Name
		if item.IsDir {
			if f.dirExpanded(item.Path) {
				marker = "- "
			} else {
				marker = "+ "
			}
			name += "/"
		}

		label := strings.Repeat("  ", item.Depth) + marker + name
		renderedLines = append(renderedLines, truncateCells(prefix+truncateCells(label, itemWidth), innerWidth))
	}

	return RenderPane(width, height, strings.Join(renderedLines, "\n"), focused)
}
