package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type EditorPane struct {
	Content string
	Cursor  CursorPosition
}

type CursorPosition struct {
	Row    int
	Column int
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
		fmt.Fprintf(&b, "%s\n", e.renderLine(line, lineCursorVisible))
	}

	return RenderPane(width, height, b.String(), focused)
}

func (e *EditorPane) MoveCursorUp() {
	if e.Cursor.Row > 0 {
		e.Cursor.Row--
	}

	e.clampCursorColumn()
}

func (e *EditorPane) MoveCursorDown() {
	lines := e.lines()
	if e.Cursor.Row < len(lines)-1 {
		e.Cursor.Row++
	}

	e.clampCursorColumn()
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

func (e *EditorPane) clampCursorColumn() {
	lineLength := e.currentLineLength()
	if e.Cursor.Column > lineLength {
		e.Cursor.Column = lineLength
	}
}

func (e EditorPane) renderLine(line string, cursorVisible bool) string {
	if !cursorVisible {
		return line
	}

	runes := []rune(line)
	column := e.Cursor.Column

	if column < 0 {
		column = 0
	}

	if column >= len(runes) {
		return string(runes) + cursorStyle.Render(" ")
	}

	before := string(runes[:column])
	cursor := cursorStyle.Render(string(runes[column]))
	after := string(runes[column+1:])

	return before + cursor + after
}

func (e *EditorPane) InsertRune(r rune) {
	lines := e.lines()

	if e.Cursor.Row < 0 {
		e.Cursor.Row = 0
	}

	if e.Cursor.Row >= len(lines) {
		e.Cursor.Row = len(lines) - 1
	}

	currentLine := []rune(lines[e.Cursor.Row])

	if e.Cursor.Column < 0 {
		e.Cursor.Column = 0
	}

	if e.Cursor.Column > len(currentLine) {
		e.Cursor.Column = len(currentLine)
	}

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
