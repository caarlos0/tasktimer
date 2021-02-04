package cmd

import "github.com/spf13/cobra"

type completionCmd struct {
	cmd *cobra.Command
}

func newCompletionCmd() *completionCmd {
	cmd := &cobra.Command{
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
		Hidden:                true,
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
	return &completionCmd{cmd: cmd}
}
