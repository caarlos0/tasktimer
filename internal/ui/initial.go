package ui

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/caarlos0/tasktimer/internal/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
)

func Init(db *badger.DB) tea.Model {
	var task = textinput.NewModel()
	task.Placeholder = "New task description..."
	task.Focus()
	task.CharLimit = 156
	task.Width = 50

	return mainModel{
		clock: clockModel{time.Now()},
		list: taskListModel{
			db: db,
		},
		db:   db,
		task: task,
	}
}

type mainModel struct {
	clock clockModel
	task  textinput.Model
	list  taskListModel
	db    *badger.DB
	err   error
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.list.Init(), m.clock.Init(), textinput.Blink, updateTaskListCmd(m.db))
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case errMsg:
		m.err = msg.error
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
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	m.task, cmd = m.task.Update(msg)
	cmds = append(cmds, cmd)
	m.clock, cmd = m.clock.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

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

func (m mainModel) View() string {
	if m.err != nil {
		return "\n" + m.err.Error() + "\n"
	}
	var s = m.clock.View() + " - Project Foo\n\n"
	s += m.task.View()
	return s + "\n\n---\n\n" + m.list.View() + "\n"
}
