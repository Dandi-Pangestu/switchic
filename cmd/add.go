package cmd

import "github.com/spf13/cobra"

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Enable an agent, skill, or rule",
}

var addAgentCmd = &cobra.Command{
	Use:   "agent <name>",
	Short: "Enable an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMutation(cmd, kindAgent, args[0], true)
	},
}

var addSkillCmd = &cobra.Command{
	Use:   "skill <name>",
	Short: "Enable a skill",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMutation(cmd, kindSkill, args[0], true)
	},
}

var addRuleCmd = &cobra.Command{
	Use:   "rule <name>",
	Short: "Enable a rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMutation(cmd, kindRule, args[0], true)
	},
}

func init() {
	addCmd.AddCommand(addAgentCmd, addSkillCmd, addRuleCmd)
	rootCmd.AddCommand(addCmd)
}
