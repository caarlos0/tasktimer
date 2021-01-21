package cmd

import (
	"github.com/caarlos0/tasktimer/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func Execute(version string, exit func(int), args []string) {
	newRootCmd(version, exit).Execute(args)
}

type rootCmd struct {
	cmd     *cobra.Command
	project string
	exit    func(int)
}

func (c rootCmd) Execute(args []string) {
	c.cmd.SetArgs(args)
	if err := c.cmd.Execute(); err != nil {
		c.exit(1)
	}
}

func newRootCmd(version string, exit func(int)) *rootCmd {
	var root = &rootCmd{
		exit: exit,
	}
	var cmd = &cobra.Command{
		Use:   "tt",
		Short: "Task Timer (tt) is a dead simple TUI task timer",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, f, err := setup(root.project)
			if err != nil {
				return err
			}
			defer db.Close()
			defer f.Close()

			var p = tea.NewProgram(ui.Init(db, root.project))
			p.EnterAltScreen()
			defer p.ExitAltScreen()
			return p.Start()
		},
	}

	cmd.PersistentFlags().StringVarP(&root.project, "project", "p", "default", "Project name")

	cmd.AddCommand(newRerportCmd().cmd, newCompletionCmd().cmd, newPathsCmd().cmd)

	root.cmd = cmd
	return root
}
