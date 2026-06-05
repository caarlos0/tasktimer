package cmd

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/tasktimer/internal/ui"
	"github.com/charmbracelet/glamour"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

type reportCmd struct {
	cmd   *cobra.Command
	since string
	until string
}

func newReportCmd() *reportCmd {
	r := &reportCmd{}
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

			var since, until time.Time
			if r.since != "" {
				since, err = time.Parse("2006-01-02", r.since)
				if err != nil {
					return fmt.Errorf("--since: invalid date %q, use YYYY-MM-DD", r.since)
				}
			}
			if r.until != "" {
				until, err = time.Parse("2006-01-02", r.until)
				if err != nil {
					return fmt.Errorf("--until: invalid date %q, use YYYY-MM-DD", r.until)
				}
				until = until.Add(24 * time.Hour)
			}

			var buf bytes.Buffer
			if err := ui.WriteProjectMarkdown(db, project, &buf, since, until); err != nil {
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

	cmd.Flags().StringVar(&r.since, "since", "", "Show tasks on or after this date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&r.until, "until", "", "Show tasks on or before this date (YYYY-MM-DD)")
	r.cmd = cmd
	return r
}
