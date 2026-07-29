package app

type Pane int

const (
	PaneFiles Pane = iota
	PaneEditor
)

func (p Pane) String() string {
	switch p {
	case PaneFiles:
		return "files"
	case PaneEditor:
		return "editor"
	default:
		return "unknown"
	}
}

func (m Model) filesPaneAvailable() bool {
	return m.filesVisible && m.width >= minEditorWidth+minFilesWidth
}

func (m *Model) normalizeFocus() {
	if m.focused == PaneFiles && !m.filesPaneAvailable() {
		m.focused = PaneEditor
		m.cursorVisible = true
	}
}

func (m *Model) focusNextPane() {
	switch m.focused {
	case PaneFiles:
		m.focused = PaneEditor
	case PaneEditor:
		if m.filesPaneAvailable() {
			m.focused = PaneFiles
		} else {
			m.focused = PaneEditor
		}
	}

	m.cursorVisible = m.focused == PaneEditor || m.focused == PaneFiles
}

func (m *Model) toggleFilesPane() {
	m.filesVisible = !m.filesVisible
	if !m.filesVisible {
		m.focused = PaneEditor
		m.cursorVisible = true
	} else if m.filesPaneAvailable() {
		m.focused = PaneFiles
		m.cursorVisible = true
	} else {
		m.focused = PaneEditor
		m.cursorVisible = true
	}
}
