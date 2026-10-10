package cli

import "github.com/spf13/cobra"

var systemTokenCmd = &cobra.Command{
	Use:   "system-token",
	Short: "Manage system (read-only, SDK) tokens",
}

func init() {
	rootCmd.AddCommand(systemTokenCmd)
}
