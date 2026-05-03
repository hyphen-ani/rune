package cli

import (
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{

	Use:   "put [key] [value]",
	Short: "Store and encrypt a secret value under a specified key in the vault",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		token, err := config.LoadToken()
		if err != nil {
			Error("[AUTHENTICATION FAILED] Please login first using: rune login <token>")
			return
		}
		c := client.New("http://localhost:8080", token)
		ns := normalizeNamespace(namespace)
		err = c.Put(args[0], args[1], ns)
		if err != nil {
			Error(err.Error())
			return
		}
		Success("Secret Stored Successfully")
	},
}

func init() {
	putCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Namespace of the secret")
	rootCmd.AddCommand(putCmd)

}
