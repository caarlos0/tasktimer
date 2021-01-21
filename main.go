package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/tasktimer/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/dgraph-io/badger/v3"
	"github.com/mattn/go-isatty"
	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/cobra"
)

var project string

func main() {
	rootCmd.PersistentFlags().StringVarP(&project, "project", "p", "default", "Project name")

	rootCmd.AddCommand(reportCmd, completionsCmd, pathsCmd)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "tt",
	Short: "Task Timer (tt) is a dead simple TUI task timer",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, f, err := setup()
		if err != nil {
			return err
		}
		defer db.Close()
		defer f.Close()

		var p = tea.NewProgram(ui.Init(db, project))
		p.EnterAltScreen()
		defer p.ExitAltScreen()
		return p.Start()
	},
}

var reportCmd = &cobra.Command{
	Use:     "report",
	Aliases: []string{"r"},
	Short:   "Print a markdown report of the given project to STDOUT",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, f, err := setup()
		if err != nil {
			return err
		}
		defer db.Close()
		defer f.Close()

		var b bytes.Buffer
		if err := ui.WriteProjectMarkdown(db, project, &b); err != nil {
			return err
		}

		var out = b.String()
		if isatty.IsTerminal(os.Stdout.Fd()) {
			outs, err := glamour.Render(out, "auto")
			if err != nil {
				return err
			}
			out = outs
		}

		fmt.Print(out)
		return nil
	},
}

var pathsCmd = &cobra.Command{
	Use:   "paths",
	Short: "Print the paths being used for logs, data et al",
	RunE: func(cmd *cobra.Command, args []string) error {
		logfile, dbfile, err := paths()
		if err != nil {
			return err
		}
		fmt.Println("Database path:", dbfile)
		fmt.Println("Log path:     ", logfile)
		return nil
	},
}

func paths() (string, string, error) {
	var home = gap.NewScope(gap.User, "tasktimer")

	logfile, err := home.LogPath(project + ".log")
	if err != nil {
		return "", "", err
	}

	dbfile, err := home.DataPath(project + ".db")
	if err != nil {
		return "", "", err
	}

	return logfile, dbfile, nil
}

func setup() (*badger.DB, io.Closer, error) {
	logfile, dbfile, err := paths()
	if err != nil {
		return nil, nil, err
	}

	if err := os.MkdirAll(filepath.Dir(logfile), 0754); err != nil {
		return nil, nil, err
	}

	f, err := tea.LogToFile(logfile, "tasktimer")
	if err != nil {
		return nil, nil, err
	}

	// TODO: maybe sync writes?
	var options = badger.DefaultOptions(dbfile).
		WithLogger(badgerStdLoggerAdapter{}).
		WithLoggingLevel(badger.ERROR)
	db, err := badger.Open(options)
	return db, f, err
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
  $ source <(tt completion bash)
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
