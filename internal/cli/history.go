package cli

import (
	"fmt"

	"rune/internal/config"
	"rune/internal/namespace"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var historyNamespace string

var historyCmd = &cobra.Command{
	Use:   "history [key]",
	Short: "Show the version history of a secret",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		token, err := config.LoadToken()
		if err != nil {
			Error("[AUTHENTICATION FAILED] Please login first using: rune login <token>")
			return
		}

		ns := namespace.Normalize(historyNamespace)

		c := client.New(
			"http://localhost:8080",
			token,
		)

		versions, err := c.ListVersions(
			args[0],
			ns,
		)

		if err != nil {
			Error(err.Error())
			return
		}

		if len(versions) == 0 {
			Info("No versions found")
			return
		}

		fmt.Println()
		fmt.Println("  SECRET HISTORY")
		fmt.Println("  ─────────────────────────────────────────────────────────")

		fmt.Printf(
			"  %-10s %-25s %-25s\n",
			"VERSION",
			"CREATED",
			"ROTATED",
		)

		fmt.Println("  ─────────────────────────────────────────────────────────")

		for _, version := range versions {

			rotated := version.RotatedAt

			if rotated == "" {
				rotated = "-"
			}

			fmt.Printf(
				"  %-10d %-25s %-25s\n",
				version.Version,
				version.CreatedAt,
				rotated,
			)
		}

		fmt.Println("  ─────────────────────────────────────────────────────────")
		fmt.Println()
	},
}

func init() {
	historyCmd.Flags().StringVarP(
		&historyNamespace,
		"namespace",
		"n",
		"default",
		"Namespace containing the secret",
	)

	rootCmd.AddCommand(historyCmd)
}
