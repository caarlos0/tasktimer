package ui

import (
	"log"
	"time"

	"github.com/caarlos0/tasktimer/internal/model"
	"github.com/caarlos0/tasktimer/internal/store"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dgraph-io/badger/v3"
)

type taskListModel struct {
	db   *badger.DB
	list list.Model
}

func (m taskListModel) Init() tea.Cmd {
	return updateTaskListCmd(m.db)
}

var docStyle = lipgloss.NewStyle().Margin(6, 2, 0, 2)

func (m taskListModel) Update(msg tea.Msg) (taskListModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	case updateTaskListMsg:
		cmds = append(cmds, updateTaskListCmd(m.db))
	case taskListUpdatedMsg:
		var items []list.Item
		for _, t := range msg.tasks {
			items = append(items, item{
				title: t.Title,
				start: t.StartAt,
				end:   t.EndAt,
			})
		}
		m.list.SetItems(items)
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, updateProjectTimerCmd(msg.tasks), cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m taskListModel) View() string {
	return m.list.View()
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

// model

type item struct {
	title      string
	start, end time.Time
}

func (i item) Title() string { return i.title }
func (i item) Description() string {
	end := time.Now()
	if !i.end.IsZero() {
		end = i.end
	}
	return end.Sub(i.start).Round(time.Second).String()
}
func (i item) FilterValue() string { return i.title }
