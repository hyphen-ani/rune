package cli

import (
	"fmt"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var unsealCmd = &cobra.Command{
	Use:   "unseal",
	Short: "Unseal the vault using a passphrase to restore access to encrypted secrets",
	Run: func(cmd *cobra.Command, args []string) {
		var passphrase string
		fmt.Print("Enter passphrase: ")
		fmt.Scanln(&passphrase)

		c := client.New("http://localhost:8080", "")
		err := c.Unseal(passphrase)

		if err != nil {
			Error(err.Error())
			return
		}
		Success("Vault Unsealed Successfully")
	},
}

func init() {
	rootCmd.AddCommand(unsealCmd)
}
