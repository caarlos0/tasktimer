package ui

import (
	"log"
	"strings"
	"time"

	"github.com/caarlos0/tasktimer/internal/store"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dgraph-io/badger/v3"
)

func Init(db *badger.DB, project string) tea.Model {
	var input = textinput.NewModel()
	input.Placeholder = "New task description..."
	input.Focus()
	input.CharLimit = 156
	input.Width = 50

	return mainModel{
		clock: clockModel{time.Now()},
		list: taskListModel{
			db: db,
		},
		timer:   projectTimerModel{},
		db:      db,
		input:   input,
		project: project,
	}
}

type mainModel struct {
	clock   clockModel
	input   textinput.Model
	list    taskListModel
	timer   projectTimerModel
	db      *badger.DB
	project string
	err     error
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.list.Init(), m.clock.Init(), textinput.Blink)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case errMsg:
		m.err = msg.error
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, closeTasks(m.db, tea.Quit)
		case "esc":
			m.input.Blur()
			cmds = append(cmds, closeTasks(m.db, updateTaskListCmd(m.db)))
		case "enter":
			if !m.input.Focused() {
				m.input.Focus()
			} else {
				log.Println("start/stop input")
				cmds = append(cmds, closeTasks(m.db, createTask(m.db, strings.TrimSpace(m.input.Value()))))
				m.input.SetValue("")
			}
		}
	}

	var cmd tea.Cmd
	m.timer, cmd = m.timer.Update(msg)
	cmds = append(cmds, cmd)
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)
	m.clock, cmd = m.clock.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	if m.err != nil {
		return "\n" + m.err.Error() + "\n"
	}
	var s = m.clock.View() + " - " + m.project + " - " + m.timer.View() + "\n\n"
	s += m.input.View() + "\n\n"
	return s + m.list.View() + "\n"
}

// cmds

func closeTasks(db *badger.DB, then tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		log.Println("closing tasks")
		if err := store.CloseTasks(db); err != nil {
			return errMsg{err}
		}
		return then()
	}
}

func createTask(db *badger.DB, t string) tea.Cmd {
	return func() tea.Msg {
		if err := store.CreateTask(db, t); err != nil {
			return errMsg{err}
		}
		return updateTaskListMsg{}
	}
}
