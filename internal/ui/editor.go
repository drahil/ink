package ui

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/drahil/ink/internal/editor"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type EditorPane struct {
	Buffer editor.Buffer
	Path   string
}

type lineHighlighter func(string) string

func NewEditorPane() EditorPane {
	return EditorPane{
		Buffer: editor.NewBuffer("<?php\n\necho 'hello';\n"),
		Path:   "",
	}
}

func (e EditorPane) View(width, height int, focused bool, cursorVisible bool) string {
	lines := e.Buffer.Lines()
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
		lineCursorVisible := focused && cursorVisible && row == e.Buffer.Cursor.Row
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

func (e *EditorPane) MoveCursorUp() {
	e.Buffer.MoveCursorUp()
}

func (e *EditorPane) MoveCursorDown() {
	e.Buffer.MoveCursorDown()
}

func (e *EditorPane) MoveCursorLeft() {
	e.Buffer.MoveCursorLeft()
}

func (e *EditorPane) MoveCursorRight() {
	e.Buffer.MoveCursorRight()
}

func (e EditorPane) visibleLineRange(lines []string, visibleHeight int) (int, int) {
	if visibleHeight <= 0 {
		return 0, 0
	}

	start := 0
	if e.Buffer.Cursor.Row >= visibleHeight {
		start = e.Buffer.Cursor.Row - visibleHeight + 1
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

	line, column := expandTabs(line, e.Buffer.Cursor.Column)

	if !cursorVisible {
		return truncateCells(highlight(truncateCells(line, width)), width)
	}

	if column < 0 {
		column = 0
	}

	start := 0
	if column >= width {
		start = column - width + 1
	}

	segment := ansi.Cut(line, start, start+width)
	cursorColumn := column - start

	if cursorColumn < 0 {
		cursorColumn = 0
	}

	before := ansi.Cut(segment, 0, cursorColumn)
	cursor := ansi.Cut(segment, cursorColumn, cursorColumn+1)
	afterStart := cursorColumn + lipgloss.Width(cursor)
	after := ansi.Cut(segment, afterStart, width)

	if cursor == "" {
		cursor = " "
	}

	rendered := highlightSegment(highlight, before) + cursorStyle.Render(cursor) + highlightSegment(highlight, after)
	return truncateCells(rendered, width)
}

func highlightSegment(highlight lineHighlighter, segment string) string {
	rendered := highlight(segment)
	missingWidth := lipgloss.Width(segment) - lipgloss.Width(rendered)
	if missingWidth <= 0 {
		return rendered
	}

	return rendered + strings.Repeat(" ", missingWidth)
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

func (e EditorPane) lineHighlighter() lineHighlighter {
	lexer := lexers.Match(e.Path)
	if lexer == nil {
		lexer = lexers.Analyse(e.Buffer.Content)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	formatter := formatters.Get("terminal16m")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	return func(line string) string {
		if line == "" {
			return ""
		}

		iterator, err := lexer.Tokenise(nil, line)
		if err != nil {
			return line
		}

		var b bytes.Buffer
		if err := formatter.Format(&b, style, iterator); err != nil {
			return line
		}

		return strings.TrimRight(b.String(), "\n")
	}
}

func (e *EditorPane) InsertRune(r rune) {
	e.Buffer.InsertRune(r)
}

func (e *EditorPane) Backspace() {
	e.Buffer.Backspace()
}

func (e *EditorPane) InsertNewline() {
	e.Buffer.InsertNewline()
}

func (e *EditorPane) OpenContent(path, content string) {
	e.Buffer.Open(content)
	e.Path = path
}

func (e *EditorPane) MoveToFirstMatch(query string) bool {
	return e.Buffer.MoveToFirstMatch(query)
}
