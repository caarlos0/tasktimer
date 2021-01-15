package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type clockModel struct {
	t time.Time
}

func (m clockModel) Init() tea.Cmd {
	return clockTickCmd
}

func (m clockModel) Update(msg tea.Msg) (clockModel, tea.Cmd) {
	switch msg.(type) {
	case clockTickMsg:
		m.t = time.Now()
		return m, clockTickCmd
	}
	return m, nil
}

func (m clockModel) View() string {
	return m.t.Format("15:04:05")
}

// msgs and cmds

type clockTickMsg struct{}

func clockTickCmd() tea.Msg {
	time.Sleep(time.Second)
	return clockTickMsg{}
}
