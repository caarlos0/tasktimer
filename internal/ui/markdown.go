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

	if len(tasks) == 0 {
		return fmt.Errorf("project %s has no tasks", project)
	}

	_, _ = fmt.Fprintln(w, "# "+project+"\n")
	_, _ = fmt.Fprintf(
		w,
		"> Total time **%s**, timed between **%s** and **%s**\n\n",
		sumTasksTimes(tasks).Round(time.Second).String(),
		tasks[len(tasks)-1].StartAt.Format("2006-01-02"),
		tasks[0].EndAt.Format("2006-01-02"),
	)

	for _, task := range tasks {
		_, _ = fmt.Fprintf(
			w,
			"- [x] #%d %s - %s\n",
			task.ID+1,
			task.Title,
			task.EndAt.Sub(task.StartAt),
		)
	}

	return nil
}
