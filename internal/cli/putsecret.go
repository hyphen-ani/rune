package cli

import (
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{

	Use:   "put [key] [value]",
	Short: "Store and encrypt a secret value under a specified key in the vault",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New("http://localhost:8080")
		err := c.Put(args[0], args[1])
		if err != nil {
			Error(err.Error())
			return
		}
		Success("Secret Stored Successfully")
	},
}

func init() {

	rootCmd.AddCommand(putCmd)

}
