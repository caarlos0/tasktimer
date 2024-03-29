package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type editCmd struct {
	cmd *cobra.Command
}

func newEditCmd() *editCmd {
	cmd := &cobra.Command{
		Use:     "edit",
		Short:   "Syntactic sugar for to-json | $EDITOR | from-json",
		Aliases: []string{"e"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			tmp := filepath.Join(os.TempDir(), fmt.Sprintf("tt-%d.json", time.Now().Unix()))

			if err := newToJSONCmd().cmd.RunE(cmd, []string{tmp}); err != nil {
				return err
			}

			editor := strings.Fields(os.Getenv("EDITOR"))
			if len(editor) == 0 {
				return fmt.Errorf("no $EDITOR set")
			}

			editorCmd := editor[0]
			var editorArgs []string
			if len(editor) > 1 {
				editorArgs = append(editorArgs, editor[1:]...)
			}
			editorArgs = append(editorArgs, tmp)

			edit := exec.Command(editorCmd, editorArgs...)
			edit.Stderr = os.Stderr
			edit.Stdout = os.Stdout
			edit.Stdin = os.Stdin
			if err := edit.Run(); err != nil {
				return err
			}

			return newFromJSONCmd().cmd.RunE(cmd, []string{tmp})
		},
	}

	return &editCmd{cmd: cmd}
}
