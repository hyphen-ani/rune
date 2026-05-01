package cli

import (
	"fmt"
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var tokenListCmd = &cobra.Command{

	Use:   "list",
	Short: "List all tokens",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := config.LoadToken()
		if err != nil {
			Error("[AUTHENTICATION FAILED] Please login first using: rune login <token>")
			return
		}
		c := client.New("http://localhost:8080", token)
		tokens, err := c.ListTokens()
		if err != nil {
			Error(err.Error())
			return
		}
		if len(tokens) == 0 {
			Info("No tokens found")
			return
		}
		for _, t := range tokens {
			status := "active"
			if t.Revoked {
				status = "revoked"
			}
			fmt.Printf("• %s  %s  %s\n", t.ID, t.Name, status)
		}
	},
}

func init() {
	tokenCmd.AddCommand(tokenListCmd)
}
