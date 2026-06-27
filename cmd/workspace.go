package cmd

import (
	"github.com/spf13/cobra"
)

// workspaceCmd is the parent for `switchic workspace ...` subcommands.
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage a multi-repo workspace",
}

func init() { rootCmd.AddCommand(workspaceCmd) }
