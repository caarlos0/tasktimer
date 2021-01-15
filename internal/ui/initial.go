package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
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
	task.Placeholder = "New task description..."
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
	clock clockModel
	task  textinput.Model
	db    *badger.DB
	tasks []model.Task
	err   error
}

func (m initialModel) Init() tea.Cmd {
	return tea.Batch(m.clock.Init(), textinput.Blink, updateTaskListCmd(m.db))
}

func (m initialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case errMsg:
		m.err = msg.error
	case updateTaskListMsg:
		cmds = append(cmds, updateTaskListCmd(m.db))
	case taskListUpdatedMsg:
		m.tasks = msg.tasks
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, closeTasks(m.db, tea.Quit)
		case "enter":
			log.Println("start/stop task")
			cmds = append(cmds, closeTasks(m.db, createTask(m.db, strings.TrimSpace(m.task.Value()))))
			m.task.SetValue("")
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

type taskListUpdatedMsg struct {
	tasks []model.Task
}

func updateTaskListCmd(db *badger.DB) tea.Cmd {
	return func() tea.Msg {
		log.Println("updating task list")
		var tasks []model.Task
		if err := db.View(func(txn *badger.Txn) error {
			var it = txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()
			for it.Seek([]byte(prefix)); it.ValidForPrefix([]byte(prefix)); it.Next() {
				var item = it.Item()
				err := item.Value(func(v []byte) error {
					var task model.Task
					if err := json.Unmarshal(v, &task); err != nil {
						return err
					}
					tasks = append(tasks, task)
					return nil
				})
				if err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return errMsg{err}
		}
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].StartAt.After(tasks[j].StartAt)
		})
		log.Println("loaded", len(tasks), "tasks")
		return taskListUpdatedMsg{tasks}
	}
}

type updateTaskListMsg struct{}

func closeTasks(db *badger.DB, then tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		log.Println("closing tasks")
		if err := db.Update(func(txn *badger.Txn) error {
			var it = txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()
			for it.Seek([]byte(prefix)); it.ValidForPrefix([]byte(prefix)); it.Next() {
				var item = it.Item()
				var k = item.Key()
				err := item.Value(func(v []byte) error {
					var task model.Task
					if err := json.Unmarshal(v, &task); err != nil {
						return err
					}
					if !task.EndAt.IsZero() {
						return nil
					}
					task.EndAt = time.Now()
					log.Println("closing", task.Title)
					return txn.Set(k, task.Bytes())
				})
				if err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return errMsg{err}
		}
		return then()
	}
}

var prefix = "tasks."

func createTask(db *badger.DB, t string) tea.Cmd {
	return func() tea.Msg {
		if t == "" {
			return updateTaskListMsg{}
		}

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
		return updateTaskListMsg{}
	}
}

func (m initialModel) View() string {
	if m.err != nil {
		return "\n" + m.err.Error() + "\n"
	}
	var s = m.clock.View() + " - Project Foo\n\n"
	s += m.task.View()
	return s + "\n\n---\n\n" + taskView(m.tasks) + "\n"
}

func taskView(tt []model.Task) string {
	var s string
	for _, t := range tt {
		var z = time.Now()
		var icon = iconOngoing
		var decorate = bold
		if !t.EndAt.IsZero() {
			z = t.EndAt
			icon = iconDone
			decorate = faint
		}
		s += decorate(fmt.Sprintf("%s %s (%s)", icon, t.Title, z.Sub(t.StartAt).Round(time.Second))) + "\n"
	}
	return s
}
