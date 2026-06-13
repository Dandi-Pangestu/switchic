package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Dandi-Pangestu/switchic/internal/output"
	"github.com/Dandi-Pangestu/switchic/internal/util"
	"github.com/Dandi-Pangestu/switchic/internal/workspace"
)

var workspaceAddRole string

var workspaceAddCmd = &cobra.Command{
	Use:   "add <repo-path>",
	Short: "Register a repo in the workspace manifest",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root := util.FindWorkspaceRoot(util.Cwd())
		if root == "" {
			return errors.New("no switchic.workspace.yaml found — run `switchic workspace init` first")
		}
		m, err := workspace.Load(util.WorkspacePath(root))
		if err != nil {
			return err
		}
		if err := m.AddRepo(args[0], workspaceAddRole); err != nil {
			if errors.Is(err, util.ErrAlreadyExists) {
				return fmt.Errorf("a repo with that name is already in the workspace")
			}
			return err
		}
		if err := workspace.Save(util.WorkspacePath(root), m); err != nil {
			return err
		}
		output.Info(cmd.OutOrStdout(), "Added repo %q to workspace %q", args[0], m.Name)
		if miss := workspace.MissingRepos(root, m); len(miss) > 0 {
			output.Info(cmd.OutOrStdout(), "Warning: %d repo path(s) do not exist on disk yet.", len(miss))
		}
		return nil
	},
}

func init() {
	workspaceAddCmd.Flags().StringVar(&workspaceAddRole, "role", "", "role label (e.g. backend, frontend, contracts)")
	workspaceCmd.AddCommand(workspaceAddCmd)
}
