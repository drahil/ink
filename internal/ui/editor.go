package ui

import (
	"fmt"
	"strings"
)

type EditorPane struct {
	Content string
	Cursor  CursorPosition
}

type CursorPosition struct {
	Row    int
	Column int
}

func NewEditorPane() EditorPane {
	return EditorPane{
		Content: "<?php\n\necho 'hello';\n",
		Cursor: CursorPosition{
			Row:    0,
			Column: 0,
		},
	}
}

func (e EditorPane) View(width, height int, focused bool) string {
	var b strings.Builder

	b.WriteString(PaneTitle("editor", focused))
	b.WriteString("\n\n")
	fmt.Fprintf(&b, "cursor row: %d\n", e.Cursor.Row)
	fmt.Fprintf(&b, "cursor col: %d\n\n", e.Cursor.Column)

	for row, line := range e.lines() {
		marker := "  "
		if row == e.Cursor.Row {
			marker = "> "
		}

		fmt.Fprintf(&b, "%s%s\n", marker, line)
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
