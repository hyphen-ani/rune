package cli

import (
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var tokenRevokeCmd = &cobra.Command{
	Use:   "revoke [token-id]",
	Short: "Revoke an existing token, immediately disabling its access",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tokenID := args[0]
		token, err := config.LoadToken()
		if err != nil {
			Error("[AUTHENTICATION FAILED] Please login first using: rune login <token>")
			return
		}
		c := client.New("http://localhost:8080", token)
		err = c.RevokeToken(tokenID)
		if err != nil {
			Error(err.Error())
			return
		}
		Success("Token revoked")
	},
}

func init() {
	tokenCmd.AddCommand(tokenRevokeCmd)
}
