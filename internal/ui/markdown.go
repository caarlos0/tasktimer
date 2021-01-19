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
	_, _ = fmt.Fprintf(w, "> Total time: %s\n", sumTasksTimes(tasks).Round(time.Second).String())
	_, _ = fmt.Fprintf(w, "> From: %s\n", tasks[len(tasks)-1].StartAt.Format(time.UnixDate))
	_, _ = fmt.Fprintf(w, "> To: %s\n\n", tasks[0].EndAt.Format(time.UnixDate))

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
