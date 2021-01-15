package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/caarlos0/tasktimer/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dgraph-io/badger/v3"
)

type taskListModel struct {
	db    *badger.DB
	tasks []model.Task
}

func (m taskListModel) Init() tea.Cmd {
	return updateTaskListCmd(m.db)
}

func (m taskListModel) Update(msg tea.Msg) (taskListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case updateTaskListMsg:
		return m, updateTaskListCmd(m.db)
	case taskListUpdatedMsg:
		m.tasks = msg.tasks
	}
	return m, nil
}

func (m taskListModel) View() string {
	var s string
	for _, t := range m.tasks {
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

type updateTaskListMsg struct{}

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
