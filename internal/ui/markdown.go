package ui

import (
	"fmt"
	"io"
	"time"

	"github.com/caarlos0/tasktimer/internal/store"
	"github.com/dgraph-io/badger/v3"
)

// WriteProjectMarkdown writes the project task list in markdown format to the given
// io.Writer. Pass zero-value time.Time for since/until to include all tasks.
func WriteProjectMarkdown(db *badger.DB, project string, w io.Writer, since, until time.Time) error {
	tasks, err := store.GetTaskList(db)
	if err != nil {
		return err
	}

	if !since.IsZero() || !until.IsZero() {
		filtered := tasks[:0]
		for _, t := range tasks {
			if !since.IsZero() && t.StartAt.Before(since) {
				continue
			}
			if !until.IsZero() && !t.StartAt.Before(until) {
				continue
			}
			filtered = append(filtered, t)
		}
		tasks = filtered
	}

	if len(tasks) == 0 {
		return fmt.Errorf("project %s has no tasks", project)
	}

	_, _ = fmt.Fprintln(w, "# "+project+"\n")
	now := time.Now()
	_, _ = fmt.Fprintf(
		w,
		"> Total time **%s**, timed between **%s** and **%s**\n\n",
		SumTasksTimes(tasks, time.Time{}).Round(time.Second).String(),
		tasks[len(tasks)-1].StartAt.Format("2006-01-02"),
		effectiveEndAt(tasks[0].EndAt, now).Format("2006-01-02"),
	)

	for _, task := range tasks {
		duration := taskDuration(task, now).Round(time.Second)
		endStr := "in progress"
		if !task.EndAt.IsZero() {
			endStr = task.EndAt.Format("2006-01-02 15:04:05")
		} else if !task.PausedAt.IsZero() {
			endStr = "paused"
		}
		_, _ = fmt.Fprintf(
			w,
			"- **#%d** %s - _%s_ - _%s (%s)_\n",
			task.ID+1,
			task.Title,
			task.StartAt.Format("2006-01-02 15:04:05"),
			endStr,
			duration,
		)
	}

	return nil
}
