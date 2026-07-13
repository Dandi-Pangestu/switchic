package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/Dandi-Pangestu/switchic/internal/output"
	"github.com/Dandi-Pangestu/switchic/internal/util"
	"github.com/Dandi-Pangestu/switchic/internal/workspace"
)

var workspaceLinkCmd = &cobra.Command{
	Use:   "link",
	Short: "Create/refresh repos/ symlinks for out-of-tree repos",
	RunE: func(cmd *cobra.Command, args []string) error {
		root := util.FindWorkspaceRoot(util.Cwd())
		if root == "" {
			return errors.New("no switchic.workspace.yaml found")
		}
		m, err := workspace.Load(util.WorkspacePath(root))
		if err != nil {
			return err
		}

		errs := workspace.LinkAll(root, m)
		if err := workspace.EnsureGitignoreEntry(root); err != nil {
			output.Info(cmd.OutOrStdout(), "Warning: could not update .gitignore: %v", err)
		}

		linked := 0
		for _, r := range m.Repos {
			if _, failed := errs[r.Name]; failed {
				continue
			}
			linked++
		}
		output.Info(cmd.OutOrStdout(), "Linked %d repo(s) under %s/", linked, workspace.LinkDir)

		if len(errs) > 0 {
			output.Info(cmd.OutOrStdout(), "")
			output.Info(cmd.OutOrStdout(), "Could not link:")
			for _, r := range m.Repos {
				if err, failed := errs[r.Name]; failed {
					output.Bullet(cmd.OutOrStdout(), "%s: %v", r.Name, err)
				}
			}
		}
		return nil
	},
}

func init() { workspaceCmd.AddCommand(workspaceLinkCmd) }
