package cli

import (
	"fmt"
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var namespaceDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a namespace",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		token, err := config.LoadToken()
		if err != nil {
			Error("[AUTHENTICATION FAILED] Please login first using: rune login [token]")
			return
		}

		if name == "default" {
			Error("Cannot delete default namespace")
			return
		}

		fmt.Printf("Are you sure you want to delete namespace '%s'? (y/N): ", name)
		var confirm string
		fmt.Scanln(&confirm)

		if confirm != "y" && confirm != "Y" {
			Info("Operation Cancelled")
			return
		}

		c := client.New("http://localhost:8080", token)

		err = c.DeleteNamespace(name)
		if err != nil {
			Error(err.Error())
			return
		}

		Success("Namespace deleted")
	},
}

func init() {
	namespaceCmd.AddCommand(namespaceDeleteCmd)
}
