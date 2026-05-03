package cli

import (
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var namespaceCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a namespace",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		token, err := config.LoadToken()
		if err != nil {
			Error("[INVALID REQUEST] Please Login First")
			return
		}

		c := client.New("http://localhost:8080", token)

		err = c.CreateNamespace(name)
		if err != nil {
			Error(err.Error())
			return
		}

		Success("[NAMESPACE CREATED] Created namespace")
	},
}

func init() {
	namespaceCmd.AddCommand(namespaceCreateCmd)
}
