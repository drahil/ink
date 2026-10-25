package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type promptKind int

const (
	promptNone promptKind = iota
	promptQuit
)

// prompt is a one-line question shown in place of the status bar. While it
// is open, every key goes to it instead of the panes.
type prompt struct {
	kind    promptKind
	message string
}

func (p prompt) active() bool {
	return p.kind != promptNone
}

func (m *Model) openPrompt(kind promptKind, message string) {
	m.prompt = prompt{kind: kind, message: message}
}

func (m *Model) closePrompt() {
	m.prompt = prompt{}
}

func (m *Model) handlePromptKey(msg tea.KeyMsg) tea.Cmd {
	switch m.prompt.kind {
	case promptQuit:
		return m.handleQuitPromptKey(msg)
	}

	return nil
}

func (m *Model) requestQuit() tea.Cmd {
	if !m.editor.IsModified {
		return tea.Quit
	}

	name := m.editor.Path
	if name == "" {
		name = "untitled"
	}
	m.openPrompt(promptQuit, fmt.Sprintf("save changes to %s? (y)es (n)o (esc) cancel", name))
	return nil
}

func (m *Model) handleQuitPromptKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "y", "Y":
		m.closePrompt()
		cmd := m.save()
		if cmd == nil {
			// Nothing to write (e.g. no file name); status explains why.
			return nil
		}
		m.quitAfterSave = true
		return cmd
	case "n", "N":
		return tea.Quit
	case "esc", "ctrl+c", "ctrl+q":
		m.closePrompt()
		m.status = m.editorCursorStatus()
	}

	return nil
}
