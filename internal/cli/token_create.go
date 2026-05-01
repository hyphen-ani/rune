package cli

import (
	"fmt"
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var tokenName string

var tokenCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new access token with an optional name for identification",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := config.LoadToken()
		if err != nil {
			Error("[AUTHENTICATION FAILED] Please login first using: rune login <token>")
			return
		}
		c := client.New("http://localhost:8080", token)
		newToken, record, err := c.CreateToken(tokenName)
		if err != nil {
			Error(err.Error())
			return
		}

		fmt.Println()
		fmt.Println("  ID:   ", record.ID)
		fmt.Println("  Name: ", record.Name)
		fmt.Println("  Token:", newToken)
		fmt.Println()
	},
}

func init() {
	tokenCreateCmd.Flags().StringVar(&tokenName, "name", "default", "Name of the token")
	tokenCmd.AddCommand(tokenCreateCmd)
}
