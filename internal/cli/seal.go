package cli

import (
	"fmt"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var sealCmd = &cobra.Command{

	Use:   "seal",
	Short: "Seal the vault",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New("http://localhost:8080")
		err := c.Seal()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Vault Sealed")
	},
}

func init() {

	rootCmd.AddCommand(sealCmd)

}
