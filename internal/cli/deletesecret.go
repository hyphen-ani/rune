package cli

import (
	"fmt"
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [key]",
	Short: "Delete a secret",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		key := args[0]
		token, err := config.LoadToken()
		if err != nil {
			Error("[AUTHENTICATION FAILED] Please login first using: rune login [token]")
			return
		}

		ns := normalizeNamespace(namespace)
		fmt.Printf("Are you sure you want to delete secret '%s' (namespace: %s)? (y/N): ", key, ns)

		var confirm string
		fmt.Scanln(&confirm)

		if confirm != "y" && confirm != "Y" {
			Info("Operation Cancelled")
			return
		}

		c := client.New("http://localhost:8080", token)

		err = c.Delete(key, ns)
		if err != nil {
			Error(err.Error())
			return
		}

		Success("Secret deleted")
	},
}

func init() {
	deleteCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace")
	rootCmd.AddCommand(deleteCmd)
}
