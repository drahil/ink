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
