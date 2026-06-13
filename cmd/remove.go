package cmd

import "github.com/spf13/cobra"

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Disable an agent, skill, or rule",
}

var removeAgentCmd = &cobra.Command{
	Use:   "agent <name>",
	Short: "Disable an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMutation(cmd, kindAgent, args[0], false)
	},
}

var removeSkillCmd = &cobra.Command{
	Use:   "skill <name>",
	Short: "Disable a skill",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMutation(cmd, kindSkill, args[0], false)
	},
}

var removeRuleCmd = &cobra.Command{
	Use:   "rule <name>",
	Short: "Disable a rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMutation(cmd, kindRule, args[0], false)
	},
}

func init() {
	removeCmd.AddCommand(removeAgentCmd, removeSkillCmd, removeRuleCmd)
	rootCmd.AddCommand(removeCmd)
}
