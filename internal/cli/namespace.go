package cli

import "github.com/spf13/cobra"

var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Manage namespaces",
}

func init() {
	rootCmd.AddCommand(namespaceCmd)
}
