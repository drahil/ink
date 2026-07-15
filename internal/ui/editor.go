package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type EditorPane struct {
	Content string
	Cursor  CursorPosition
	Path    string
}

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
	var b strings.Builder

	contentWidth := max(1, width-6)

	b.WriteString(e.title(focused, contentWidth))
	b.WriteString("\n\n")
	fmt.Fprintf(&b, "cursor row: %d\n", e.Cursor.Row)
	fmt.Fprintf(&b, "cursor col: %d\n\n", e.Cursor.Column)

	lines := e.lines()
	start, end := e.visibleLineRange(lines, height)
	for row := start; row < end; row++ {
		line := lines[row]
		lineCursorVisible := focused && cursorVisible && row == e.Cursor.Row
		fmt.Fprintf(&b, "%s\n", e.renderVisibleLine(line, contentWidth, lineCursorVisible))
	}

	return RenderPane(width, height, b.String(), focused)
}

func (e EditorPane) title(focused bool, width int) string {
	if e.Path == "" {
		return PaneTitle("editor", focused)
	}

	if focused {
		return truncateRunes(fmt.Sprintf("[ %s * ]", e.Path), width)
	}

	return truncateRunes(fmt.Sprintf("[ %s ]", e.Path), width)
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

func (e EditorPane) visibleLineRange(lines []string, height int) (int, int) {
	visibleHeight := height - 7
	if visibleHeight < 1 {
		visibleHeight = 1
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

func (e EditorPane) renderVisibleLine(line string, width int, cursorVisible bool) string {
	if width <= 0 {
		return ""
	}

	if !cursorVisible {
		return truncateRunes(line, width)
	}

	runes := []rune(line)
	column := e.Cursor.Column
	if column < 0 {
		column = 0
	}

	start := 0
	if column >= width {
		start = column - width + 1
	}

	if start > len(runes) {
		start = len(runes)
	}

	end := start + width
	if end > len(runes) {
		end = len(runes)
	}

	segment := string(runes[start:end])
	cursor := CursorPosition{
		Column: column - start,
	}

	return cursor.RenderLine(segment, true)
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
