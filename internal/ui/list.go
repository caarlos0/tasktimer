package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/tasktimer/internal/model"
	"github.com/caarlos0/tasktimer/internal/store"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dgraph-io/badger/v3"
)

type taskListModel struct {
	db       *badger.DB
	viewport viewport.Model
	ready    bool
}

func (m taskListModel) Init() tea.Cmd {
	return updateTaskListCmd(m.db)
}

func (m taskListModel) Update(msg tea.Msg) (taskListModel, tea.Cmd) {
	const offset = 7
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.Model{
				Width:  msg.Width,
				Height: msg.Height - offset,
			}
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - offset
		}
	case updateTaskListMsg:
		cmds = append(cmds, updateTaskListCmd(m.db))
	case taskListUpdatedMsg:
		m.viewport.SetContent(taskList(msg.tasks))
		cmds = append(cmds, updateProjectTimerCmd(msg.tasks))
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m taskListModel) View() string {
	return m.viewport.View()
}

func taskList(tasks []model.Task) string {
	var s string
	for _, t := range tasks {
		var z = time.Now()
		var icon = iconOngoing
		var decorate = bold
		if !t.EndAt.IsZero() {
			z = t.EndAt
			icon = iconDone
			decorate = faint
		}
		s += decorate(fmt.Sprintf("%s #%d %s (%s)", icon, t.ID+1, t.Title, secondaryForeground(z.Sub(t.StartAt).Round(time.Second).String()))) + "\n"
	}
	return s
}

// msgs

type updateTaskListMsg struct{}

type taskListUpdatedMsg struct {
	tasks []model.Task
}

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

// cmds

func updateTaskListCmd(db *badger.DB) tea.Cmd {
	return func() tea.Msg {
		log.Println("updating input list")
		tasks, err := store.GetTaskList(db)
		if err != nil {
			return errMsg{err}
		}
		return taskListUpdatedMsg{tasks}
	}
}
