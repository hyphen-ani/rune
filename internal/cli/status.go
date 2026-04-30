package cli

import (
	"fmt"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{

	Use:   "status",
	Short: "Check vault status",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New("http://localhost:8080")
		status, err := c.Status()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Vault is", status)
	},
}

func init() {

	rootCmd.AddCommand(statusCmd)

}
