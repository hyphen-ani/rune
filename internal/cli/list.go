package cli

import (
	"fmt"
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stored keys",
	Run: func(cmd *cobra.Command, args []string) {

		token, err := config.LoadToken()
		if err != nil {
			Error("[AUTHENTICATION FAILED] Please login first using: rune login <token>")
			return
		}

		c := client.New("http://localhost:8080", token)
		keys, err := c.List()
		if err != nil {
			Error(err.Error())
			return
		}

		if len(keys) == 0 {
			Info("No stored keys found")
			return
		}

		for _, key := range keys {
			fmt.Println("[SECRET]", key)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
