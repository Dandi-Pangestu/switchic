package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Dandi-Pangestu/switchic/internal/app"
	"github.com/Dandi-Pangestu/switchic/internal/config"
	"github.com/Dandi-Pangestu/switchic/internal/output"
	"github.com/Dandi-Pangestu/switchic/internal/platform"
	"github.com/Dandi-Pangestu/switchic/internal/util"
	"github.com/Dandi-Pangestu/switchic/internal/workspace"
)

var replaceFlag bool

var switchCmd = &cobra.Command{
	Use:   "switch [platform]",
	Short: "Switch the active platform and regenerate its config files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		adapter, err := platform.Get(name)
		if err != nil {
			return err
		}

		ctx, err := app.LoadContext(util.Cwd())
		if err != nil {
			return err
		}
		if !ctx.IsWorkspace && !ctx.HasProject() {
			return fmt.Errorf("no .switchic/config.yaml or workspace manifest found — run `switchic init` first")
		}

		// Persist the platform change in whichever config drives this run.
		if ctx.IsWorkspace {
			ctx.Workspace.Platform = name
			if err := workspace.Save(util.WorkspacePath(ctx.WorkspaceRoot), ctx.Workspace); err != nil {
				return err
			}
		}
		if ctx.HasProject() {
			ctx.Project.Platform = name
			if err := config.Save(util.ProjectConfigPath(ctx.WorkingDir), ctx.Project); err != nil {
				return err
			}
		}

		resolved, err := app.Resolve(ctx.PrimaryRoot(), ctx.Project, ctx.Workspace, ctx.IsWorkspace)
		if err != nil {
			return err
		}

		pCtx := resolved.ToContext()
		pCtx.Replace = replaceFlag
		written, err := adapter.Generate(pCtx)
		if err != nil {
			return err
		}

		output.Info(cmd.OutOrStdout(), "Switched to %s. Generated %d file(s):", name, len(written))
		for _, p := range written {
			output.Bullet(cmd.OutOrStdout(), "%s", p)
		}
		return nil
	},
}

func init() {
	switchCmd.Flags().BoolVar(&replaceFlag, "replace", false, "overwrite existing files even if user-written (no switchic banner)")
	rootCmd.AddCommand(switchCmd)
}
