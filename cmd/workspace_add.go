package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Dandi-Pangestu/switchic/internal/output"
	"github.com/Dandi-Pangestu/switchic/internal/util"
	"github.com/Dandi-Pangestu/switchic/internal/workspace"
)

var (
	workspaceAddRole        string
	workspaceAddNotes       string
	workspaceAddContextFile string
)

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
		if err := m.AddRepo(args[0], workspaceAddRole, workspaceAddNotes, workspaceAddContextFile); err != nil {
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

		newRepo := m.Repos[len(m.Repos)-1]
		if err := workspace.LinkRepo(root, newRepo); err != nil {
			output.Info(cmd.OutOrStdout(), "Warning: could not create repos/ symlink: %v", err)
		} else if err := workspace.EnsureGitignoreEntry(root); err != nil {
			output.Info(cmd.OutOrStdout(), "Warning: could not update .gitignore: %v", err)
		}
		return nil
	},
}

func init() {
	workspaceAddCmd.Flags().StringVar(&workspaceAddRole, "role", "", "role label (e.g. backend, frontend, contracts)")
	workspaceAddCmd.Flags().StringVar(&workspaceAddNotes, "notes", "", "short description of this repo")
	workspaceAddCmd.Flags().StringVar(&workspaceAddContextFile, "context-file", "", "path to this repo's context file, relative to workspace root (overrides platform default, e.g. ../hub_core/docs/CLAUDE.md)")
	workspaceCmd.AddCommand(workspaceAddCmd)
}
