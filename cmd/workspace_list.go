package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/Dandi-Pangestu/switchic/internal/output"
	"github.com/Dandi-Pangestu/switchic/internal/util"
	"github.com/Dandi-Pangestu/switchic/internal/workspace"
)

var workspaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repos in the workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		root := util.FindWorkspaceRoot(util.Cwd())
		if root == "" {
			return errors.New("no switchic.workspace.yaml found")
		}
		m, err := workspace.Load(util.WorkspacePath(root))
		if err != nil {
			return err
		}
		output.Info(cmd.OutOrStdout(), "Workspace: %s   platform=%s   repos=%d", m.Name, m.Platform, len(m.Repos))
		t := output.Table{}
		for _, r := range m.Repos {
			role := r.Role
			if role == "" {
				role = "-"
			}
			t.Row(r.Name, r.Path+"  ["+role+"]")
		}
		t.Render(cmd.OutOrStdout())

		if miss := workspace.MissingRepos(root, m); len(miss) > 0 {
			output.Info(cmd.OutOrStdout(), "")
			output.Info(cmd.OutOrStdout(), "Missing on disk:")
			for _, r := range miss {
				output.Bullet(cmd.OutOrStdout(), "%s -> %s", r.Name, r.Path)
			}
		}
		return nil
	},
}

func init() { workspaceCmd.AddCommand(workspaceListCmd) }
