package cmd

import (
	"fmt"

	"github.com/muesli/coral"
)

type pathsCmd struct {
	cmd *coral.Command
}

func newPathsCmd() *pathsCmd {
	cmd := &coral.Command{
		Use:   "paths",
		Short: "Print the paths being used for logs, data et al",
		Args:  coral.NoArgs,
		RunE: func(cmd *coral.Command, args []string) error {
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
