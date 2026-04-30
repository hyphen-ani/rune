package cli

import (
	"fmt"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{

	Use:   "status",
	Short: "Check whether the vault is currently sealed or unsealed",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New("http://localhost:8080")
		status, err := c.Status()
		if err != nil {
			Error(err.Error())
			return
		}

		if status == "sealed" {
			fmt.Println("🔒 Vault is Sealed")
		} else {
			fmt.Println("🔓 Vault is Unsealed")
		}
	},
}

func init() {

	rootCmd.AddCommand(statusCmd)

}
