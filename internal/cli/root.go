package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "rune",
	Short: "Nano Secret Manager",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
