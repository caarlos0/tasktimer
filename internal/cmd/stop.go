package cmd

import (
	"fmt"

	"github.com/caarlos0/tasktimer/internal/store"
	"github.com/spf13/cobra"
)

type stopCmd struct {
	cmd *cobra.Command
}

func newStopCmd() *stopCmd {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop all running timers for the given project",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			project := cmd.Parent().Flag("project").Value.String()
			db, f, err := setup(project)
			if err != nil {
				return err
			}
			defer db.Close()
			defer f.Close()

			tasks, err := store.GetTaskList(db)
			if err != nil {
				return err
			}

			var active int
			for _, t := range tasks {
				if t.EndAt.IsZero() {
					fmt.Printf("stopped: %s\n", t.Title)
					active++
				}
			}

			if active == 0 {
				fmt.Println("no running tasks")
				return nil
			}

			return store.CloseTasks(db)
		},
	}
	return &stopCmd{cmd: cmd}
}
