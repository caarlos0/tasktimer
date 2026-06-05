package cmd

import (
	"fmt"
	"strings"

	"github.com/caarlos0/tasktimer/internal/store"
	"github.com/spf13/cobra"
)

type startCmd struct {
	cmd *cobra.Command
}

func newStartCmd() *startCmd {
	cmd := &cobra.Command{
		Use:   "start <title>",
		Short: "Stop any running timer and start a new task",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project := cmd.Parent().Flag("project").Value.String()
			db, f, err := setup(project)
			if err != nil {
				return err
			}
			defer db.Close()
			defer f.Close()

			if err := store.CloseTasks(db); err != nil {
				return err
			}

			title := strings.Join(args, " ")
			if err := store.CreateTask(db, title); err != nil {
				return err
			}

			fmt.Printf("started: %s\n", title)
			return nil
		},
	}
	return &startCmd{cmd: cmd}
}
