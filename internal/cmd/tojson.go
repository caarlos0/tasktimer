package cmd

import (
	"os"

	"github.com/caarlos0/tasktimer/internal/ui"
	"github.com/muesli/coral"
)

type toJSONCmd struct {
	cmd *coral.Command
}

func newToJSONCmd() *toJSONCmd {
	cmd := &coral.Command{
		Use:   "to-json",
		Short: "Exports the database as JSON",
		Args:  coral.MaximumNArgs(1),
		RunE: func(cmd *coral.Command, args []string) error {
			project := cmd.Parent().Flag("project").Value.String()
			db, f, err := setup(project)
			if err != nil {
				return err
			}
			defer db.Close()
			defer f.Close()

			if len(args) > 0 {
				f, err := os.OpenFile(args[0], os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0o666)
				if err != nil {
					return err
				}
				defer f.Close()
				return ui.WriteProjectJSON(db, project, f)
			}

			return ui.WriteProjectJSON(db, project, os.Stdout)
		},
	}

	return &toJSONCmd{cmd: cmd}
}
