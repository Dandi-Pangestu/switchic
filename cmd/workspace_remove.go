package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Dandi-Pangestu/switchic/internal/output"
	"github.com/Dandi-Pangestu/switchic/internal/util"
	"github.com/Dandi-Pangestu/switchic/internal/workspace"
)

var workspaceRemoveCmd = &cobra.Command{
	Use:   "remove <repo-name>",
	Short: "Remove a repo from the workspace manifest",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root := util.FindWorkspaceRoot(util.Cwd())
		if root == "" {
			return errors.New("no switchic.workspace.yaml found")
		}
		m, err := workspace.Load(util.WorkspacePath(root))
		if err != nil {
			return err
		}
		if err := m.RemoveRepo(args[0]); err != nil {
			if errors.Is(err, util.ErrNotFound) {
				return fmt.Errorf("no repo named %q in workspace %q", args[0], m.Name)
			}
			return err
		}
		if err := workspace.Save(util.WorkspacePath(root), m); err != nil {
			return err
		}
		output.Info(cmd.OutOrStdout(), "Removed repo %q from workspace %q", args[0], m.Name)
		return nil
	},
}

func init() { workspaceCmd.AddCommand(workspaceRemoveCmd) }
