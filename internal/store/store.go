package store

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/caarlos0/tasktimer/internal/model"
	"github.com/dgraph-io/badger/v3"
)

var (
	prefix     = []byte("tasks.")
	sequenceID = []byte("tasks_seq")
)

func GetTaskList(db *badger.DB) ([]model.Task, error) {
	var tasks []model.Task
	if err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				var task model.Task
				if err := json.Unmarshal(v, &task); err != nil {
					return err
				}
				sort.Slice(task.Durations, func(i, j int) bool {
					return task.Durations[i].StartAt.Before(task.Durations[j].StartAt)
				})
				tasks = append(tasks, task)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return tasks, err
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].FirstStartedAt.After(tasks[j].FirstStartedAt)
	})
	return tasks, nil
}

func CloseTasks(db *badger.DB) error {
	return db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				var task model.Task
				if err := json.Unmarshal(v, &task); err != nil {
					return err
				}

				if len(task.Durations) == 0 {
					return nil
				}

				lastDuration := task.Durations[len(task.Durations)-1]
				if !lastDuration.EndAt.IsZero() {
					return nil
				}
				lastDuration.EndAt = time.Now().Truncate(time.Second)
				task.LastEndedAt = time.Now().Truncate(time.Second)
				lastDuration.Duration = lastDuration.EndAt.Sub(lastDuration.StartAt)
				task.Total += lastDuration.Duration
				log.Println("closing", task.Title)
				log.Println("total", task.Total)
				return txn.Set(k, task.Bytes())
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func GetTaskId(id uint64) string {
	return string(prefix) + strconv.FormatUint(id, 10)
}

func CreateTask(db *badger.DB, t string) error {
	if t == "" {
		return nil
	}

	return db.Update(func(txn *badger.Txn) error {
		seq, err := db.GetSequence(sequenceID, 100)
		if err != nil {
			return err
		}
		defer seq.Release()
		s, err := seq.Next()
		if err != nil {
			return err
		}
		id := string(prefix) + strconv.FormatUint(s, 10)
		log.Println("creating task:", id, "->", t)
		return txn.Set([]byte(id), model.Task{
			ID:             s,
			Title:          t,
			FirstStartedAt: time.Now().Truncate(time.Second),
			Durations: []*model.TaskDuration{
				{StartAt: time.Now().Truncate(time.Second)},
			},
		}.Bytes())
	})
}

func NewTaskDuration(db *badger.DB, id uint64) error {
	sid := GetTaskId(id)
	return db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(sid))
		if err != nil {
			return err
		}
		var task model.Task
		if err := item.Value(func(v []byte) error {
			return json.Unmarshal(v, &task)
		}); err != nil {
			return err
		}
		task.Durations = append(task.Durations, &model.TaskDuration{
			StartAt: time.Now().Truncate(time.Second),
		})
		log.Println("new duration for", task.Title)

		return txn.Set([]byte(sid), task.Bytes())
	})
}

func LoadTasks(db *badger.DB, tasks []model.ExportedTask) error {
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].FirstStartedAt.Before(tasks[j].FirstStartedAt)
	})
	return db.Update(func(txn *badger.Txn) error {
		seq, err := db.GetSequence(sequenceID, 100)
		if err != nil {
			return err
		}
		defer seq.Release()

		for _, t := range tasks {
			s, err := seq.Next()
			if err != nil {
				return err
			}
			id := string(prefix) + strconv.FormatUint(s, 10)
			log.Println("creating task:", id, "->", t)
			if err := txn.Set([]byte(id), model.Task{
				ID:             s,
				Title:          t.Title,
				FirstStartedAt: t.FirstStartedAt,
				LastEndedAt:    t.LastEndedAt,
			}.Bytes()); err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}
		}

		return nil
	})
}
