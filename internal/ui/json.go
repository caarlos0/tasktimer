package ui

import (
	"encoding/json"
	"io"

	"github.com/caarlos0/tasktimer/internal/model"
	"github.com/caarlos0/tasktimer/internal/store"
	"github.com/dgraph-io/badger/v3"
)

// WriteProjectJSON writes the project task list in JSON format to the given
// io.Writer.
func WriteProjectJSON(db *badger.DB, project string, w io.Writer) error {
	tasks, err := store.GetTaskList(db)
	if err != nil {
		return err
	}

	var expTasks []model.ExportedTask
	for _, t := range tasks {
		expTasks = append(expTasks, model.ExportedTask{
			Title:   t.Title,
			StartAt: t.StartAt,
			EndAt:   t.EndAt,
		})
	}
	bts, err := json.MarshalIndent(expTasks, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write(bts)

	return err
}

func WriteProjectJSONGrouped(db *badger.DB, project string, w io.Writer) error {
	tasks, err := store.GetTaskList(db)
	if err != nil {
		return err
	}

	expTasksMap := model.ExportedTasksGrouped{}
	for _, t := range tasks {
		expTasksMap[t.Title] = append(expTasksMap[t.Title], model.ExportedTask{
			StartAt: t.StartAt,
			EndAt:   t.EndAt,
		})
	}

	bts, err := json.MarshalIndent(expTasksMap, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write(bts)

	return err
}
