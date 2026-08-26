package cli

import (
	"fmt"
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stored keys in the vault (filtered by namespace)",
	Run: func(cmd *cobra.Command, args []string) {

		token, err := config.LoadToken()
		if err != nil {
			Error("[AUTHENTICATION FAILED] Please login first using: rune login <token>")
			return
		}

		c := client.New("http://localhost:8080", token)
		ns := normalizeNamespace(namespaceName)
		keys, err := c.List(ns)
		if err != nil {
			Error(err.Error())
			return
		}

		if len(keys) == 0 {
			Info("No stored keys found")
			return
		}

		fmt.Println("\nStored Secrets:\n")
		for _, key := range keys {
			fmt.Printf("  [SECRET] %s\n", key)
		}
	},
}

func init() {
	listCmd.Flags().StringVarP(&namespaceName, "namespace", "n", "", "Namespace")
	rootCmd.AddCommand(listCmd)
}
