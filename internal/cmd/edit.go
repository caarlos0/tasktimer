package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

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
			tmp, err := ioutil.TempFile("", "")
			if err != nil {
				return err
			}
			if err := tmp.Close(); err != nil {
				return err
			}

			if err := newToJSONCmd().cmd.RunE(cmd, []string{tmp.Name()}); err != nil {
				return err
			}

			editor := os.Getenv("EDITOR")
			if editor == "" {
				return fmt.Errorf("no $EDITOR set")
			}

			log.Printf("%s %s\n", editor, tmp.Name())
			edit := exec.Command(editor, tmp.Name())
			edit.Stderr = os.Stderr
			edit.Stdout = os.Stdout
			edit.Stdin = os.Stdin
			if err := edit.Run(); err != nil {
				return err
			}

			return newFromJSONCmd().cmd.RunE(cmd, []string{tmp.Name()})
		},
	}

	return &editCmd{cmd: cmd}
}
