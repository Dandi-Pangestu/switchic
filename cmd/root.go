// Package cmd hosts the Cobra command tree for switchic.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags. The default value is what users
// will see when building from source without explicit version info.
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "switchic",
	Short: "Plug-and-play AI coding platform switcher with workspace + cost controls",
	Long: `switchic keeps your project ready to work with any supported AI coding
platform. The MVP focuses on Claude. Initialize a repo or workspace, choose
which agents / skills / rules should be active, and generate platform-specific
context files with a single command.

Common flows:
  switchic init                     # one-time per repo
  switchic switch claude            # regenerate Claude context
  switchic status                   # see active components
  switchic workspace init           # multi-repo support
  switchic add agent reviewer       # turn a component on
  switchic remove skill summarize   # turn a component off
  switchic run                      # print the resolved workflow plan
`,
	SilenceUsage: true,
}

// Execute runs the root command. main() is a thin wrapper around this.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
