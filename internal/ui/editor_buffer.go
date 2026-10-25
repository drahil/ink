package ui

import (
	"regexp"
	"strings"
)

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

func (e *EditorPane) Tab() {
	lines := e.lines()
	currentLine := []rune(lines[e.Cursor.Row])
	tabSpaces := []rune("    ")

	before := currentLine[:e.Cursor.Column]
	after := currentLine[e.Cursor.Column:]
	nextLine := make([]rune, 0, len(currentLine)+len(tabSpaces))
	nextLine = append(nextLine, before...)
	nextLine = append(nextLine, tabSpaces...)
	nextLine = append(nextLine, after...)
	lines[e.Cursor.Row] = string(nextLine)
	e.Content = strings.Join(lines, "\n")
	e.Cursor.Column += len(tabSpaces)
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
	e.IsModified = false
	e.Cursor = CursorPosition{
		Row:    0,
		Column: 0,
	}
	e.IsSearchActive = false
	e.SearchQuery = ""
	e.searchMatchEnd = -1
}

func (e *EditorPane) InitSearch() {
	e.IsSearchActive = true
	e.searchMatchEnd = -1
}

func (e *EditorPane) AddToSearchQuery(r rune) {
	e.SearchQuery += string(r)
	e.searchMatchEnd = -1
}

func (e *EditorPane) BackspaceInSearch() {
	query := []rune(e.SearchQuery)
	if len(query) == 0 {
		return
	}

	e.SearchQuery = string(query[:len(query)-1])
	e.searchMatchEnd = -1
}

func (e *EditorPane) CloseSearch() {
	e.IsSearchActive = false
	e.searchMatchEnd = -1
}

func (e *EditorPane) Search() bool {
	if e.SearchQuery == "" {
		return false
	}

	pattern := regexp.MustCompile("(?i)" + regexp.QuoteMeta(e.SearchQuery))
	matches := pattern.FindAllStringIndex(e.Content, -1)
	if len(matches) == 0 {
		e.searchMatchEnd = -1
		return false
	}

	match := matches[0]
	for _, candidate := range matches {
		if candidate[0] >= e.searchMatchEnd {
			match = candidate
			break
		}
	}

	e.searchMatchEnd = match[1]
	before := e.Content[:match[0]]
	e.Cursor.Row = strings.Count(before, "\n")
	if lastNewline := strings.LastIndex(before, "\n"); lastNewline >= 0 {
		before = before[lastNewline+1:]
	}
	e.Cursor.Column = len([]rune(before))
	return true
}
