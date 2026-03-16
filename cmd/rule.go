package cmd

import "github.com/spf13/cobra"

var ruleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Manage local rules",
}

var ruleListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List rules",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	ruleCmd.AddCommand(ruleListCmd)
	rootCmd.AddCommand(ruleCmd)
}
