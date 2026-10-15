package cli

import (
	"fmt"
	"rune/internal/auth"
	"rune/internal/config"
	"rune/pkg/client"
	"strings"

	"github.com/spf13/cobra"
)

var tokenName string
var tokenNamespaceFlags []string // repeatable: "ns:read,write" or "ns" (defaults to all perms)

var tokenCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new access token",
	Long: `Create a new access token with fine-grained namespace permissions.

Examples:
  rune token create ci-deploy -n production:read,write
  rune token create ci-full   -n production:read,write,delete,rotate
  rune token create multi     -n production:read -n staging:read,write`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token, err := config.LoadToken()
		if err != nil {
			Error("[AUTHENTICATION FAILED] Please login first using: rune login <token>")
			return
		}

		name := tokenName
		if len(args) > 0 {
			name = args[0]
		}
		if name == "" {
			name = "default"
		}

		if len(tokenNamespaceFlags) == 0 {
			Error("[INVALID]: Provide at least one namespace with -n (e.g. -n default:read,write)")
			return
		}

		namespaces := parseNamespaceFlags(tokenNamespaceFlags)
		if namespaces == nil {
			return
		}

		c := client.New(serverURL(), token)
		newToken, record, err := c.CreateToken(name, namespaces)
		if err != nil {
			Error(err.Error())
			return
		}

		Success("Token created successfully")
		fmt.Println()
		fmt.Println("  ID:         ", record.ID)
		fmt.Println("  Name:       ", record.Name)
		fmt.Println("  Created At: ", record.CreatedAt)
		fmt.Println()
		for _, np := range record.Namespaces {
			fmt.Printf("  %-20s %s\n", np.Namespace, strings.Join(np.Permissions, ", "))
		}
		fmt.Println()
		fmt.Println("  Token:", newToken)
		fmt.Println()
	},
}

// parseNamespaceFlags converts ["ns1:read,write", "ns2:read"] to []auth.NamespacePermission.
// A value without a colon (e.g. "default") gets all four permissions.
// Returns nil and prints an error message if any flag is invalid.
func parseNamespaceFlags(flags []string) []auth.NamespacePermission {
	allPerms := []string{auth.PermRead, auth.PermWrite, auth.PermDelete, auth.PermRotate}
	var result []auth.NamespacePermission

	for _, f := range flags {
		parts := strings.SplitN(f, ":", 2)
		ns := strings.TrimSpace(parts[0])
		if ns == "" {
			Error(fmt.Sprintf("[INVALID]: empty namespace in flag %q", f))
			return nil
		}

		var perms []string
		if len(parts) == 1 {
			perms = allPerms
		} else {
			for _, p := range strings.Split(parts[1], ",") {
				p = strings.TrimSpace(p)
				if !auth.ValidPermissions[p] {
					Error(fmt.Sprintf("[INVALID]: unknown permission %q (valid: read, write, delete, rotate)", p))
					return nil
				}
				perms = append(perms, p)
			}
		}
		if len(perms) == 0 {
			Error(fmt.Sprintf("[INVALID]: no permissions specified for namespace %q", ns))
			return nil
		}

		result = append(result, auth.NamespacePermission{
			Namespace:   ns,
			Permissions: perms,
		})
	}

	return result
}

func init() {
	tokenCreateCmd.Flags().StringVar(&tokenName, "name", "", "Name of the token")
	tokenCreateCmd.Flags().StringArrayVarP(&tokenNamespaceFlags, "namespace", "n", nil,
		"Namespace and permissions: -n ns:read,write (repeatable for multiple namespaces)")
	tokenCmd.AddCommand(tokenCreateCmd)
}
