package cmd

import (
	"os"

	"github.com/caarlos0/tasktimer/internal/ui"
	"github.com/spf13/cobra"
)

type toJSONCmd struct {
	cmd *cobra.Command
}

func newToJSONCmd() *toJSONCmd {
	cmd := &cobra.Command{
		Use:   "to-json",
		Short: "Exports the database as JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			project := cmd.Parent().Flag("project").Value.String()
			db, f, err := setup(project)
			if err != nil {
				return err
			}
			defer db.Close()
			defer f.Close()

			return ui.WriteProjectJSON(db, project, os.Stdout)
		},
	}

	return &toJSONCmd{cmd: cmd}
}
