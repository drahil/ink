package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type TextInput struct {
	Value  string
	Cursor int
}

func (i *TextInput) InsertRune(r rune) {
	runes := []rune(i.Value)
	if i.Cursor < 0 {
		i.Cursor = 0
	}
	if i.Cursor > len(runes) {
		i.Cursor = len(runes)
	}

	next := make([]rune, 0, len(runes)+1)
	next = append(next, runes[:i.Cursor]...)
	next = append(next, r)
	next = append(next, runes[i.Cursor:]...)
	i.Value = string(next)
	i.MoveRight()
}

func (i *TextInput) Backspace() {
	if i.Cursor <= 0 {
		return
	}

	runes := []rune(i.Value)
	if i.Cursor > len(runes) {
		i.Cursor = len(runes)
	}
	if i.Cursor <= 0 {
		return
	}

	next := make([]rune, 0, len(runes)-1)
	next = append(next, runes[:i.Cursor-1]...)
	next = append(next, runes[i.Cursor:]...)
	i.Value = string(next)
	i.MoveLeft()
}

func (i *TextInput) Clear() {
	i.Value = ""
	i.Cursor = 0
}

func (i *TextInput) MoveLeft() {
	if i.Cursor > 0 {
		i.Cursor--
	}
}

func (i *TextInput) MoveRight() {
	if i.Cursor < len([]rune(i.Value)) {
		i.Cursor++
	}
}

func (i TextInput) RenderWithin(width int, cursorVisible bool) string {
	if width <= 0 {
		return ""
	}

	if !cursorVisible {
		return truncateCells(i.Value, width)
	}

	column := i.Cursor
	if column < 0 {
		column = 0
	}

	start := 0
	if column >= width {
		start = column - width + 1
	}

	segment := ansi.Cut(i.Value, start, start+width)
	cursorColumn := column - start
	before := ansi.Cut(segment, 0, cursorColumn)
	cursorText := ansi.Cut(segment, cursorColumn, cursorColumn+1)
	afterStart := cursorColumn + lipgloss.Width(cursorText)
	after := ansi.Cut(segment, afterStart, width)

	if cursorText == "" {
		cursorText = " "
	}

	return truncateCells(before+cursorStyle.Render(cursorText)+after, width)
}
