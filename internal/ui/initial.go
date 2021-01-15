package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/tasktimer/internal/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
)

var _ tea.Model = initialModel{}

func Init(db *badger.DB) tea.Model {
	var task = textinput.NewModel()
	task.Placeholder = "This UI shrug"
	task.Focus()
	task.CharLimit = 156
	task.Width = 50

	return initialModel{
		clock: clockModel{time.Now()},
		db:    db,
		task:  task,
	}
}

type initialModel struct {
	clock   clockModel
	task    textinput.Model
	db      *badger.DB
	current model.Task
}

func (m initialModel) Init() tea.Cmd {
	return tea.Batch(m.clock.Init(), textinput.Blink)
}

func (m initialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.task.Focused() {
				log.Println("start task")
				cmds = append(cmds, createTask(m.db, m.task.Value()))
				m.task.SetValue("")
				m.task.Blur()
			}
		case "n":
			if !m.task.Focused() {
				m.task.Focus()
			}
		}
	}

	var cmd tea.Cmd
	m.task, cmd = m.task.Update(msg)
	cmds = append(cmds, cmd)
	clock, cmd := m.clock.Update(msg)
	m.clock = clock.(clockModel)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

type updateTaskList struct{}

func createTask(db *badger.DB, t string) tea.Cmd {
	return func() tea.Msg {
		var id = "tasks." + uuid.New().String()
		if err := db.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(id), model.Task{
				Title:   t,
				StartAt: time.Now(),
			}.Bytes())
		}); err != nil {
			return errMsg{err}
		}
		log.Println("created task:", id, "->", t)
		return updateTaskList{}
	}
}

func (m initialModel) View() string {
	return fmt.Sprintf(
		"%s\n\nWhat are you working on?\n\n%s\n\n%s\n\n%s",
		m.clock.View(),
		m.task.View(),
		taskView(m.current),
		"(esc to quit)",
	) + "\n"
}

func taskView(t model.Task) string {
	if !t.EndAt.IsZero() {
		return fmt.Sprintf("%s (%s)", t.Title, t.EndAt.Sub(t.StartAt).Round(time.Second))
	}
	if !t.StartAt.IsZero() {
		return fmt.Sprintf("%s (%s...)", t.Title, time.Since(t.StartAt).Round(time.Second))
	}
	return ""
}
