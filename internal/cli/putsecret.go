package cli

import (
	"fmt"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{

	Use:   "put [key] [value]",
	Short: "Store a secret",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New("http://localhost:8080")
		err := c.Put(args[0], args[1])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Stored successfully")
	},
}

func init() {

	rootCmd.AddCommand(putCmd)

}
