package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// workspaceCmd is the parent for `switchic workspace ...` subcommands.
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage a multi-repo workspace",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stderr, "workspace: this feature is still under development and not yet available")
		os.Exit(1)
	},
}

func init() { rootCmd.AddCommand(workspaceCmd) }
