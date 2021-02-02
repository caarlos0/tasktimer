package ui

import (
	"encoding/json"
	"io"

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

	bts, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	_, err = w.Write(bts)

	return err
}
