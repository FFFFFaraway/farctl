package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var name string

// helloCmd represents the hello command
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "hello world",
	Long:  `hello world again`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello", name)
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
	helloCmd.Flags().StringVarP(&name, "Name", "n", "world", "Your name")
}
