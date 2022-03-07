package cmd

import (
	"github.com/caarlos0/tasktimer/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/coral"
)

func Execute(version string, exit func(int), args []string) {
	newRootCmd(version, exit).Execute(args)
}

type rootCmd struct {
	cmd     *coral.Command
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
	root := &rootCmd{
		exit: exit,
	}
	cmd := &coral.Command{
		Use:          "tt",
		Short:        "Task Timer (tt) is a dead simple TUI task timer",
		Version:      version,
		SilenceUsage: true,
		RunE: func(cmd *coral.Command, args []string) error {
			db, f, err := setup(root.project)
			if err != nil {
				return err
			}
			defer db.Close()
			defer f.Close()

			p := tea.NewProgram(ui.Init(db, root.project))
			p.EnterAltScreen()
			defer p.ExitAltScreen()
			return p.Start()
		},
	}

	cmd.PersistentFlags().StringVarP(&root.project, "project", "p", "default", "Project name")

	cmd.AddCommand(
		newRerportCmd().cmd,
		newCompletionCmd().cmd,
		newPathsCmd().cmd,
		newToJSONCmd().cmd,
		newFromJSONCmd().cmd,
		newListCmd().cmd,
		newEditCmd().cmd,
		newManCmd().cmd,
	)

	root.cmd = cmd
	return root
}
