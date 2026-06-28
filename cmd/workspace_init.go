package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Dandi-Pangestu/switchic/internal/output"
	"github.com/Dandi-Pangestu/switchic/internal/project"
	"github.com/Dandi-Pangestu/switchic/internal/util"
)

var (
	workspaceInitName  string
	workspaceInitNotes string
)

var workspaceInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create switchic.workspace.yaml in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd := util.Cwd()
		name := workspaceInitName
		if name == "" {
			name = filepath.Base(cwd)
		}
		m, err := project.InitWorkspace(cwd, name, workspaceInitNotes)
		if err != nil {
			if errors.Is(err, util.ErrAlreadyExists) {
				return fmt.Errorf("switchic.workspace.yaml already exists in %s", cwd)
			}
			return err
		}
		output.Info(cmd.OutOrStdout(), "Initialized workspace %q in %s", m.Name, cwd)
		output.Info(cmd.OutOrStdout(), "Next: `switchic workspace add <repo-path>` for each repo.")
		return nil
	},
}

func init() {
	workspaceInitCmd.Flags().StringVar(&workspaceInitName, "name", "", "workspace name (defaults to directory name)")
	workspaceInitCmd.Flags().StringVar(&workspaceInitNotes, "notes", "", "short description of this workspace")
	workspaceCmd.AddCommand(workspaceInitCmd)
}
