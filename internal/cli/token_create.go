package cli

import (
	"fmt"
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var tokenName string
var tokenNamespace string

var tokenCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new access token with an optional name for identification",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token, err := config.LoadToken()
		if err != nil {
			Error("[AUTHENTICATION FAILED] Please login first using: rune login <token>")
			return
		}

		name := tokenName

		if len(args) > 0 {
			name = args[0]
		}

		if name == "" {
			name = "default"
		}

		c := client.New("http://localhost:8080", token)
		newToken, record, err := c.CreateToken(name, tokenNamespace)
		if err != nil {
			Error(err.Error())
			return
		}

		Success("Token created successfully")

		fmt.Println()
		fmt.Println("  ID:   ", record.ID)
		fmt.Println("  Name: ", record.Name)
		fmt.Println("  Created At: ", record.CreatedAt)
		fmt.Println("  Token:", newToken)
		fmt.Println()
	},
}

func init() {
	tokenCreateCmd.Flags().StringVar(&tokenName, "name", "", "Name of the token")
	tokenCreateCmd.Flags().StringVarP(&tokenNamespace, "namespace", "n", "default", "Namespace of the token is allowed to access")
	tokenCmd.AddCommand(tokenCreateCmd)
}
