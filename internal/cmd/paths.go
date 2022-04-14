package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type pathsCmd struct {
	cmd *cobra.Command
}

func newPathsCmd() *pathsCmd {
	cmd := &cobra.Command{
		Use:   "paths",
		Short: "Print the paths being used for logs, data et al",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			project := cmd.Parent().Flag("project").Value.String()
			logfile, dbfile, err := paths(project)
			if err != nil {
				return err
			}
			fmt.Println("Database path:", dbfile)
			fmt.Println("Log path:     ", logfile)
			return nil
		},
	}
	return &pathsCmd{cmd: cmd}
}
