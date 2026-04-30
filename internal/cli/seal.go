package cli

import (
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var sealCmd = &cobra.Command{

	Use:   "seal",
	Short: "Seal the vault by clearing the encryption key from memory, preventing all access to secrets",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New("http://localhost:8080")
		err := c.Seal()
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
