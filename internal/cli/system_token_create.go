package cli

import (
	"fmt"
	"rune/internal/config"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var sysTokenName      string
var sysTokenNamespace string

var systemTokenCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a read-only system token for SDK/programmatic access",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token, err := config.LoadToken()
		if err != nil {
			Error("[AUTHENTICATION FAILED] Please login first using: rune login <token>")
			return
		}

		name := sysTokenName
		if len(args) > 0 {
			name = args[0]
		}
		if name == "" {
			name = "system-default"
		}

		c := client.New("http://localhost:8080", token)
		newToken, record, err := c.CreateSystemToken(name, sysTokenNamespace)
		if err != nil {
			Error(err.Error())
			return
		}

		Success("System token created successfully")

		fmt.Println()
		fmt.Println("  ID:         ", record.ID)
		fmt.Println("  Name:       ", record.Name)
		fmt.Println("  Namespace:  ", record.Namespace)
		fmt.Println("  Read-only:  ", record.ReadOnly)
		fmt.Println("  Created At: ", record.CreatedAt)
		fmt.Println("  Token:      ", newToken)
		fmt.Println()
		fmt.Println("  Use this token in your SDK config. It cannot write or delete secrets.")
	},
}

func init() {
	systemTokenCreateCmd.Flags().StringVar(&sysTokenName, "name", "", "Name of the system token")
	systemTokenCreateCmd.Flags().StringVarP(&sysTokenNamespace, "namespace", "n", "default", "Namespace this token is scoped to")
	systemTokenCmd.AddCommand(systemTokenCreateCmd)
}
