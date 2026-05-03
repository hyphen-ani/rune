package cli

import (
	"fmt"
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var tokenListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tokens with their status (active or revoked)",

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

		// 🧾 Header
		fmt.Println("\nTokens:\n")
		fmt.Printf("  %-28s %-20s %-25s %-10s\n", "ID", "NAME", "CREATED AT", "STATUS")
		fmt.Println("  --------------------------------------------------------------------------------------")

		for _, t := range tokens {

			status := "active"
			color := "\033[32m"

			if t.Revoked {
				status = "revoked"
				color = "\033[31m"
			}

			reset := "\033[0m"

			fmt.Printf(
				"  %-28s %-20s %-25s %s%-10s%s\n",
				t.ID,
				t.Name,
				t.CreatedAt,
				color,
				status,
				reset,
			)
		}

		fmt.Println()
	},
}

func init() {
	tokenCmd.AddCommand(tokenListCmd)
}
