package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/caarlos0/tasktimer/internal/model"
	"github.com/caarlos0/tasktimer/internal/store"
	"github.com/spf13/cobra"
)

type fromJSONCmd struct {
	cmd *cobra.Command
}

func newFromJSONCmd() *fromJSONCmd {
	cmd := &cobra.Command{
		Use:   "from-json",
		Short: "Imports a JSON into a project - WARNING: it will wipe the project first, use with care!",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project := cmd.Parent().Flag("project").Value.String()
			db, f, err := setup(project)
			if err != nil {
				return err
			}
			defer db.Close()
			defer f.Close()

			input, err := ioutil.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", args[0], err)
			}

			var tasks []model.ExportedTask
			if err := json.Unmarshal(input, &tasks); err != nil {
				return fmt.Errorf("input json is not in the correct format: %w", err)
			}

			tmp, err := ioutil.TempFile("", "tasktimer-"+project)
			if err != nil {
				return fmt.Errorf("failed to create backup file: %w", err)
			}
			if _, err := db.Backup(tmp, 0); err != nil {
				return fmt.Errorf("failed to backup to %s: %w", tmp.Name(), err)
			}

			log.Printf("backup made to %s\n", tmp.Name())

			if err := db.DropAll(); err != nil {
				return fmt.Errorf("failed to clear database: %w", err)
			}

			return store.LoadTasks(db, tasks)
		},
	}

	return &fromJSONCmd{cmd: cmd}
}
