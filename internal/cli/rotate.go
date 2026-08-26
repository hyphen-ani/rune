package cli

import (
	"fmt"
	"rune/internal/config"
	"rune/internal/namespace"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var rotateNamespace string
var rotateCmd = &cobra.Command{
	Use:   "rotate [key]",
	Short: "Generate a new value for an existing secret",

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		token, err := config.LoadToken()
		if err != nil {
			Error("[AUTHENTICATION FAILED] Please login first using: rune login <token>")
			return
		}

		ns := namespace.Normalize(rotateNamespace)
		c := client.New("http://localhost:8080", token)

		value, err := c.Rotate(args[0], ns)
		if err != nil {
			Error(err.Error())
			return
		}

		fmt.Println()

		fmt.Println("  SECRET ROTATED")
		fmt.Println("  ─────────────────────────────────────────────")
		fmt.Printf("  %-12s %s\n", "Namespace:", ns)
		fmt.Printf("  %-12s %s\n", "Key:", args[0])
		fmt.Printf("  %-12s %s\n", "New Value:", value)
		fmt.Println("  ─────────────────────────────────────────────")
		fmt.Println()
		Success("Secret rotated successfully")

	},
}

func init() {
	rotateCmd.Flags().StringVarP(&rotateNamespace, "namespace", "n", "default", "Namespace containing the secret")
	rootCmd.AddCommand(rotateCmd)
}
