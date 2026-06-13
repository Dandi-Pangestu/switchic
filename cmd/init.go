package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Dandi-Pangestu/switchic/internal/output"
	"github.com/Dandi-Pangestu/switchic/internal/project"
	"github.com/Dandi-Pangestu/switchic/internal/util"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize switchic in the current repo (.switchic/config.yaml)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd := util.Cwd()
		p, err := project.Init(cwd)
		if err != nil {
			if errors.Is(err, util.ErrAlreadyExists) {
				return fmt.Errorf(".switchic/config.yaml already exists in %s", cwd)
			}
			return err
		}
		output.Info(cmd.OutOrStdout(), "Initialized switchic project in %s", cwd)
		output.Info(cmd.OutOrStdout(), "  platform : %s", p.Platform)
		output.Info(cmd.OutOrStdout(), "  workflows: %s", strings.Join(p.Workflows.Active, ", "))
		if p.Language != "" {
			output.Info(cmd.OutOrStdout(), "  language : %s (detected)", p.Language)
		}
		output.Info(cmd.OutOrStdout(), "")
		output.Info(cmd.OutOrStdout(), "Next: run `switchic switch claude` to generate CLAUDE.md.")
		return nil
	},
}

func init() { rootCmd.AddCommand(initCmd) }
