package cli

import "github.com/spf13/cobra"

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage authentication tokens",
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
