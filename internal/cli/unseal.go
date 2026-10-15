package cli

import (
	"fmt"
	"os"
	"rune/pkg/client"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var unsealCmd = &cobra.Command{
	Use:   "unseal",
	Short: "Unseal the vault using a passphrase to restore access to encrypted secrets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter passphrase: ")
		passphraseBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			Error("Failed to read passphrase: " + err.Error())
			return
		}
		passphrase := string(passphraseBytes)

		c := client.New(serverURL(), "")
		err = c.Unseal(passphrase)

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
