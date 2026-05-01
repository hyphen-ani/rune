package cli

import (
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var sealCmd = &cobra.Command{

	Use:   "seal",
	Short: "Seal the vault by clearing the encryption key from memory, preventing all access to secrets",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := config.LoadToken()
		if err != nil {
			Error("[AUTHENTICATION FAILED] Please login first using: rune login <token>")
			return
		}
		c := client.New("http://localhost:8080", token)
		err = c.Seal()
		if err != nil {
			Error(err.Error())
			return
		}
		Success("Vault Sealed Successfully")
	},
}

func init() {

	rootCmd.AddCommand(sealCmd)

}
