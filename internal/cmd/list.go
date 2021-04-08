package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/mattn/go-isatty"
	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/cobra"
)

type listCmd struct {
	cmd *cobra.Command
}

func newListCmd() *listCmd {
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
						_, _ = fmt.Fprintln(&buf, "- "+strings.Replace(filepath.Base(path), ".db", "", 1))
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

	return &listCmd{cmd: cmd}
}
