package cli

import (
	"github.com/spf13/cobra"
)

var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Manage namespaces for organizing and isolating secrets across applications and environments",
	Long: `Namespaces provide logical separation of secrets in Rune.
	They allow the same key to exist across different applications or environments without conflict (e.g., db/password in "springboot" vs "backend").
	Use namespaces to structure secrets for multi-tenant systems and better organization.`,
}

func init() {
	rootCmd.AddCommand(namespaceCmd)
}
