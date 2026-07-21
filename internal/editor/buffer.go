package editor

import "strings"

type Position struct {
	Row    int
	Column int
}

type Buffer struct {
	Content string
	Cursor  Position
}

func NewBuffer(content string) Buffer {
	return Buffer{
		Content: content,
		Cursor: Position{
			Row:    0,
			Column: 0,
		},
	}
}

func (b *Buffer) Open(content string) {
	b.Content = content
	b.Cursor = Position{
		Row:    0,
		Column: 0,
	}
}

func (b Buffer) Lines() []string {
	lines := strings.Split(b.Content, "\n")
	if len(lines) == 0 {
		return []string{""}
	}

	return lines
}

func (b Buffer) CurrentLineLength() int {
	lines := b.Lines()
	if b.Cursor.Row < 0 || b.Cursor.Row >= len(lines) {
		return 0
	}

	return len([]rune(lines[b.Cursor.Row]))
}

func (b *Buffer) MoveCursorUp() {
	if b.Cursor.Row > 0 {
		b.Cursor.Row--
	}

	b.Cursor.ClampColumn(b.CurrentLineLength())
}

func (b *Buffer) MoveCursorDown() {
	lines := b.Lines()
	if b.Cursor.Row < len(lines)-1 {
		b.Cursor.Row++
	}

	b.Cursor.ClampColumn(b.CurrentLineLength())
}

func (b *Buffer) MoveCursorLeft() {
	if b.Cursor.Column > 0 {
		b.Cursor.Column--
	}
}

func (b *Buffer) MoveCursorRight() {
	if b.Cursor.Column < b.CurrentLineLength() {
		b.Cursor.Column++
	}
}

func (b *Buffer) MoveCursorToNextLine() {
	b.Cursor.Row++
	b.Cursor.Column = 0
}

func (b *Buffer) MoveCursorToPreviousLine() {
	b.Cursor.Row--
	b.Cursor.Column = b.CurrentLineLength()
}

func (b *Buffer) InsertRune(r rune) {
	lines := b.Lines()
	currentLine := []rune(lines[b.Cursor.Row])

	before := currentLine[:b.Cursor.Column]
	after := currentLine[b.Cursor.Column:]
	nextLine := make([]rune, 0, len(currentLine)+1)
	nextLine = append(nextLine, before...)
	nextLine = append(nextLine, r)
	nextLine = append(nextLine, after...)
	lines[b.Cursor.Row] = string(nextLine)
	b.Content = strings.Join(lines, "\n")
	b.MoveCursorRight()
}

func (b *Buffer) Backspace() {
	lines := b.Lines()

	if b.Cursor.Column == 0 {
		if b.Cursor.Row == 0 {
			return
		}

		previousLine := lines[b.Cursor.Row-1]
		currentLine := lines[b.Cursor.Row]
		b.MoveCursorToPreviousLine()

		nextLines := make([]string, 0, len(lines)-1)
		nextLines = append(nextLines, lines[:b.Cursor.Row]...)
		nextLines = append(nextLines, previousLine+currentLine)
		nextLines = append(nextLines, lines[b.Cursor.Row+2:]...)

		b.Content = strings.Join(nextLines, "\n")
		return
	}

	currentLine := []rune(lines[b.Cursor.Row])
	before := currentLine[:b.Cursor.Column-1]
	after := currentLine[b.Cursor.Column:]
	nextLine := make([]rune, 0, len(currentLine)-1)
	nextLine = append(nextLine, before...)
	nextLine = append(nextLine, after...)
	lines[b.Cursor.Row] = string(nextLine)
	b.Content = strings.Join(lines, "\n")
	b.MoveCursorLeft()
}

func (b *Buffer) InsertNewline() {
	lines := b.Lines()
	currentLine := []rune(lines[b.Cursor.Row])

	before := currentLine[:b.Cursor.Column]
	after := currentLine[b.Cursor.Column:]

	nextLines := make([]string, 0, len(lines)+1)
	nextLines = append(nextLines, lines[:b.Cursor.Row]...)
	nextLines = append(nextLines, string(before), string(after))
	nextLines = append(nextLines, lines[b.Cursor.Row+1:]...)

	b.Content = strings.Join(nextLines, "\n")
	b.MoveCursorToNextLine()
}

func (b *Buffer) MoveToFirstMatch(query string) bool {
	lines := b.Lines()
	for row, line := range lines {
		column := runeIndexOf(line, query)
		if column == -1 {
			continue
		}

		b.Cursor = Position{
			Row:    row,
			Column: column,
		}
		return true
	}

	return false
}

func (p *Position) ClampColumn(lineLength int) {
	if p.Column > lineLength {
		p.Column = lineLength
	}
}

func runeIndexOf(value, query string) int {
	if query == "" {
		return -1
	}

	byteIndex := strings.Index(strings.ToLower(value), strings.ToLower(query))
	if byteIndex == -1 {
		return -1
	}

	return len([]rune(value[:byteIndex]))
}
