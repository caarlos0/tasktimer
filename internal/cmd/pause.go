package cmd

import (
	"fmt"

	"github.com/caarlos0/tasktimer/internal/store"
	"github.com/spf13/cobra"
)

type pauseCmd struct {
	cmd *cobra.Command
}

func newPauseCmd() *pauseCmd {
	cmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause the currently running task",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			project := cmd.Parent().Flag("project").Value.String()
			db, f, err := setup(project)
			if err != nil {
				return err
			}
			defer db.Close()
			defer f.Close()

			task, err := store.GetRunningTask(db)
			if err != nil {
				return fmt.Errorf("pause: %w", err)
			}

			if err := store.PauseTask(db); err != nil {
				return err
			}

			fmt.Println("paused:", task.Title)
			return nil
		},
	}
	return &pauseCmd{cmd: cmd}
}
