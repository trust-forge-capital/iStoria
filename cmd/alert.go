package cmd

import "github.com/spf13/cobra"

var alertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Manage local alerts",
}

var alertListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List alerts",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	alertCmd.AddCommand(alertListCmd)
	rootCmd.AddCommand(alertCmd)
}
