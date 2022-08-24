package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var names []string

// helloCmd represents the hello command
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "hello world",
	Long:  `hello world again`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, n := range names {
			fmt.Println("hello", n)
		}
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
	helloCmd.Flags().StringArrayVarP(&names, "Name", "n", []string{"world"}, "Your name")
}
