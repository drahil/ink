package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func (e EditorPane) View(width, height int, focused bool, cursorVisible bool) string {
	lines := e.lines()
	innerWidth, innerHeight := PaneInnerSize(width, height, focused)
	gutterDigits := len(fmt.Sprintf("%d", len(lines)))
	gutterWidth := gutterDigits + 2
	contentWidth := max(0, innerWidth-gutterWidth)
	highlight := e.lineHighlighter()
	renderedLines := make([]string, 0, innerHeight)

	if innerHeight >= 1 {
		renderedLines = append(renderedLines, e.title(focused, innerWidth))
	}
	if innerHeight >= 2 {
		renderedLines = append(renderedLines, "")
	}

	visibleHeight := max(0, innerHeight-len(renderedLines))
	start, end := e.visibleLineRange(lines, visibleHeight)
	for row := start; row < end; row++ {
		line := lines[row]
		lineCursorVisible := focused && cursorVisible && row == e.Cursor.Row
		prefix := fmt.Sprintf("%*d  ", gutterDigits, row+1)
		renderedLine := prefix + e.renderVisibleLine(line, contentWidth, lineCursorVisible, highlight)
		renderedLines = append(renderedLines, truncateCells(renderedLine, innerWidth))
	}

	return RenderPane(width, height, strings.Join(renderedLines, "\n"), focused)
}

func (e EditorPane) title(focused bool, width int) string {
	if e.Path == "" {
		return PaneTitle("editor", focused)
	}

	if focused {
		return truncateCellsWithTail(fmt.Sprintf("[ %s * ]", e.Path), width, "...")
	}

	return truncateCellsWithTail(fmt.Sprintf("[ %s ]", e.Path), width, "...")
}

func (e EditorPane) visibleLineRange(lines []string, visibleHeight int) (int, int) {
	if visibleHeight <= 0 {
		return 0, 0
	}
	start := e.Cursor.Row - visibleHeight/2
	if start < 0 {
		start = 0
	}

	end := start + visibleHeight
	if end > len(lines) {
		end = len(lines)
		start = end - visibleHeight
		if start < 0 {
			start = 0
		}
	}

	return start, end
}

func (e EditorPane) renderVisibleLine(line string, width int, cursorVisible bool, highlight lineHighlighter) string {
	if width <= 0 {
		return ""
	}

	line, column := expandTabs(line, e.Cursor.Column)

	if column < 0 {
		column = 0
	}

	start := 0
	if column >= width {
		start = column - width + 1
	}

	segment := ansi.Cut(line, start, start+width)
	highlightedSegment := highlight(segment)

	if !cursorVisible {
		return truncateCells(highlightedSegment, width)
	}

	cursorColumn := column - start
	if cursorColumn < 0 {
		cursorColumn = 0
	}

	before := ansi.Cut(highlightedSegment, 0, cursorColumn)
	cursor := ansi.Cut(segment, cursorColumn, cursorColumn+1)
	afterStart := cursorColumn + lipgloss.Width(cursor)
	after := ansi.Cut(highlightedSegment, afterStart, width)

	if cursor == "" {
		cursor = " "
	}

	rendered := before + cursorStyle.Render(cursor) + after
	return truncateCells(rendered, width)
}

func expandTabs(line string, cursorColumn int) (string, int) {
	var b strings.Builder
	expandedColumn := cursorColumn
	visualColumn := 0
	sourceColumn := 0

	for _, r := range line {
		if r == '\t' {
			spaces := 4 - visualColumn%4
			b.WriteString(strings.Repeat(" ", spaces))
			if sourceColumn < cursorColumn {
				expandedColumn += spaces - 1
			}
			visualColumn += spaces
			sourceColumn++
			continue
		}

		b.WriteRune(r)
		visualColumn++
		sourceColumn++
	}

	return b.String(), expandedColumn
}
