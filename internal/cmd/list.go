package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/caarlos0/tasktimer/internal/store"
	"github.com/caarlos0/tasktimer/internal/ui"
	"github.com/charmbracelet/glamour"
	"github.com/mattn/go-isatty"
	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/cobra"
)

type listCmd struct {
	cmd     *cobra.Command
	verbose bool
}

func newListCmd() *listCmd {
	l := &listCmd{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all projects",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			home := gap.NewScope(gap.User, "tasktimer")
			datas, err := home.DataDirs()
			if err != nil {
				return err
			}

			var buf bytes.Buffer
			for _, data := range datas {
				if _, err := os.Stat(data); err != nil && os.IsNotExist(err) {
					continue
				}
				if err := filepath.Walk(data, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if filepath.Ext(path) == ".db" {
						name := strings.Replace(filepath.Base(path), ".db", "", 1)
						if l.verbose {
							db, f, err := setup(name)
							if err != nil {
								_, _ = fmt.Fprintf(&buf, "- %s (error: %v)\n", name, err)
								return filepath.SkipDir
							}
							tasks, err := store.GetTaskList(db)
							db.Close()
							f.Close()
							if err != nil {
								_, _ = fmt.Fprintf(&buf, "- %s (error: %v)\n", name, err)
								return filepath.SkipDir
							}
							total := ui.SumTasksTimes(tasks, time.Time{})
							h := int(total.Hours())
							m := int(total.Minutes()) % 60
							_, _ = fmt.Fprintf(&buf, "- %s (%d tasks, %dh %dm)\n", name, len(tasks), h, m)
						} else {
							_, _ = fmt.Fprintln(&buf, "- "+name)
						}
						return filepath.SkipDir
					}
					return nil
				}); err != nil {
					return err
				}
			}

			if isatty.IsTerminal(os.Stdout.Fd()) {
				rendered, err := glamour.RenderWithEnvironmentConfig(buf.String())
				if err != nil {
					return err
				}
				fmt.Print(rendered)
				return nil
			}

			fmt.Print(buf.String())
			return nil
		},
	}

	cmd.Flags().BoolVarP(&l.verbose, "verbose", "v", false, "Show task count and total time per project")
	l.cmd = cmd
	return l
}
