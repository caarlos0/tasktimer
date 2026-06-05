package cmd

import (
	"fmt"

	"github.com/caarlos0/tasktimer/internal/store"
	"github.com/spf13/cobra"
)

type resumeCmd struct {
	cmd *cobra.Command
}

func newResumeCmd() *resumeCmd {
	cmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume the currently paused task",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			project := cmd.Parent().Flag("project").Value.String()
			db, f, err := setup(project)
			if err != nil {
				return err
			}
			defer db.Close()
			defer f.Close()

			task, err := store.GetPausedTask(db)
			if err != nil {
				return fmt.Errorf("resume: %w", err)
			}

			if err := store.ResumeTask(db); err != nil {
				return err
			}

			fmt.Println("resumed:", task.Title)
			return nil
		},
	}
	return &resumeCmd{cmd: cmd}
}
