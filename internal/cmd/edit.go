package cmd

import (
	"fmt"
	"log"
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

			editor := os.Getenv("EDITOR")
			if editor == "" {
				return fmt.Errorf("no $EDITOR set")
			}

			log.Printf("%s %s\n", editor, tmp)

			editorArgs := []string{tmp}
			if strings.ContainsAny(editor, " ") {
				editorParts := strings.Split(editor, " ")
				editor = editorParts[0]
				editorArgs = append(editorArgs, editorParts[1:]...)
			}

			edit := exec.Command(editor, editorArgs...)
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
