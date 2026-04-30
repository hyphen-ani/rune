package cli

import (
	"fmt"
	"rune/pkg/client"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{

	Use:   "get [key]",
	Short: "Retrieve a secret",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New("http://localhost:8080")
		val, err := c.Get(args[0])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(val)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
