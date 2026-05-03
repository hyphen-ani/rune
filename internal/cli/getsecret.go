package cli

import (
	"fmt"
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var namespace string
var getCmd = &cobra.Command{

	Use:   "get [key]",
	Short: "Retrieve and decrypt a stored secret by key from the vault (requires vault to be unsealed)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token, err := config.LoadToken()
		if err != nil {
			Error("[AUTHENTICATION FAILED] Please login first using: rune login <token>")
			return
		}
		c := client.New("http://localhost:8080", token)
		ns := normalizeNamespace(namespace)
		val, err := c.Get(args[0], ns)
		if err != nil {
			Error(err.Error())
			return
		}
		fmt.Println("[SECRET]", val)
		Success("Secret retrieved")
	},
}

func init() {
	getCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "The vault namespace to use")
	rootCmd.AddCommand(getCmd)
}
