package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current Rune CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rune", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
