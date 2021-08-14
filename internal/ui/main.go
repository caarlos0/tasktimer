package ui

import (
	"log"
	"strings"

	"github.com/caarlos0/tasktimer/internal/store"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dgraph-io/badger/v3"
)

func Init(db *badger.DB, project string) tea.Model {
	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.Placeholder = "New task description..."
	input.Focus()
	input.CharLimit = 250
	input.Width = 50

	l := list.NewModel([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Task List"
	l.KeyMap.Quit.SetEnabled(false)

	return mainModel{
		list: taskListModel{
			db:   db,
			list: l,
		},
		timer:   projectTimerModel{},
		db:      db,
		input:   input,
		project: project,
	}
}

type mainModel struct {
	input   textinput.Model
	list    taskListModel
	timer   projectTimerModel
	db      *badger.DB
	project string
	err     error
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.list.Init(), textinput.Blink)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg2 := msg.(type) {
	case errMsg:
		m.err = msg2.error
	case tea.KeyMsg:
		switch msg2.String() {
		case "ctrl+c":
			return m, tea.Sequentially(closeTasks(m.db), tea.Quit)
		case "esc":
			if m.list.list.SettingFilter() {
				m.list.list, cmd = m.list.list.Update(msg)
				cmds = append(cmds, cmd)
			} else {
				if m.input.Focused() {
					m.input.Blur()
				}
				log.Println("stop timer")
				cmds = append(cmds, tea.Sequentially(closeTasks(m.db), updateTaskListCmd(m.db)))
			}
		case "enter":
			if m.list.list.SettingFilter() {
				m.list.list, cmd = m.list.list.Update(msg)
				cmds = append(cmds, cmd)
			} else {
				if !m.input.Focused() {
					m.input.Focus()
				} else {
					log.Println("start/stop timer")
					cmds = append(cmds, tea.Sequentially(closeTasks(m.db), createTask(m.db, strings.TrimSpace(m.input.Value()))))
					m.input.SetValue("")
				}
			}
		default:
			if m.input.Focused() {
				// only send key presses to input if it is focused
				m.input, cmd = m.input.Update(msg)
				cmds = append(cmds, cmd)
				msg = nil
			} else {
				m.list.list, cmd = m.list.list.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	default:
		// if its not a keypress, we gotta update the list
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
		m.list.list, cmd = m.list.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.timer, cmd = m.timer.Update(msg)
	cmds = append(cmds, cmd)
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	if m.err != nil {
		return "\n" +
			errorFaintForeground.Render("Oops, something went wrong:") +
			"\n\n" +
			errorForegroundPadded.Render(m.err.Error()) +
			"\n\n" +
			errorFaintForeground.Render("Check the logs for more details...")
	}
	return secondaryForeground.Render("project: ") +
		activeForegroundBold.Render(m.project) +
		separator + m.timer.View() + "\n\n" +
		m.input.View() + "\n\n" +
		m.list.View() + "\n"
}

// cmds

func closeTasks(db *badger.DB) tea.Cmd {
	return func() tea.Msg {
		log.Println("closing tasks")
		if err := store.CloseTasks(db); err != nil {
			return errMsg{err}
		}
		return nil
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
