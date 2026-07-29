package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type cursorBlinkMsg time.Time

type autosaveMsg struct {
	version int
}

type saveCompleteMsg struct {
	version int
	err     error
}

const cursorBlinkInterval = 500 * time.Millisecond

func blinkCursor() tea.Cmd {
	return tea.Tick(cursorBlinkInterval, func(t time.Time) tea.Msg {
		return cursorBlinkMsg(t)
	})
}
