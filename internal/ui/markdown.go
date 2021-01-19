package ui

import (
	"fmt"
	"io"
	"time"

	"github.com/caarlos0/tasktimer/internal/store"
	"github.com/dgraph-io/badger/v3"
)

// WriteProjectMarkdown writes the project task list in markdown format to the given
// io.Writer.
func WriteProjectMarkdown(db *badger.DB, project string, w io.Writer) error {
	tasks, err := store.GetTaskList(db)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintln(w, "# "+project+"\n")
	_, _ = fmt.Fprintf(w, "> Total time: %s\n\n", sumTasksTimes(tasks).Round(time.Second).String())
	for _, task := range tasks {
		_, _ = fmt.Fprintf(
			w,
			"- %s (#%d) - %s\n",
			task.Title,
			task.ID,
			task.EndAt.Sub(task.StartAt),
		)
	}

	return nil
}
