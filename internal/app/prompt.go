package app

import (
	"fmt"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
)

type promptKind int

const (
	promptNone promptKind = iota
	promptQuit
)

type prompt struct {
	kind    promptKind
	title   string
	items   []string
	message string
}

func (p prompt) active() bool {
	return p.kind != promptNone
}

func (m *Model) openPrompt(kind promptKind, title string, items []string, message string) {
	m.prompt = prompt{kind: kind, title: title, items: items, message: message}
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
	var modified []string
	if m.editor.IsModified {
		name := m.editor.Path
		if name == "" {
			name = "untitled"
		}
		modified = append(modified, name)
	}
	for path, pane := range m.allOpenedFiles {
		if pane.IsModified {
			modified = append(modified, path)
		}
	}

	if len(modified) == 0 {
		return tea.Quit
	}

	sort.Strings(modified)
	title := fmt.Sprintf("unsaved changes in %d files:", len(modified))
	if len(modified) == 1 {
		title = "unsaved changes in 1 file:"
	}
	m.openPrompt(promptQuit, title, modified, "save all? (y)es (n)o (esc) cancel")
	return nil
}

func (m *Model) handleQuitPromptKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "y", "Y":
		m.closePrompt()
		if m.editor.IsModified && m.editor.Path == "" {
			m.status = "untitled buffer has no file name; quit and press n to discard it"
			return nil
		}
		cmd := m.saveAll()
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
