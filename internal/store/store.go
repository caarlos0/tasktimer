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
				task.EndAt = time.Now().Truncate(time.Second)
				log.Println("closing", task.Title)
				return txn.Set(k, task.Bytes())
			})
			if err != nil {
				return err
			}
		}
		return nil
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
		return txn.Set([]byte(id), model.Task{
			ID:      s,
			Title:   t,
			StartAt: time.Now().Truncate(time.Second),
		}.Bytes())
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
			if err := txn.Set([]byte(id), model.Task{
				ID:      s,
				Title:   t.Title,
				StartAt: t.StartAt,
				EndAt:   t.EndAt,
			}.Bytes()); err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}
		}

		return nil
	})
}
