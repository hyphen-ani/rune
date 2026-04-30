package cli

import (
	"fmt"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var unsealCmd = &cobra.Command{
	Use:   "unseal",
	Short: "Unseals the vault",
	Run: func(cmd *cobra.Command, args []string) {
		var passphrase string
		fmt.Print("Enter passphrase: ")
		fmt.Scanln(&passphrase)

		c := client.New("http://localhost:8080")
		err := c.Unseal(passphrase)

		if err != nil {
			fmt.Println("Failed:", err)
			return
		}
		fmt.Println("Vault Unsealed")
	},
}

func init() {
	rootCmd.AddCommand(unsealCmd)
}
