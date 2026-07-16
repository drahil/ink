package ui

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type EditorPane struct {
	Content string
	Cursor  CursorPosition
	Path    string
}

type lineHighlighter func(string) string

var cursorStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("15")).
	Foreground(lipgloss.Color("0"))

func NewEditorPane() EditorPane {
	return EditorPane{
		Content: "<?php\n\necho 'hello';\n",
		Cursor: CursorPosition{
			Row:    0,
			Column: 0,
		},
		Path: "",
	}
}

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

func (e *EditorPane) MoveCursorUp() {
	if e.Cursor.Row > 0 {
		e.Cursor.Row--
	}

	e.Cursor.ClampCursorColumn(e.currentLineLength())
}

func (e *EditorPane) MoveCursorDown() {
	lines := e.lines()
	if e.Cursor.Row < len(lines)-1 {
		e.Cursor.Row++
	}

	e.Cursor.ClampCursorColumn(e.currentLineLength())
}

func (e *EditorPane) MoveCursorLeft() {
	if e.Cursor.Column > 0 {
		e.Cursor.Column--
	}
}

func (e *EditorPane) MoveCursorRight() {
	if e.Cursor.Column < e.currentLineLength() {
		e.Cursor.Column++
	}
}

func (e *EditorPane) MoveCursorToNextLine() {
	e.Cursor.Row++
	e.Cursor.Column = 0
}

func (e *EditorPane) MoveCursorToPreviousLine() {
	e.Cursor.Row--
	e.Cursor.Column = e.currentLineLength()
}

func (e EditorPane) lines() []string {
	lines := strings.Split(e.Content, "\n")
	if len(lines) == 0 {
		return []string{""}
	}

	return lines
}

func (e EditorPane) currentLineLength() int {
	lines := e.lines()
	if e.Cursor.Row < 0 || e.Cursor.Row >= len(lines) {
		return 0
	}

	return len([]rune(lines[e.Cursor.Row]))
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

	rendered := highlight(before) + cursorStyle.Render(cursor) + highlight(after)
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

func (e EditorPane) lineHighlighter() lineHighlighter {
	lexer := lexers.Match(e.Path)
	if lexer == nil {
		lexer = lexers.Analyse(e.Content)
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
	lines := e.lines()
	currentLine := []rune(lines[e.Cursor.Row])

	before := currentLine[:e.Cursor.Column]
	after := currentLine[e.Cursor.Column:]
	nextLine := make([]rune, 0, len(currentLine)+1)
	nextLine = append(nextLine, before...)
	nextLine = append(nextLine, r)
	nextLine = append(nextLine, after...)
	lines[e.Cursor.Row] = string(nextLine)
	e.Content = strings.Join(lines, "\n")
	e.MoveCursorRight()
}

func (e *EditorPane) Backspace() {
	lines := e.lines()

	if e.Cursor.Column == 0 {
		if e.Cursor.Row == 0 {
			return
		}

		previousLine := lines[e.Cursor.Row-1]
		currentLine := lines[e.Cursor.Row]
		e.MoveCursorToPreviousLine()

		nextLines := make([]string, 0, len(lines)-1)
		nextLines = append(nextLines, lines[:e.Cursor.Row]...)
		nextLines = append(nextLines, previousLine+currentLine)
		nextLines = append(nextLines, lines[e.Cursor.Row+2:]...)

		e.Content = strings.Join(nextLines, "\n")
		return
	}

	currentLine := []rune(lines[e.Cursor.Row])
	before := currentLine[:e.Cursor.Column-1]
	after := currentLine[e.Cursor.Column:]
	nextLine := make([]rune, 0, len(currentLine)-1)
	nextLine = append(nextLine, before...)
	nextLine = append(nextLine, after...)
	lines[e.Cursor.Row] = string(nextLine)
	e.Content = strings.Join(lines, "\n")
	e.MoveCursorLeft()
}

func (e *EditorPane) InsertNewline() {
	lines := e.lines()
	currentLine := []rune(lines[e.Cursor.Row])

	before := currentLine[:e.Cursor.Column]
	after := currentLine[e.Cursor.Column:]

	nextLines := make([]string, 0, len(lines)+1)
	nextLines = append(nextLines, lines[:e.Cursor.Row]...)
	nextLines = append(nextLines, string(before), string(after))
	nextLines = append(nextLines, lines[e.Cursor.Row+1:]...)

	lines = nextLines
	e.Content = strings.Join(lines, "\n")
	e.MoveCursorToNextLine()
}

func (e *EditorPane) OpenContent(path, content string) {
	e.Content = content
	e.Path = path
	e.Cursor = CursorPosition{
		Row:    0,
		Column: 0,
	}
}
