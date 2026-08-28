package cli

import (
	"fmt"
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var getVersion int
var namespaceName string
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
		ns := normalizeNamespace(namespaceName)

		var val string

		if getVersion > 0 {
			val, err = c.GetVersion(args[0], ns, getVersion)
		} else {
			val, err = c.Get(args[0], ns)
		}
		if err != nil {
			Error(err.Error())
			return
		}

		fmt.Println("[SECRET]", val)
		Success("Secret retrieved")
	},
}

func init() {
	getCmd.Flags().IntVar(&getVersion, "version", 0, "Retrieve a specific secret version")
	getCmd.Flags().StringVarP(&namespaceName, "namespace", "n", "default", "Namespace containing the secret")
	rootCmd.AddCommand(getCmd)
}
