package cli

import (
	"rune/internal/config"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{

	Use:   "login [token]",
	Short: "Save authentication token",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]
		err := config.SaveToken(token)
		if err != nil {
			Error("Failed to save token")
			return
		}
		Success("Token saved successfully")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
