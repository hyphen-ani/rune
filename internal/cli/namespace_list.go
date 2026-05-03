package cli

import (
	"fmt"
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var namespaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all namespaces",

	Run: func(cmd *cobra.Command, args []string) {

		token, err := config.LoadToken()
		if err != nil {
			Error("[INVALID REQUEST] Please Login First")
			return
		}

		c := client.New("http://localhost:8080", token)

		namespaces, err := c.ListNamespace()
		if err != nil {
			Error(err.Error())
			return
		}

		if len(namespaces) == 0 {
			Info("[NO NAMESPACES] No namespaces found")
			return
		}

		fmt.Println("\nNamespaces:\n")
		for _, namespace := range namespaces {
			fmt.Printf("[NAMESPACES] %s\n", namespace)
		}

		fmt.Println()
	},
}

func init() {
	namespaceCmd.AddCommand(namespaceListCmd)
}
