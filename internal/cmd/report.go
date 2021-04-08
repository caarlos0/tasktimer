package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/caarlos0/tasktimer/internal/ui"
	"github.com/charmbracelet/glamour"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

type reportCmd struct {
	cmd *cobra.Command
}

func newRerportCmd() *reportCmd {
	cmd := &cobra.Command{
		Use:     "report",
		Aliases: []string{"r"},
		Short:   "Print a markdown report of the given project to STDOUT",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			project := cmd.Parent().Flag("project").Value.String()
			db, f, err := setup(project)
			if err != nil {
				return err
			}
			defer db.Close()
			defer f.Close()

			var buf bytes.Buffer
			if err := ui.WriteProjectMarkdown(db, project, &buf); err != nil {
				return err
			}

			md := buf.String()

			if isatty.IsTerminal(os.Stdout.Fd()) {
				rendered, err := glamour.RenderWithEnvironmentConfig(md)
				if err != nil {
					return err
				}
				md = rendered
			}

			fmt.Print(md)
			return nil
		},
	}

	return &reportCmd{cmd: cmd}
}
