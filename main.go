package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/tasktimer/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dgraph-io/badger/v3"
	"github.com/spf13/cobra"
)

var (
	project string
	output  string
)

func main() {
	rootCmd.PersistentFlags().StringVarP(&project, "project", "p", "default", "Project name")
	reportCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Where to save the markdown report file (default \"{project}.md\")")

	rootCmd.AddCommand(reportCmd, completionsCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "tt",
	Short: "Task Timer (tt) is a dead simple TUI task timer",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := setup()
		if err != nil {
			return err
		}
		defer db.Close()

		var p = tea.NewProgram(ui.Init(db, project))
		p.EnterAltScreen()
		defer p.ExitAltScreen()
		return p.Start()
	},
}

var reportCmd = &cobra.Command{
	Use:     "report",
	Aliases: []string{"r"},
	Short:   "Print a markdown report of the given project",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := setup()
		if err != nil {
			return err
		}
		defer db.Close()

		if output == "" {
			output = project + ".md"
		}

		f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		fmt.Println("writing project to", output)
		return ui.WriteProjectMarkdown(db, project, f)
	},
}

func setup() (*badger.DB, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	var folder = filepath.Join(home, "tasktimer")
	if err := os.MkdirAll(folder, 0764); err != nil {
		return nil, err
	}

	log.SetFlags(0)

	var logfile = filepath.Join(folder, project+".log")
	log.Println("logging to", logfile)

	f, err := tea.LogToFile(logfile, "")
	defer f.Close()

	var options = badger.DefaultOptions(filepath.Join(folder, project+".db")).
		WithLogger(badgerStdLoggerAdapter{}).
		WithLoggingLevel(badger.ERROR)
	return badger.Open(options)
}

type badgerStdLoggerAdapter struct{}

func (b badgerStdLoggerAdapter) Errorf(s string, i ...interface{}) {
	log.Printf("[ERR] "+s, i...)
}

func (b badgerStdLoggerAdapter) Warningf(s string, i ...interface{}) {
	log.Printf("[WARN] "+s, i...)
}

func (b badgerStdLoggerAdapter) Infof(s string, i ...interface{}) {
	log.Printf("[INFO] "+s, i...)
}

func (b badgerStdLoggerAdapter) Debugf(s string, i ...interface{}) {
	log.Printf("[DEBUG] "+s, i...)
}

var completionsCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Print shell autocompletion scripts for tt",
	Long: `To load completions:
Bash:
$ source <(goreleaser completion bash)
# To load completions for each session, execute once:
Linux:
  $ tt completion bash > /etc/bash_completion.d/tt
MacOS:
  $ tt completion bash > /usr/local/etc/bash_completion.d/tt
Zsh:
  # If shell completion is not already enabled in your environment you will need
  # to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  # To load completions for each session, execute once:
  $ tt completion zsh > "${fpath[1]}/_tt"
  # You will need to start a new shell for this setup to take effect.
Fish:
  $ tt completion fish | source
  # To load completions for each session, execute once:
  $ tt completion fish > ~/.config/fish/completions/tt.fish
`,
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish"},
	Args:                  cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			err = cmd.Root().GenZshCompletion(cmd.OutOrStdout())
		case "fish":
			err = cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
		}
		return err
	},
}
