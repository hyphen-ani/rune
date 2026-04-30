package cli

import (
	"fmt"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{

	Use:   "get [key]",
	Short: "Retrieve and decrypt a stored secret by key from the vault (requires vault to be unsealed)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New("http://localhost:8080")
		val, err := c.Get(args[0])
		if err != nil {
			Error(err.Error())
			return
		}
		fmt.Println("[SECRET]", val)
		Success("Secret retrieved")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
