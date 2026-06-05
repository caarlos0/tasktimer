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

func setTask(txn *badger.Txn, key []byte, task model.Task) error {
	bts, err := task.Bytes()
	if err != nil {
		return err
	}
	return txn.Set(key, bts)
}

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
		return tasks[i].StartAt.After(tasks[j].StartAt)
	})
	return tasks, nil
}

func GetRunningTask(db *badger.DB) (model.Task, error) {
	tasks, err := GetTaskList(db)
	if err != nil {
		return model.Task{}, err
	}
	for _, t := range tasks {
		if t.EndAt.IsZero() && t.PausedAt.IsZero() {
			return t, nil
		}
	}
	return model.Task{}, fmt.Errorf("no running task")
}

func GetPausedTask(db *badger.DB) (model.Task, error) {
	tasks, err := GetTaskList(db)
	if err != nil {
		return model.Task{}, err
	}
	for _, t := range tasks {
		if t.EndAt.IsZero() && !t.PausedAt.IsZero() {
			return t, nil
		}
	}
	return model.Task{}, fmt.Errorf("no paused task")
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
				if !task.EndAt.IsZero() {
					return nil
				}
				if !task.PausedAt.IsZero() {
					task.EndAt = task.PausedAt.Truncate(time.Second)
					task.PausedAt = time.Time{}
				} else {
					task.EndAt = time.Now().Truncate(time.Second)
				}
				log.Println("closing", task.Title)
				return setTask(txn, k, task)
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func PauseTask(db *badger.DB) error {
	return db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			var found bool
			err := item.Value(func(v []byte) error {
				var task model.Task
				if err := json.Unmarshal(v, &task); err != nil {
					return err
				}
				if task.EndAt.IsZero() && task.PausedAt.IsZero() {
					task.PausedAt = time.Now().Truncate(time.Second)
					found = true
					return setTask(txn, k, task)
				}
				return nil
			})
			if err != nil {
				return err
			}
			if found {
				return nil
			}
		}
		return fmt.Errorf("no running task to pause")
	})
}

func ResumeTask(db *badger.DB) error {
	return db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			var found bool
			err := item.Value(func(v []byte) error {
				var task model.Task
				if err := json.Unmarshal(v, &task); err != nil {
					return err
				}
				if task.EndAt.IsZero() && !task.PausedAt.IsZero() {
					task.PausedFor += time.Now().Truncate(time.Second).Sub(task.PausedAt)
					task.PausedAt = time.Time{}
					found = true
					return setTask(txn, k, task)
				}
				return nil
			})
			if err != nil {
				return err
			}
			if found {
				return nil
			}
		}
		return fmt.Errorf("no paused task to resume")
	})
}

func DeleteTask(db *badger.DB, id uint64) error {
	key := []byte(string(prefix) + strconv.FormatUint(id, 10))
	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
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
		return setTask(txn, []byte(id), model.Task{
			ID:      s,
			Title:   t,
			StartAt: time.Now().Truncate(time.Second),
		})
	})
}

func LoadTasks(db *badger.DB, tasks []model.ExportedTask) error {
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].StartAt.Before(tasks[j].StartAt)
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
			if err := setTask(txn, []byte(id), model.Task{
				ID:      s,
				Title:   t.Title,
				StartAt: t.StartAt,
				EndAt:   t.EndAt,
			}); err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}
		}

		return nil
	})
}
