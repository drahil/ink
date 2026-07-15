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

	b.WriteString(PaneTitle("editor", focused))
	b.WriteString("\n\n")
	fmt.Fprintf(&b, "cursor row: %d\n", e.Cursor.Row)
	fmt.Fprintf(&b, "cursor col: %d\n\n", e.Cursor.Column)

	for row, line := range e.lines() {
		lineCursorVisible := focused && cursorVisible && row == e.Cursor.Row
		fmt.Fprintf(&b, "%s\n", e.Cursor.RenderLine(line, lineCursorVisible))
	}

	return RenderPane(width, height, b.String(), focused)
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
