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
	var grouped bool

	cmd := &cobra.Command{
		Use:   "to-json",
		Short: "Exports the database as JSON",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project := cmd.Parent().Flag("project").Value.String()
			db, f, err := setup(project)
			if err != nil {
				return err
			}
			defer db.Close()
			defer f.Close()

			var fi *os.File = os.Stdout
			if len(args) > 0 {
				fi, err = os.OpenFile(args[0], os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0o666)
				if err != nil {
					return err
				}
				defer fi.Close()
			}

			if grouped {
				return ui.WriteProjectJSONGrouped(db, project, fi)
			}

			return ui.WriteProjectJSON(db, project, fi)
		},
	}

	jcmd := toJSONCmd{cmd: cmd}
	jcmd.cmd.Flags().BoolVarP(&grouped, "grouped", "g", false, "grouped by desc")
	return &jcmd
}
